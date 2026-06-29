package monitor

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Checker struct {
	store  *MetricsStore
	alerts *AlertEngine
	peers  []string
}

func NewChecker(store *MetricsStore, alerts *AlertEngine) *Checker {
	return &Checker{
		store:  store,
		alerts: alerts,
	}
}

func (c *Checker) SetPeers(peers []string) {
	c.peers = peers
}

func (c *Checker) RunAll() {
	c.CheckHub()
	c.CheckContainers()
	if len(c.peers) > 0 {
		for _, peer := range c.peers {
			c.CheckPeerAlive(peer)
			c.CheckLatency(peer)
		}
	}
}

func (c *Checker) CheckHub() {
	// CPU
	cpu, err := readCPU()
	r := CheckResult{
		Type: CheckCPU, OK: err == nil, Value: cpu, Label: "hub",
		Timestamp: time.Now().Unix(), Message: errorMsg(err),
	}
	c.store.Push(r)
	c.alerts.Evaluate(r)

	// RAM
	ramPct, used, total, err := readRAM()
	msg := fmt.Sprintf("%s / %s", used, total)
	if err != nil {
		msg = err.Error()
	}
	r = CheckResult{
		Type: CheckRAM, OK: err == nil, Value: ramPct, Label: "hub",
		Timestamp: time.Now().Unix(), Message: msg,
	}
	c.store.Push(r)
	c.alerts.Evaluate(r)

	// Disk
	diskPct, dUsed, dTotal, err := readDisk("/data")
	if err != nil {
		diskPct, dUsed, dTotal, err = readDisk("/")
	}
	msg = fmt.Sprintf("%s / %s", dUsed, dTotal)
	if err != nil {
		msg = err.Error()
	}
	r = CheckResult{
		Type: CheckDisk, OK: err == nil, Value: diskPct, Label: "hub",
		Timestamp: time.Now().Unix(), Message: msg,
	}
	c.store.Push(r)
	c.alerts.Evaluate(r)

	// Public IP
	ip, err := readPublicIP()
	r = CheckResult{
		Type: CheckPublicIP, OK: err == nil, Value: 0, Label: "hub",
		Timestamp: time.Now().Unix(), Message: ip,
	}
	if err != nil {
		r.Message = err.Error()
	}
	c.store.Push(r)
}

func (c *Checker) CheckPeerAlive(label string) {
	alive, seconds := checkWGPeer(label)
	r := CheckResult{
		Type: CheckPeerAlive, OK: alive, Value: seconds, Label: label,
		Timestamp: time.Now().Unix(),
	}
	if !alive {
		r.Message = "no recent handshake"
	}
	c.store.Push(r)
	c.alerts.Evaluate(r)
}

func (c *Checker) CheckLatency(label string) {
	latency, err := pingPeer(label)
	r := CheckResult{
		Type: CheckLatency, OK: err == nil, Value: latency, Label: label,
		Timestamp: time.Now().Unix(), Message: errorMsg(err),
	}
	c.store.Push(r)
	c.alerts.Evaluate(r)
}

func (c *Checker) CheckContainers() {
	containers, err := listContainers()
	if err != nil {
		r := CheckResult{
			Type: CheckContainer, OK: false, Label: "docker",
			Timestamp: time.Now().Unix(), Message: fmt.Sprintf("docker ps: %v", err),
		}
		c.store.Push(r)
		return
	}

	for _, ct := range containers {
		r := CheckResult{
			Type: CheckContainer, OK: ct.Running, Value: float64(ct.Restarts), Label: ct.Name,
			Timestamp: time.Now().Unix(), Message: ct.Status,
		}
		if !ct.Running {
			r.Message = fmt.Sprintf("container %s is %s", ct.Name, ct.Status)
		}
		c.store.Push(r)
		c.alerts.Evaluate(r)

		if ct.Running {
			r2 := CheckResult{
				Type: CheckContainerCPU, OK: true, Value: ct.CPUPercent, Label: ct.Name,
				Timestamp: time.Now().Unix(),
			}
			c.store.Push(r2)

			r3 := CheckResult{
				Type: CheckContainerRAM, OK: true, Value: ct.MemPercent, Label: ct.Name,
				Timestamp: time.Now().Unix(), Message: ct.MemUsage,
			}
			c.store.Push(r3)
		}
	}
}

func readCPU() (float64, error) {
	out, err := exec.Command("sh", "-c", "top -bn1 2>/dev/null | grep 'Cpu(s)' | awk '{print $2}'").Output()
	if err != nil {
		out, err = exec.Command("sh", "-c", "cat /proc/loadavg | awk '{print $1}'").Output()
		if err != nil {
			return 0, err
		}
	}
	return strconv.ParseFloat(strings.TrimSpace(string(out)), 64)
}

func readRAM() (pct float64, used, total string, err error) {
	out, err := exec.Command("sh", "-c", "free -m | awk '/Mem:/{printf \"%.1f %sMB %sMB\", $3/$2*100, $3, $2}'").Output()
	if err != nil {
		return 0, "", "", err
	}
	parts := strings.Fields(string(out))
	if len(parts) < 2 {
		return 0, "", "", fmt.Errorf("parse free output")
	}
	pct, _ = strconv.ParseFloat(parts[0], 64)
	return pct, parts[1], parts[2], nil
}

func readDisk(path string) (pct float64, used, total string, err error) {
	out, err := exec.Command("df", "-B1", path).Output()
	if err != nil {
		return 0, "", "", err
	}
	lines := strings.Split(string(out), "\n")
	if len(lines) < 2 {
		return 0, "", "", fmt.Errorf("parse df output")
	}
	fields := strings.Fields(lines[1])
	if len(fields) < 5 {
		return 0, "", "", fmt.Errorf("parse df fields")
	}
	usedB, _ := strconv.ParseInt(fields[2], 10, 64)
	totalB, _ := strconv.ParseInt(fields[1], 10, 64)
	if totalB == 0 {
		return 0, "", "", nil
	}
	pct = float64(usedB) / float64(totalB) * 100
	return pct, formatBytes(usedB), formatBytes(totalB), nil
}

func readPublicIP() (string, error) {
	out, err := exec.Command("curl", "-s", "https://ifconfig.me").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func checkWGPeer(label string) (alive bool, seconds float64) {
	out, err := exec.Command("wg", "show", "dump").Output()
	if err != nil {
		return false, 0
	}
	for _, line := range strings.Split(string(out), "\n") {
		parts := strings.Fields(line)
		if len(parts) < 6 {
			continue
		}
		hs, _ := strconv.ParseInt(parts[4], 10, 64)
		if hs == 0 {
			continue
		}
		now := time.Now().Unix()
		diff := float64(now - hs)
		return diff < 180, diff
	}
	return false, 0
}

func pingPeer(label string) (float64, error) {
	out, err := exec.Command("wg", "show", "dump").Output()
	if err != nil {
		return 0, err
	}
	var peerIP string
	for _, line := range strings.Split(string(out), "\n") {
		parts := strings.Fields(line)
		if len(parts) >= 5 && parts[0] != "" {
			allowedIPs := parts[3]
			peerIP = strings.Split(allowedIPs, ",")[0]
			peerIP = strings.Split(peerIP, "/")[0]
			break
		}
	}
	if peerIP == "" {
		return 0, fmt.Errorf("no peer IP found")
	}
	out2, err := exec.Command("ping", "-c", "1", "-W", "2", peerIP).Output()
	if err != nil {
		return 0, err
	}
	line := string(out2)
	idx := strings.LastIndex(line, "=")
	if idx < 0 {
		return 0, fmt.Errorf("parse ping output")
	}
	end := strings.Index(line[idx+1:], " ")
	if end < 0 {
		end = len(line[idx+1:])
	}
	latency, err := strconv.ParseFloat(line[idx+1:idx+1+end], 64)
	return latency, err
}

func listContainers() ([]ContainerInfo, error) {
	out, err := exec.Command("docker", "ps", "-a",
		"--format", `{"name":"{{.Names}}","image":"{{.Image}}","status":"{{.Status}}","running":{{if eq .State "running"}}true{{else}}false{{end}}}`,
	).Output()
	if err != nil {
		return nil, err
	}

	var containers []ContainerInfo
	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	for scanner.Scan() {
		var ct ContainerInfo
		if err := json.Unmarshal(scanner.Bytes(), &ct); err != nil {
			continue
		}

		out2, err2 := exec.Command("docker", "stats", "--no-stream",
			"--format", `{"cpu":"{{.CPUPerc}}","mem":"{{.MemPerc}}","mem_usage":"{{.MemUsage}}"}`,
			ct.Name).Output()
		if err2 == nil {
			var stats struct {
				CPU     string `json:"cpu"`
				Mem     string `json:"mem"`
				MemUsage string `json:"mem_usage"`
			}
			if json.Unmarshal(out2, &stats) == nil {
				ct.CPUPercent, _ = strconv.ParseFloat(strings.TrimSuffix(stats.CPU, "%"), 64)
				ct.MemPercent, _ = strconv.ParseFloat(strings.TrimSuffix(stats.Mem, "%"), 64)
				ct.MemUsage = stats.MemUsage
			}
		}

		out3, _ := exec.Command("docker", "inspect", "-f", "{{.RestartCount}}", ct.Name).Output()
		ct.Restarts, _ = strconv.Atoi(strings.TrimSpace(string(out3)))

		containers = append(containers, ct)
	}
	return containers, scanner.Err()
}

func errorMsg(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

func ListContainersExternal() ([]ContainerInfo, error) {
	return listContainers()
}

func formatBytes(b int64) string {
	switch {
	case b > 1<<30:
		return fmt.Sprintf("%.1fG", float64(b)/(1<<30))
	case b > 1<<20:
		return fmt.Sprintf("%.1fM", float64(b)/(1<<20))
	case b > 1<<10:
		return fmt.Sprintf("%.1fK", float64(b)/(1<<10))
	default:
		return fmt.Sprintf("%d", b)
	}
}

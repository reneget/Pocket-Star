package tui

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/anomalyco/my-pretty-star/modules/monitor"
	"github.com/anomalyco/my-pretty-star/pkg/log"
	"github.com/anomalyco/my-pretty-star/pkg/module"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type PeerInfo struct {
	Name      string
	IP        string
	Connected bool
	Rx        string
	Tx        string
}

type StatusInfo struct {
	Connected bool
	HubAddr   string
	MyIP      string
	Uptime    string
}

type HealthData struct {
	CPU     float64
	RAM     float64
	RAMUsed string
	RAMTotal string
	Disk    float64
	DiskUsed string
	DiskTotal string
	PubIP   string
}

type model struct {
	ready    bool
	tab      int
	err      error
	ctx      *module.Context
	quitting bool

	status  StatusInfo
	peers   []PeerInfo

	healthData  HealthData
	containers  []monitor.ContainerInfo

	logView viewport.Model
	logBuf  *log.RingBuffer

	commandMode bool
	command     string

	width  int
	height int

	showHelp bool
}

func initialModel(ctx *module.Context) model {
	vp := viewport.New(0, 0)
	vp.Style = lipgloss.NewStyle().PaddingLeft(1)

	return model{
		tab:     0,
		ctx:     ctx,
		logBuf:  log.Default,
		logView: vp,
		status: StatusInfo{
			Connected: false,
			HubAddr:   "N/A",
			MyIP:      "N/A",
			Uptime:    "N/A",
		},
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		tickCmd(),
		tea.EnterAltScreen,
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
		m.logView.Width = msg.Width - 6
		m.logView.Height = msg.Height - 12
		return m, nil

	case tea.KeyMsg:
		return m.handleKeyMsg(msg)

	case tea.MouseMsg:
		return m.handleMouseMsg(msg)

	case tickMsg:
		return m, tea.Batch(
			tickCmd(),
			func() tea.Msg { return m.refreshStatus() },
			m.refreshHealth(),
			m.refreshContainers(),
		)

	case statusMsg:
		m.status = StatusInfo(msg)
		return m, nil

	case peersMsg:
		m.peers = []PeerInfo(msg)
		return m, nil

	case healthMsg:
		m.healthData = HealthData(msg)
		return m, nil

	case containersMsg:
		m.containers = []monitor.ContainerInfo(msg)
		return m, nil

	case errMsg:
		m.err = error(msg)
		return m, nil

	case tea.QuitMsg:
		m.quitting = true
		return m, tea.Quit
	}

	return m, nil
}

func (m *model) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.showHelp {
		if msg.String() == "?" || msg.String() == "esc" {
			m.showHelp = false
		}
		return m, nil
	}

	if m.commandMode {
		switch msg.Type {
		case tea.KeyEnter:
			m.executeCommand()
			m.commandMode = false
			m.command = ""
			return m, nil
		case tea.KeyEscape:
			m.commandMode = false
			m.command = ""
			return m, nil
		case tea.KeyBackspace:
			if len(m.command) > 0 {
				m.command = m.command[:len(m.command)-1]
			}
			return m, nil
		case tea.KeyRunes:
			m.command += string(msg.Runes)
			return m, nil
		}
		return m, nil
	}

	switch msg.String() {
	case "q", "ctrl+c":
		m.quitting = true
		return m, tea.Quit
	case "?":
		m.showHelp = true
		return m, nil
	case "1":
		m.tab = 0
	case "2":
		m.tab = 1
	case "3":
		m.tab = 2
	case "r":
		return m, func() tea.Msg { return m.refreshStatus() }
	case "/":
		m.commandMode = true
		m.command = ""
		return m, nil
	case "tab":
		m.tab = (m.tab + 1) % 3
	}

	return m, nil
}

func (m *model) handleMouseMsg(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	if msg.Type != tea.MouseLeft {
		return m, nil
	}

	x, y := msg.X, msg.Y

	if y == 1 {
		if x >= 4 && x <= 14 {
			m.tab = 0
		} else if x >= 16 && x <= 26 {
			m.tab = 1
		} else if x >= 28 && x <= 38 {
			m.tab = 2
		} else if x >= 40 && x <= 50 {
			m.showHelp = true
		}
	}

	return m, nil
}

func (m *model) executeCommand() {
	parts := strings.Fields(m.command)
	if len(parts) == 0 {
		return
	}
	cmd := parts[0]

	switch cmd {
	case "q", "quit", "exit":
		m.quitting = true
	case "c", "connect":
		log.Default.Push("connect: not yet implemented")
	case "d", "disconnect":
		log.Default.Push("disconnect: not yet implemented")
	case "r", "refresh":
		m.refreshStatus()
	case "h", "help":
		m.showHelp = true
	case "tab":
		if len(parts) > 1 {
			switch parts[1] {
			case "1":
				m.tab = 0
			case "2":
				m.tab = 1
			case "3":
				m.tab = 2
			}
		}
	default:
		log.Default.Push(fmt.Sprintf("unknown command: /%s", cmd))
	}
}

func (m model) View() string {
	if !m.ready {
		return "Loading..."
	}

	if m.quitting {
		return "Goodbye!\n"
	}

	if m.showHelp {
		return m.helpView()
	}

	s := strings.Builder{}

	s.WriteString(m.headerView())
	s.WriteString("\n\n")

	switch m.tab {
	case 0:
		s.WriteString(m.dashboardView())
	case 1:
		s.WriteString(m.peersView())
	case 2:
		s.WriteString(m.logsView())
	}

	s.WriteString("\n")
	s.WriteString(m.footerView())

	return appStyle.Render(s.String())
}

func (m model) headerView() string {
	title := headerStyle.Render("★ Pocket Star")

	var statusStr string
	if m.status.Connected {
		statusStr = connectedStyle.Render("● Connected")
	} else {
		statusStr = disconnectedStyle.Render("○ Disconnected")
	}

	right := lipgloss.JoinHorizontal(
		lipgloss.Top,
		statusStyle.Render("v0.1.0"),
		statusStyle.Render("│"),
		statusStr,
	)

	topLine := lipgloss.JoinHorizontal(
		lipgloss.Top,
		title,
		lipgloss.NewStyle().Width(m.width-50).Render(""),
		right,
	)

	tabs := lipgloss.JoinHorizontal(
		lipgloss.Top,
		styleTab(0, m.tab, " Dashboard "),
		styleTab(1, m.tab, " Peers "),
		styleTab(2, m.tab, " Logs "),
		tabStyle.Render("│ ? Help │"),
	)

	divider := helpStyle.Width(m.width - 8).Render(strings.Repeat("─", m.width-8))

	return fmt.Sprintf("%s\n%s\n%s", topLine, divider, tabs)
}

func styleTab(idx, active int, label string) string {
	if idx == active {
		return activeTabStyle.Render(label)
	}
	return tabStyle.Render(label)
}

func (m model) footerView() string {
	var cmdLine string
	if m.commandMode {
		cmdLine = commandPromptStyle.Render("/") + commandInputStyle.Render(m.command+" ")
	} else {
		cmdLine = helpStyle.Render("q: quit  ?: help  /: command  tab: switch  r: refresh")
	}

	divider := strings.Repeat("─", m.width-6)
	return fmt.Sprintf("%s\n%s", divider, cmdLine)
}

func tickCmd() tea.Cmd {
	return tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m model) refreshStatus() tea.Msg {
	peers := fetchPeers()
	if len(peers) > 0 {
		m.logBuf.Push(fmt.Sprintf("status: %d peer(s) connected", len(peers)))
	}
	return peersMsg(peers)
}

func (m model) refreshHealth() tea.Cmd {
	return func() tea.Msg {
		return healthMsg(m.fetchSystemHealth())
	}
}

func (m model) refreshContainers() tea.Cmd {
	return func() tea.Msg {
		cts, _ := monitor.ListContainersExternal()
		return containersMsg(cts)
	}
}

func (m model) fetchSystemHealth() HealthData {
	h := HealthData{}

	out, err := exec.Command("sh", "-c", `top -bn1 2>/dev/null | grep 'Cpu(s)' | awk '{print $2}'`).Output()
	if err == nil {
		s := strings.TrimSpace(string(out))
		h.CPU, _ = strconv.ParseFloat(s, 64)
	}

	out, err = exec.Command("sh", "-c", `free -m | awk '/Mem:/{printf "%.1f %sMB %sMB", $3/$2*100, $3, $2}'`).Output()
	if err == nil {
		parts := strings.Fields(string(out))
		if len(parts) >= 3 {
			h.RAM, _ = strconv.ParseFloat(parts[0], 64)
			h.RAMUsed = parts[1]
			h.RAMTotal = parts[2]
		}
	}

	out, err = exec.Command("sh", "-c", "df -B1 /data 2>/dev/null || df -B1 /").Output()
	if err == nil {
		lines := strings.Split(string(out), "\n")
		if len(lines) >= 2 {
			fields := strings.Fields(lines[1])
			if len(fields) >= 5 {
				totalB, _ := strconv.ParseInt(fields[1], 10, 64)
				usedB, _ := strconv.ParseInt(fields[2], 10, 64)
				if totalB > 0 {
					h.Disk = float64(usedB) / float64(totalB) * 100
					h.DiskUsed = fmt.Sprintf("%.1fG", float64(usedB)/(1<<30))
					h.DiskTotal = fmt.Sprintf("%.1fG", float64(totalB)/(1<<30))
				}
			}
		}
	}

	out, err = exec.Command("curl", "-s", "https://ifconfig.me").Output()
	if err == nil {
		h.PubIP = strings.TrimSpace(string(out))
	}

	return h
}

func fetchPeers() []PeerInfo {
	out, err := exec.Command("wg", "show", "dump").Output()
	if err != nil {
		return nil
	}

	var peers []PeerInfo
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		parts := strings.Fields(line)
		if len(parts) < 5 {
			continue
		}
		p := PeerInfo{
			Name:      parts[0][:min(15, len(parts[0]))],
			IP:        parts[1],
			Connected: parts[2] != "0",
			Rx:        formatBytes(parts[4]),
			Tx:        formatBytes(parts[5]),
		}
		peers = append(peers, p)
	}
	return peers
}

func formatBytes(s string) string {
	var b int64
	fmt.Sscanf(s, "%d", &b)
	switch {
	case b > 1<<30:
		return fmt.Sprintf("%.1f GiB", float64(b)/(1<<30))
	case b > 1<<20:
		return fmt.Sprintf("%.1f MiB", float64(b)/(1<<20))
	case b > 1<<10:
		return fmt.Sprintf("%.1f KiB", float64(b)/(1<<10))
	default:
		return fmt.Sprintf("%d B", b)
	}
}

type tickMsg time.Time
type statusMsg StatusInfo
type peersMsg []PeerInfo
type healthMsg HealthData
type containersMsg []monitor.ContainerInfo
type errMsg error

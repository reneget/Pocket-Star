package vpn

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"text/template"

	"golang.org/x/crypto/curve25519"
)

const hubNetwork = "10.0.0.0/24"

type KeyPair struct {
	Private string
	Public  string
}

func GenerateKeyPair() (*KeyPair, error) {
	var private [32]byte
	if _, err := rand.Read(private[:]); err != nil {
		return nil, fmt.Errorf("generate private key: %w", err)
	}

	private[0] &= 248
	private[31] &= 127
	private[31] |= 64

	public, err := curve25519.X25519(private[:], curve25519.Basepoint)
	if err != nil {
		return nil, fmt.Errorf("derive public key: %w", err)
	}

	return &KeyPair{
		Private: base64.StdEncoding.EncodeToString(private[:]),
		Public:  base64.StdEncoding.EncodeToString(public),
	}, nil
}

type hubPeer struct {
	Name   string
	Public string
	IP     string
}

type hubData struct {
	HubPrivate string
	Address    string
	Port       int
	Peers      []hubPeer
}

type nodeData struct {
	Private        string
	Address        string
	HubPublic      string
	HubEndpoint    string
	AllowedNetwork string
}

const hubTemplate = `[Interface]
PrivateKey = {{.HubPrivate}}
Address = {{.Address}}
ListenPort = {{.Port}}

{{range .Peers}}
# {{.Name}}
[Peer]
PublicKey = {{.Public}}
AllowedIPs = {{.IP}}/32
{{end}}
`

const nodeTemplate = `[Interface]
PrivateKey = {{.Private}}
Address = {{.Address}}

[Peer]
PublicKey = {{.HubPublic}}
Endpoint = {{.HubEndpoint}}
AllowedIPs = {{.AllowedNetwork}}
PersistentKeepalive = 25
`

func GenerateHubConfig(address string, listenPort int) (hubCfg string, nodeCfgs map[string]string, err error) {
	hubKeys, err := GenerateKeyPair()
	if err != nil {
		return "", nil, err
	}

	defaultNodes := []string{"server", "work-laptop"}
	nodeCfgs = make(map[string]string)
	ipIndex := 2
	var peers []hubPeer

	for _, name := range defaultNodes {
		nodeKeys, err := GenerateKeyPair()
		if err != nil {
			return "", nil, fmt.Errorf("generate keys for %s: %w", name, err)
		}

		nodeIP := fmt.Sprintf("10.0.0.%d", ipIndex)
		peers = append(peers, hubPeer{
			Name:   name,
			Public: nodeKeys.Public,
			IP:     nodeIP,
		})

		nd := nodeData{
			Private:        nodeKeys.Private,
			Address:        nodeIP + "/24",
			HubPublic:      hubKeys.Public,
			HubEndpoint:    fmt.Sprintf("CHANGE_TO_HUB_IP:%d", listenPort),
			AllowedNetwork: hubNetwork,
		}

		var buf bytes.Buffer
		tmpl := template.Must(template.New("node").Parse(nodeTemplate))
		if err := tmpl.Execute(&buf, nd); err != nil {
			return "", nil, fmt.Errorf("render node %s: %w", name, err)
		}
		nodeCfgs[name] = buf.String()
		ipIndex++
	}

	hd := hubData{
		HubPrivate: hubKeys.Private,
		Address:    address,
		Port:       listenPort,
		Peers:      peers,
	}

	var hubBuf bytes.Buffer
	tmpl := template.Must(template.New("hub").Parse(hubTemplate))
	if err := tmpl.Execute(&hubBuf, hd); err != nil {
		return "", nil, fmt.Errorf("render hub: %w", err)
	}

	return hubBuf.String(), nodeCfgs, nil
}

func GenerateNodeConfig(name, hubEndpoint, allowedNet string) string {
	keys, err := GenerateKeyPair()
	if err != nil {
		return fmt.Sprintf("# Error generating keys: %v", err)
	}
	return fmt.Sprintf(`[Interface]
PrivateKey = %s
Address = 10.0.0.X/24

[Peer]
PublicKey = <HUB_PUBLIC_KEY>
Endpoint = %s
AllowedIPs = %s
PersistentKeepalive = 25
`, keys.Private, hubEndpoint, allowedNet)
}

func ShowStatus() error {
	cmd := exec.Command("wg", "show")
	out, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("wg show failed (is WireGuard running?): %w", err)
	}
	fmt.Print(string(out))
	return nil
}

type HubConfig struct {
	PrivateKey string
	Address    string
	Port       int
	Peers      []hubPeer
}

func ParseHubConfig(data string) (*HubConfig, error) {
	hc := &HubConfig{Port: 51820}
	lines := strings.Split(data, "\n")
	inInterface := false
	var lastComment string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "#") {
			lastComment = strings.TrimPrefix(line, "# ")
			continue
		}
		if line == "[Interface]" {
			inInterface = true
			continue
		}
		if line == "[Peer]" {
			inInterface = false
			hc.Peers = append(hc.Peers, hubPeer{Name: lastComment})
			lastComment = ""
			continue
		}
		if inInterface {
			if strings.HasPrefix(line, "PrivateKey") {
				hc.PrivateKey = extractValue(line)
			} else if strings.HasPrefix(line, "Address") {
				hc.Address = extractValue(line)
			} else if strings.HasPrefix(line, "ListenPort") {
				portStr := extractValue(line)
				if p, err := strconv.Atoi(portStr); err == nil {
					hc.Port = p
				}
			}
		}
		if len(hc.Peers) > 0 {
			idx := len(hc.Peers) - 1
			if strings.HasPrefix(line, "PublicKey") {
				hc.Peers[idx].Public = extractValue(line)
			} else if strings.HasPrefix(line, "AllowedIPs") {
				val := extractValue(line)
				hc.Peers[idx].IP = strings.Split(val, "/")[0]
			}
		}
	}
	return hc, nil
}

func extractValue(line string) string {
	parts := strings.SplitN(line, "=", 2)
	if len(parts) < 2 {
		return ""
	}
	return strings.TrimSpace(parts[1])
}

func AddPeerToHubConfig(hubConfigPath, peerName string) (nodeCfg string, peerIP string, err error) {
	data, err := os.ReadFile(hubConfigPath)
	if err != nil {
		return "", "", fmt.Errorf("read hub config: %w", err)
	}

	hc, err := ParseHubConfig(string(data))
	if err != nil {
		return "", "", fmt.Errorf("parse hub config: %w", err)
	}

	nextIP := 2
	for _, p := range hc.Peers {
		parts := strings.Split(p.IP, ".")
		if len(parts) == 4 {
			if n, err := strconv.Atoi(parts[3]); err == nil && n >= nextIP {
				nextIP = n + 1
			}
		}
	}

	peerIP = fmt.Sprintf("10.0.0.%d", nextIP)

	peerKeys, err := GenerateKeyPair()
	if err != nil {
		return "", "", fmt.Errorf("generate peer keys: %w", err)
	}

	hc.Peers = append(hc.Peers, hubPeer{
		Name:   peerName,
		Public: peerKeys.Public,
		IP:     peerIP,
	})

	hd := hubData{
		HubPrivate: hc.PrivateKey,
		Address:    hc.Address,
		Port:       hc.Port,
		Peers:      hc.Peers,
	}

	var hubBuf bytes.Buffer
	tmpl := template.Must(template.New("hub").Parse(hubTemplate))
	if err := tmpl.Execute(&hubBuf, hd); err != nil {
		return "", "", fmt.Errorf("render updated hub config: %w", err)
	}

	if err := os.WriteFile(hubConfigPath, hubBuf.Bytes(), 0600); err != nil {
		return "", "", fmt.Errorf("write hub config: %w", err)
	}

	nodeCfg = fmt.Sprintf(`[Interface]
PrivateKey = %s
Address = %s/24

[Peer]
PublicKey = %s
Endpoint = CHANGE_TO_HUB_IP:%d
AllowedIPs = 10.0.0.0/24
PersistentKeepalive = 25
`, peerKeys.Private, peerIP, hc.PrivateKey, hc.Port)

	return nodeCfg, peerIP, nil
}

func GenerateNodeConfigFull(peerName, hubEndpoint, hubPublicKey, nodePrivateKey, nodeIP string) string {
	return fmt.Sprintf(`[Interface]
PrivateKey = %s
Address = %s/24

[Peer]
PublicKey = %s
Endpoint = %s
AllowedIPs = 10.0.0.0/24
PersistentKeepalive = 25
`, nodePrivateKey, nodeIP, hubPublicKey, hubEndpoint)
}

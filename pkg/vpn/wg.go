package vpn

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os/exec"
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

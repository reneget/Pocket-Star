package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/anomalyco/my-pretty-star/pkg/docker"
	"github.com/anomalyco/my-pretty-star/pkg/vpn"
	"github.com/spf13/cobra"
)

var nodeCmd = &cobra.Command{
	Use:   "node",
	Short: "Node (server/client) management commands",
}

var nodeJoinCmd = &cobra.Command{
	Use:   "join",
	Short: "Connect this node to a hub",
	RunE: func(cmd *cobra.Command, args []string) error {
		hubAddr, _ := cmd.Flags().GetString("hub")
		useDocker, _ := cmd.Flags().GetBool("docker")
		useMonitor, _ := cmd.Flags().GetBool("monitor")
		dataDir, _ := cmd.Flags().GetString("data-dir")
		nodeName, _ := cmd.Flags().GetString("name")
		autoMode, _ := cmd.Flags().GetBool("auto")
		sshKey, _ := cmd.Flags().GetString("ssh-key")

		if hubAddr == "" {
			return fmt.Errorf("--hub is required (e.g. --hub user@1.2.3.4)")
		}
		if nodeName == "" {
			nodeName = "node"
		}
		if dataDir == "" {
			dataDir = "./pstar-data"
		}
		if err := os.MkdirAll(dataDir, 0755); err != nil {
			return fmt.Errorf("create data dir: %w", err)
		}

		if autoMode {
			return nodeJoinAuto(hubAddr, nodeName, dataDir, sshKey, useDocker, useMonitor)
		}

		fmt.Print("Enter master password (if set on hub): ")
		var masterPass string
		fmt.Scanln(&masterPass)

		cfg := vpn.GenerateNodeConfig(nodeName, hubAddr, "10.0.0.0/24")
		cfgPath := filepath.Join(dataDir, "node.conf")
		if err := os.WriteFile(cfgPath, []byte(cfg), 0600); err != nil {
			return fmt.Errorf("write node config: %w", err)
		}
		fmt.Printf("✓ Node config: %s\n", cfgPath)

		if masterPass != "" {
			importCmd := fmt.Sprintf("scp %s:./pstar-data/clients/%s.conf.enc . && pstar decrypt %s.conf.enc",
				hubAddr, nodeName, nodeName)
			fmt.Printf("📌 Import from hub: %s\n", importCmd)
		} else {
			fmt.Printf("📌 Import from hub: scp %s:./pstar-data/clients/%s.conf ./\n", hubAddr, nodeName)
		}

		if useDocker {
			compose := docker.GenerateNodeCompose(extractHost(hubAddr))
			composePath := filepath.Join(dataDir, "docker-compose.yml")
			if err := os.WriteFile(composePath, []byte(compose), 0644); err != nil {
				return fmt.Errorf("write docker-compose: %w", err)
			}
			fmt.Printf("✓ Docker Compose: %s\n", composePath)
		}

		fmt.Println("\n✔ Node configured. Apply the config to connect.")
		return nil
	},
}

var nodeStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show node connection status",
	RunE: func(cmd *cobra.Command, args []string) error {
		return vpn.ShowStatus()
	},
}

func nodeJoinAuto(hubAddr, nodeName, dataDir, sshKey string, useDocker, useMonitor bool) error {
	host := hubAddr
	if !strings.Contains(host, "@") {
		host = "root@" + host
	}

	sshBase := []string{"ssh"}
	if sshKey != "" {
		sshBase = append(sshBase, "-i", sshKey)
	}
	sshBase = append(sshBase, "-o", "StrictHostKeyChecking=accept-new")

	fmt.Printf("→ Adding peer %s to hub %s...\n", nodeName, host)
	sshCmd := append(sshBase, host, fmt.Sprintf("pstar hub add-peer --name %s --data-dir ./pstar-data", nodeName))

	out, err := exec.Command(sshCmd[0], sshCmd[1:]...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("ssh add-peer failed: %w\n%s", err, string(out))
	}
	fmt.Printf("  %s", string(out))

	scpSrc := fmt.Sprintf("%s:./pstar-data/clients/%s.conf", host, nodeName)
	scpDst := filepath.Join(dataDir, nodeName+".conf")

	scpArgs := []string{}
	if sshKey != "" {
		scpArgs = append(scpArgs, "-i", sshKey)
	}
	scpArgs = append(scpArgs, "-o", "StrictHostKeyChecking=accept-new", scpSrc, scpDst)

	fmt.Printf("→ Downloading config...\n")
	out, err = exec.Command("scp", scpArgs...).CombinedOutput()
	if err != nil {
		if strings.Contains(string(out), "No such file") {
			// try .conf.enc
			scpSrcEnc := fmt.Sprintf("%s:./pstar-data/clients/%s.conf.enc", host, nodeName)
			scpDstEnc := filepath.Join(dataDir, nodeName+".conf.enc")
			scpArgs[len(scpArgs)-2] = scpSrcEnc
			scpArgs[len(scpArgs)-1] = scpDstEnc
			out, err = exec.Command("scp", scpArgs...).CombinedOutput()
			if err != nil {
				return fmt.Errorf("scp config failed: %w\n%s", err, string(out))
			}
			fmt.Printf("  Downloaded encrypted config. Run: pstar decrypt %s\n", scpDstEnc)
		} else {
			return fmt.Errorf("scp config failed: %w\n%s", err, string(out))
		}
	} else {
		fmt.Printf("✓ Config: %s\n", scpDst)
	}

	fmt.Printf("→ Applying WireGuard config...\n")
	applyPath := filepath.Join(dataDir, nodeName+".conf")
	if _, err := os.Stat(applyPath); err == nil {
		apply := exec.Command("wg-quick", "up", applyPath)
		if out, err := apply.CombinedOutput(); err != nil {
			fmt.Printf("  wg-quick up: %s\n", string(out))
		} else {
			fmt.Printf("✓ WireGuard tunnel is up\n")
		}
	}

	if useDocker {
		fmt.Printf("→ Generating docker-compose...\n")
		compose := docker.GenerateNodeCompose(extractHost(hubAddr))
		composePath := filepath.Join(dataDir, "docker-compose.yml")
		if err := os.WriteFile(composePath, []byte(compose), 0644); err != nil {
			return fmt.Errorf("write docker-compose: %w", err)
		}
		fmt.Printf("✓ Docker Compose: %s\n", composePath)
	}

	if useMonitor {
		fmt.Printf("→ Generating compose with Pulse agent...\n")
		compose := docker.GenerateNodeComposeWithMonitor(extractHost(hubAddr))
		composePath := filepath.Join(dataDir, "docker-compose.yml")
		if err := os.WriteFile(composePath, []byte(compose), 0644); err != nil {
			return fmt.Errorf("write docker-compose: %w", err)
		}
		fmt.Printf("✓ Docker Compose (with Pulse): %s\n", composePath)
	}

	fmt.Println("\n✔ Node configured and connected!")
	return nil
}

func extractHost(addr string) string {
	if strings.Contains(addr, "@") {
		parts := strings.Split(addr, "@")
		return parts[len(parts)-1]
	}
	if strings.Contains(addr, ":") {
		parts := strings.Split(addr, ":")
		return parts[0]
	}
	return addr
}

func init() {
	nodeCmd.AddCommand(nodeJoinCmd)
	nodeCmd.AddCommand(nodeStatusCmd)

	nodeJoinCmd.Flags().String("hub", "", "Hub address (e.g. user@1.2.3.4)")
	nodeJoinCmd.Flags().String("name", "", "Node name")
	nodeJoinCmd.Flags().Bool("docker", false, "Generate docker-compose.yml")
	nodeJoinCmd.Flags().Bool("monitor", false, "Generate compose with Pulse agent")
	nodeJoinCmd.Flags().Bool("auto", false, "Auto-join via SSH/SCP")
	nodeJoinCmd.Flags().String("ssh-key", "", "SSH key path for auto-join")
	nodeJoinCmd.Flags().String("data-dir", "./pstar-data", "Data directory")
}

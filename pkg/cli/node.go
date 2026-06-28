package cli

import (
	"fmt"
	"os"
	"path/filepath"

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
		dataDir, _ := cmd.Flags().GetString("data-dir")
		nodeName, _ := cmd.Flags().GetString("name")

		if hubAddr == "" {
			return fmt.Errorf("--hub is required (e.g. --hub 1.2.3.4:51820)")
		}
		if nodeName == "" {
			nodeName = "node"
		}
		if dataDir == "" {
			dataDir = "./mps-data"
		}
		if err := os.MkdirAll(dataDir, 0755); err != nil {
			return fmt.Errorf("create data dir: %w", err)
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
			importCmd := fmt.Sprintf("scp user@%s:./mps-data/clients/%s.conf.enc . && mps decrypt %s.conf.enc",
				hubAddr, nodeName, nodeName)
			fmt.Printf("📌 Import from hub: %s\n", importCmd)
		} else {
			fmt.Printf("📌 Import from hub: scp user@%s:./mps-data/clients/%s.conf ./\n", hubAddr, nodeName)
		}

		if useDocker {
			compose := docker.GenerateNodeCompose(hubAddr)
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

func init() {
	nodeCmd.AddCommand(nodeJoinCmd)
	nodeCmd.AddCommand(nodeStatusCmd)
	nodeJoinCmd.Flags().String("hub", "", "Hub address (e.g. 1.2.3.4:51820)")
	nodeJoinCmd.Flags().String("name", "", "Node name")
	nodeJoinCmd.Flags().Bool("docker", false, "Generate docker-compose.yml")
	nodeJoinCmd.Flags().String("data-dir", "./mps-data", "Data directory")
}

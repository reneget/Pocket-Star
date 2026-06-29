package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/anomalyco/my-pretty-star/pkg/config"
	"github.com/anomalyco/my-pretty-star/pkg/docker"
	"github.com/anomalyco/my-pretty-star/pkg/vpn"
	"github.com/spf13/cobra"
)

var hubCmd = &cobra.Command{
	Use:   "hub",
	Short: "Hub (VPS) management commands",
}

var hubInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize hub on VPS",
	RunE: func(cmd *cobra.Command, args []string) error {
		useDocker, _ := cmd.Flags().GetBool("docker")
		usePass, _ := cmd.Flags().GetBool("pass")
		useMonitor, _ := cmd.Flags().GetBool("monitor")
		dataDir, _ := cmd.Flags().GetString("data-dir")

		if dataDir == "" {
			dataDir = "./pstar-data"
		}
		if err := os.MkdirAll(dataDir, 0755); err != nil {
			return fmt.Errorf("create data dir: %w", err)
		}

		hubCfg, nodeCfgs, err := vpn.GenerateHubConfig("10.0.0.1/24", 51820)
		if err != nil {
			return fmt.Errorf("generate configs: %w", err)
		}

		hubPath := filepath.Join(dataDir, "hub.conf")
		if err := os.WriteFile(hubPath, []byte(hubCfg), 0600); err != nil {
			return fmt.Errorf("write hub config: %w", err)
		}
		fmt.Printf("✓ Hub config written: %s\n", hubPath)

		var masterPass string
		if usePass {
			masterPass, err = readPassword("Enter master password: ")
			if err != nil {
				return err
			}
			confirm, err := readPassword("Confirm master password: ")
			if err != nil {
				return err
			}
			if masterPass != confirm {
				return fmt.Errorf("passwords do not match")
			}
		}

		clientsDir := filepath.Join(dataDir, "clients")
		if err := os.MkdirAll(clientsDir, 0755); err != nil {
			return fmt.Errorf("create clients dir: %w", err)
		}

		for name, nodeCfg := range nodeCfgs {
			ext := ".conf"
			cfgData := []byte(nodeCfg)

			if masterPass != "" {
				encData, err := config.Encrypt(cfgData, masterPass)
				if err != nil {
					return fmt.Errorf("encrypt config for %s: %w", name, err)
				}
				cfgData = encData
				ext = ".conf.enc"
			}

			path := filepath.Join(clientsDir, name+ext)
			if err := os.WriteFile(path, cfgData, 0600); err != nil {
				return fmt.Errorf("write %s config: %w", name, err)
			}
			fmt.Printf("✓ Config for %s: %s\n", name, path)
		}

		if useDocker {
			var compose string
			if useMonitor {
				compose = docker.GenerateHubComposeWithMonitor(masterPass != "")
			} else {
				compose = docker.GenerateHubCompose(masterPass != "")
			}
			composePath := filepath.Join(dataDir, "docker-compose.yml")
			if err := os.WriteFile(composePath, []byte(compose), 0644); err != nil {
				return fmt.Errorf("write docker-compose: %w", err)
			}
			fmt.Printf("✓ Docker Compose: %s\n", composePath)
		}

		fmt.Println("\n✔ Hub initialized. Share client configs securely with your nodes.")
		return nil
	},
}

var hubStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show hub status and connected peers",
	RunE: func(cmd *cobra.Command, args []string) error {
		return vpn.ShowStatus()
	},
}

func init() {
	hubCmd.AddCommand(hubInitCmd)
	hubCmd.AddCommand(hubStatusCmd)
	hubInitCmd.Flags().Bool("docker", false, "Generate docker-compose.yml")
	hubInitCmd.Flags().Bool("pass", false, "Protect node configs with master password")
	hubInitCmd.Flags().Bool("monitor", false, "Include monitoring (Uptime Kuma) in docker-compose")
	hubInitCmd.Flags().String("data-dir", "./pstar-data", "Data directory for configs")
}

func readPassword(prompt string) (string, error) {
	fmt.Print(prompt)
	var pass string
	_, err := fmt.Scanln(&pass)
	return pass, err
}

package cli

import (
	"fmt"
	"os"

	"github.com/anomalyco/my-pretty-star/pkg/icon"
	"github.com/anomalyco/my-pretty-star/pkg/module"
	"github.com/anomalyco/my-pretty-star/modules/tui"
	"github.com/spf13/cobra"
)

var cliMode bool

var rootCmd = &cobra.Command{
	Use:   "pstar",
	Short: "Pocket Star — site-to-site VPN orchestrator",
	Long: `Pocket Star turns your old PC into a personal server
by connecting it to a cheap VPS via site-to-site VPN (star topology).

Documentation: https://github.com/reneget/Pocket-Star`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if cliMode {
			return cmd.Help()
		}
		return tui.Run(&module.Context{
			DataDir: "./pstar-data",
			Network: "10.0.0.0/24",
		})
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&cliMode, "cli", false, "Force CLI mode instead of TUI")
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(doctorCmd)
	rootCmd.AddCommand(hubCmd)
	rootCmd.AddCommand(nodeCmd)
	rootCmd.AddCommand(moduleCmd)
	rootCmd.AddCommand(decryptCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Pocket Star v0.1.0\n")
		fmt.Printf("License: MIT\n")
		fmt.Printf("Logo:    embedded (%d bytes)\n", len(icon.Logo))
	},
}

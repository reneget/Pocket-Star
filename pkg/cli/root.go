package cli

import (
	"fmt"
	"os"

	"github.com/anomalyco/my-pretty-star/pkg/icon"
	"github.com/anomalyco/my-pretty-star/pkg/module"
	"github.com/spf13/cobra"
)

var cliMode bool
var monEnabled bool

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

		ctx := &module.Context{
			DataDir: "./pstar-data",
			Network: "10.0.0.0/24",
		}

		mgr := module.NewManager(ctx)
		if err := mgr.InitAll(); err != nil {
			return fmt.Errorf("init modules: %w", err)
		}

		if mod := mgr.Get("tui"); mod != nil {
			if tw, ok := mod.(interface{ SetModuleManager(*module.Manager) }); ok {
				tw.SetModuleManager(mgr)
			}
		}

		if monEnabled {
			if err := mgr.Enable("monitor"); err != nil {
				return fmt.Errorf("enable monitor: %w", err)
			}
		}

		if err := mgr.StartAll(); err != nil {
			return fmt.Errorf("start modules: %w", err)
		}

		if tuiMod := mgr.Get("tui"); tuiMod != nil {
			if tw, ok := tuiMod.(interface{ Wait() }); ok {
				tw.Wait()
			}
		}
		return nil
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
	rootCmd.PersistentFlags().BoolVar(&monEnabled, "monitor", false, "Enable monitoring module at startup")
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(doctorCmd)
	rootCmd.AddCommand(hubCmd)
	rootCmd.AddCommand(nodeCmd)
	rootCmd.AddCommand(moduleCmd)
	rootCmd.AddCommand(decryptCmd)
	rootCmd.AddCommand(monitorCmd)
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

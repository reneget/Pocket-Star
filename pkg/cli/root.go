package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "pstar",
	Short: "Pocket Star — site-to-site VPN orchestrator",
	Long: `Pocket Star turns your old PC into a personal server
by connecting it to a cheap VPS via site-to-site VPN (star topology).

Documentation: https://github.com/reneget/Pocket-Star`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
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
		fmt.Println("pstar v0.1.0")
	},
}

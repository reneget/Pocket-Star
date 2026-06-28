package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mps",
	Short: "My Pretty Star — site-to-site VPN orchestrator",
	Long: `My Pretty Star turns your old PC into a personal server
by connecting it to a cheap VPS via site-to-site VPN (star topology).

Documentation: https://github.com/anomalyco/my-pretty-star`,
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
		fmt.Println("mps v0.1.0")
	},
}

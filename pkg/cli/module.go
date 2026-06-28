package cli

import (
	"fmt"

	"github.com/anomalyco/my-pretty-star/pkg/module"
	"github.com/spf13/cobra"
)

var moduleCmd = &cobra.Command{
	Use:   "module",
	Short: "Manage modules",
}

var moduleListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available modules",
	RunE: func(cmd *cobra.Command, args []string) error {
		registry := module.GetRegistry()
		if len(registry) == 0 {
			fmt.Println("No modules registered.")
			return nil
		}
		for _, m := range registry {
			fmt.Printf("  %s\n", m.Name())
		}
		return nil
	},
}

func init() {
	moduleCmd.AddCommand(moduleListCmd)
}

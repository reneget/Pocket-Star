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
		ctx := &module.Context{DataDir: "./pstar-data"}
		mgr := module.NewManager(ctx)
		mgr.InitAll()
		fmt.Print(module.ModuleListString(mgr.List()))
		return nil
	},
}

func init() {
	moduleCmd.AddCommand(moduleListCmd)
}

package cli

import (
	"fmt"

	"github.com/anomalyco/my-pretty-star/modules/monitor"
	"github.com/spf13/cobra"
)

var monitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "System and peer monitoring",
}

var monitorStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current health status",
	RunE: func(cmd *cobra.Command, args []string) error {
		store := monitor.DefaultStore
		if store == nil {
			return fmt.Errorf("monitor not initialized")
		}
		checks := store.Recent(50)
		fmt.Printf("Latest %d checks:\n\n", len(checks))
		for _, r := range checks {
			status := "✔"
			if !r.OK {
				status = "✘"
			}
			msg := ""
			if r.Message != "" {
				msg = " (" + r.Message + ")"
			}
			fmt.Printf("  %s %s %s = %.1f%s\n", status, r.Type, r.Label, r.Value, msg)
		}
		return nil
	},
}

var monitorLogsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Show monitor logs from file",
	RunE: func(cmd *cobra.Command, args []string) error {
		store := monitor.DefaultStore
		if store == nil {
			return fmt.Errorf("monitor not initialized")
		}
		n, _ := cmd.Flags().GetInt("lines")
		results, err := store.ReadFileLines(n)
		if err != nil {
			return fmt.Errorf("read logs: %w", err)
		}
		for _, r := range results {
			status := "✔"
			if !r.OK {
				status = "✘"
			}
			fmt.Printf("[%d] %s %s %s = %.1f %s\n",
				r.Timestamp, status, r.Type, r.Label, r.Value, r.Message)
		}
		return nil
	},
}

var monitorConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure monitor settings",
	RunE: func(cmd *cobra.Command, args []string) error {
		store := monitor.DefaultStore
		if store == nil {
			return fmt.Errorf("monitor not initialized")
		}
		maxLog, _ := cmd.Flags().GetInt64("max-log")
		if maxLog > 0 {
			store.SetMaxBytes(maxLog * 1024 * 1024)
			fmt.Printf("Max log size set to %d MB\n", maxLog)
		} else {
			fmt.Println("Usage: pstar monitor config --max-log 20")
		}
		return nil
	},
}

var monitorContainersCmd = &cobra.Command{
	Use:   "containers",
	Short: "Show Docker container status",
	RunE: func(cmd *cobra.Command, args []string) error {
		showRestart, _ := cmd.Flags().GetBool("restart")
		containers, err := monitor.ListContainersExternal()
		if err != nil {
			return fmt.Errorf("list containers: %w", err)
		}
		for _, ct := range containers {
			if showRestart && ct.Restarts == 0 {
				continue
			}
			status := "●"
			if !ct.Running {
				status = "○"
			}
			fmt.Printf("  %s %-20s CPU:%5s RAM:%5s Restarts:%d %s\n",
				status, ct.Name,
				fmt.Sprintf("%.1f%%", ct.CPUPercent),
				fmt.Sprintf("%.1f%%", ct.MemPercent),
				ct.Restarts, ct.Status)
		}
		return nil
	},
}

var monitorAlertsCmd = &cobra.Command{
	Use:   "alerts",
	Short: "Show active alert rules",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Alert rules active:")
		fmt.Println("  peer.alive  → warning: >120s, critical: >300s")
		fmt.Println("  peer.latency → warning: >100ms, critical: >500ms")
		fmt.Println("  system.cpu  → warning: >70%, critical: >90%")
		fmt.Println("  system.ram  → warning: >80%, critical: >95%")
		fmt.Println("  system.disk → warning: >80%, critical: >95%")
		fmt.Println("  docker.cpu  → warning: >80%, critical: >95%")
		return nil
	},
}

func init() {
	monitorCmd.AddCommand(monitorStatusCmd)
	monitorCmd.AddCommand(monitorLogsCmd)
	monitorCmd.AddCommand(monitorConfigCmd)
	monitorCmd.AddCommand(monitorContainersCmd)
	monitorCmd.AddCommand(monitorAlertsCmd)

	monitorLogsCmd.Flags().IntP("lines", "n", 50, "Number of log lines to show")
	monitorConfigCmd.Flags().Int64("max-log", 0, "Max log file size in MB")
	monitorContainersCmd.Flags().Bool("restart", false, "Show only restarting containers")
}

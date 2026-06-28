package cli

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Run system diagnostics",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Running diagnostics...")
		fmt.Println()

		checks := []struct {
			name string
			check func() error
		}{
			{"WireGuard (wg)", checkWG},
			{"WireGuard tools (wg-quick)", checkWGQuick},
			{"IP forwarding", checkIPForward},
			{"TUN module", checkTUN},
		}

		allOK := true
		for _, c := range checks {
			err := c.check()
			status := "✔"
			if err != nil {
				status = "✘"
				allOK = false
			}
			fmt.Printf("  %s %s", status, c.name)
			if err != nil {
				fmt.Printf(" (%v)", err)
			}
			fmt.Println()
		}

		if allOK {
			fmt.Println("\n✔ All checks passed.")
		} else {
			fmt.Println("\n⚠ Some checks failed. See messages above.")
		}
		return nil
	},
}

func checkWG() error {
	return exec.Command("wg", "version").Run()
}

func checkWGQuick() error {
	return exec.Command("wg-quick", "--version").Run()
}

func checkIPForward() error {
	return exec.Command("sh", "-c", "sysctl net.ipv4.ip_forward | grep -q 1").Run()
}

func checkTUN() error {
	return exec.Command("sh", "-c", "test -c /dev/net/tun || modprobe tun 2>/dev/null").Run()
}

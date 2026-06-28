package cli

import (
	"fmt"
	"os"

	"github.com/anomalyco/my-pretty-star/pkg/config"
	"github.com/spf13/cobra"
)

var decryptCmd = &cobra.Command{
	Use:   "decrypt <file>",
	Short: "Decrypt an encrypted config file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		encPath := args[0]
		data, err := os.ReadFile(encPath)
		if err != nil {
			return fmt.Errorf("read file: %w", err)
		}

		fmt.Print("Enter master password: ")
		var pass string
		fmt.Scanln(&pass)

		plain, err := config.Decrypt(data, pass)
		if err != nil {
			return fmt.Errorf("decrypt failed (wrong password?): %w", err)
		}

		outPath := encPath
		if len(outPath) > 4 && outPath[len(outPath)-4:] == ".enc" {
			outPath = outPath[:len(outPath)-4]
		} else {
			outPath += ".decrypted"
		}

		if err := os.WriteFile(outPath, plain, 0600); err != nil {
			return fmt.Errorf("write decrypted file: %w", err)
		}
		fmt.Printf("✓ Decrypted: %s\n", outPath)
		return nil
	},
}

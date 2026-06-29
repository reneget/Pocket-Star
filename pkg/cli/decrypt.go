package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/anomalyco/my-pretty-star/pkg/config"
	"github.com/spf13/cobra"
)

var decryptCmd = &cobra.Command{
	Use:   "decrypt <file>",
	Short: "Decrypt an encrypted config file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read file: %w", err)
		}

		fmt.Print("Enter master password: ")
		var pass string
		fmt.Scanln(&pass)

		plain, err := config.Decrypt(data, pass)
		if err != nil {
			return fmt.Errorf("decrypt: %w", err)
		}

		outPath := filepath.Join(filepath.Dir(path), filepath.Base(path)+".decrypted")
		if err := os.WriteFile(outPath, plain, 0600); err != nil {
			return fmt.Errorf("write output: %w", err)
		}
		fmt.Printf("✓ Decrypted: %s\n", outPath)
		return nil
	},
}

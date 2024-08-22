// cobra cli setup files
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var useSha256 bool
var beQuiet bool
var tempDir string

var rootCmd = &cobra.Command{
	Use:   "safelock-cli",
	Short: "Simple tool to encrypt/decrypt files with AES encryption",
	Long:  "Simple command-line tool to encrypt and decrypt files with AES encryption",
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&useSha256, "sha256", false, "use SHA256 (faster) instead of SHA512")
	rootCmd.PersistentFlags().BoolVar(&beQuiet, "quiet", false, "disable output logs")
	rootCmd.PersistentFlags().StringVar(&tempDir, "temp-dir", os.TempDir(), "directory for temporary files")
}

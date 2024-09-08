// cobra cli setup files
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var useSha256 bool
var useSha512 bool
var beQuiet bool

var rootCmd = &cobra.Command{
	Use:     "safelock-cli",
	Short:   "Simple tool to encrypt/decrypt files with AES encryption",
	Long:    "Simple command-line tool to encrypt and decrypt files with AES encryption",
	Version: "0.4.3",
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
	rootCmd.PersistentFlags().BoolVar(&useSha256, "sha256", false, "use sha256 for hashing")
	rootCmd.PersistentFlags().BoolVar(&useSha512, "sha512", false, "use sha512 for hashing")
	rootCmd.PersistentFlags().BoolVar(&beQuiet, "quiet", false, "disable output logs")
}

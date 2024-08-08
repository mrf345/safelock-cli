package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mrf345/safelock-cli/encryption"
	"github.com/mrf345/safelock-cli/utils"
)

var encryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "encrypt [file or directory path] [encrypted file path]",
	Long:  "encrypt [file or directory path] [encrypted file path]",
	Run: func(cmd *cobra.Command, args []string) {
		const example = "example: safelock-cli encrypt test.txt encrypted.bin"
		var err error
		var password string

		switch len(args) {
		case 0:
			utils.PrintErrsAndExit("missing input and output file paths", example)
		case 1:
			utils.PrintErrsAndExit("missing output file path", example)
		case 2:
			break
		default:
			utils.PrintErrsAndExit("too many arguments", example)
		}

		if password, err = utils.GetPassword(); err != nil {
			utils.PrintErrsAndExit(err.Error())
		}

		options := encryption.GetDefaultEncryptionOptions()
		inputPath, outputPath := args[0], args[1]

		if err = encryption.Encrypt(inputPath, outputPath, password, options); err != nil {
			utils.PrintErrsAndExit(err.Error())
		}

		fmt.Printf("Encrypted: %s\n", outputPath)
	},
}

func init() {
	rootCmd.AddCommand(encryptCmd)
}

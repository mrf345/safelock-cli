package cmd

import (
	"github.com/mrf345/safelock-cli/encryption"
	"github.com/mrf345/safelock-cli/utils"
	"github.com/spf13/cobra"
)

var decryptCmd = &cobra.Command{
	Use:   "decrypt",
	Short: "Decrypt an encrypted file",
	Long:  "Decrypt an encrypted file",
	Run: func(cmd *cobra.Command, args []string) {
		const example = "example: safelock-cli decrypt encrypted.bin decrypted_files"
		var err error
		var password string

		switch len(args) {
		case 0:
			utils.PrintErrsAndExit("missing input and output file paths", example)
		case 1:
			utils.PrintErrsAndExit("missing output path", example)
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

		if err = encryption.Decrypt(inputPath, outputPath, password, options); err != nil {
			utils.PrintErrsAndExit(err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(decryptCmd)
}

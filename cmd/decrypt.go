package cmd

import (
	"context"

	"github.com/mrf345/safelock-cli/safelock"
	"github.com/mrf345/safelock-cli/utils"
	"github.com/spf13/cobra"
)

var decryptCmd = &cobra.Command{
	Use:   "decrypt",
	Short: "decrypt [encrypted file path] [directory path]",
	Long:  "decrypt [encrypted file path] [directory path]",
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		var pwd string
		var sl *safelock.Safelock
		const example = "example: safelock-cli decrypt encrypted.bin decrypted_files"

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

		if useSha256 {
			sl = safelock.NewSha256()
		} else {
			sl = safelock.New()
		}

		if pwd, err = utils.GetPassword(sl.MinPasswordLength); err != nil {
			utils.PrintErrsAndExit(err.Error())
		}

		sl.Quiet = beQuiet
		inputPath, outputPath := args[0], args[1]

		if err = sl.Decrypt(context.TODO(), inputPath, outputPath, pwd); err != nil {
			utils.PrintErrsAndExit(err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(decryptCmd)
}

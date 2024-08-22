package cmd

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/mrf345/safelock-cli/safelock"
	"github.com/mrf345/safelock-cli/utils"
)

var encryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "encrypt [file or directory path] [encrypted file path]",
	Long:  "encrypt [file or directory path] [encrypted file path]",
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		var pwd string
		var sl *safelock.Safelock
		const example = "example: safelock-cli encrypt test.txt encrypted.bin"

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

		if useSha256 {
			sl = safelock.NewSha256()
		} else {
			sl = safelock.New()
		}

		if pwd, err = utils.GetPassword(sl.MinPasswordLength); err != nil {
			utils.PrintErrsAndExit(err.Error())
		}

		sl.Quiet = beQuiet
		sl.Registry.TempDir = tempDir
		inputPath, outputPath := []string{args[0]}, args[1]

		if err = sl.Encrypt(context.TODO(), inputPath, outputPath, pwd); err != nil {
			utils.PrintErrsAndExit(err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(encryptCmd)
}

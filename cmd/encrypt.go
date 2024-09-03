package cmd

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	"github.com/mrf345/safelock-cli/safelock"
	"github.com/mrf345/safelock-cli/slErrs"
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
		var outputFile *os.File
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
		} else if useSha512 {
			sl = safelock.NewSha512()
		} else {
			sl = safelock.New()
		}

		if pwd, err = utils.GetPassword(sl.MinPasswordLength); err != nil {
			utils.PrintErrsAndExit(err.Error())
		}

		sl.Quiet = beQuiet
		inputPath, outputPath := []string{args[0]}, args[1]

		if outputFile, err = os.OpenFile(outputPath, os.O_RDWR|os.O_CREATE, 0776); err != nil {
			utils.PrintErrsAndExit((&slErrs.ErrInvalidOutputPath{
				Path: outputPath,
				Err:  err,
			}).Error())
		}
		defer outputFile.Close()

		if err = sl.Encrypt(context.TODO(), inputPath, outputFile, pwd); err != nil {
			utils.PrintErrsAndExit(err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(encryptCmd)
}

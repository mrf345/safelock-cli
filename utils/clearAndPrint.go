package utils

import (
	"fmt"

	"github.com/inancgumus/screen"
	myTypes "github.com/mrf345/safelock-cli/types"
)

func ClearAndPrint(options myTypes.EncryptionOptions, msg string, params ...any) {
	if !options.Quiet {
		screen.Clear()
		screen.MoveTopLeft()
		fmt.Printf(msg, params...)
	}
}

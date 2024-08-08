package utils

import (
	"fmt"
	"os"

	"github.com/inancgumus/screen"
)

func PrintErrsAndExit(errs ...string) {
	screen.Clear()
	screen.MoveTopLeft()

	for _, e := range errs {
		fmt.Println(e)
	}
	os.Exit(1)
}

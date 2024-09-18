package safelock

import (
	"fmt"
	"runtime"

	"github.com/inancgumus/screen"
)

func (sl *Safelock) log(msg string, params ...any) {
	if !sl.Quiet {
		if runtime.GOOS != "windows" {
			screen.Clear()
			screen.MoveTopLeft()
		}
		fmt.Printf(msg, params...)
	}
}

func (sl *Safelock) logStatus(status StatusItem) {
	if status.Event == StatusUpdate {
		sl.log("%s (%.2f%%)\n", status.Msg, status.Percent)
	}
}

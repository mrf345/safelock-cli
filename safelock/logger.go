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

func (sl *Safelock) logStatus(status string, percent float64) {
	sl.log("%s (%.2f)\n", status, percent)
}

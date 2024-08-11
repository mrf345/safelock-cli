package safelock

import (
	"fmt"

	"github.com/inancgumus/screen"
)

func (sl *Safelock) log(msg string, params ...any) {
	if !sl.Quiet {
		screen.Clear()
		screen.MoveTopLeft()
		fmt.Printf(msg, params...)
	}
}

func (sl *Safelock) logStatus(status, percent string) {
	sl.log("%s (%s)\n", status, percent)
}

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

func (sl *Safelock) logStatus(status string, percent float64) {
	sl.log("%s (%.2f)\n", status, percent)
}

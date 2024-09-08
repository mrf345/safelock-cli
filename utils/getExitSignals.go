package utils

import (
	"os"
	"os/signal"
	"syscall"
)

// get exit signal channel, and its closer function
func GetExitSignals() (<-chan os.Signal, func()) {
	signals := make(chan os.Signal, 2)
	close := func() {
		signal.Stop(signals)
		close(signals)
	}

	signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT)
	return signals, close
}

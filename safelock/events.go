package safelock

// [safelock.Safelock.StatusObs] streaming event keys type
type StatusEvent string

// [safelock.Safelock.StatusObs] streaming event keys
const (
	StatusStart  StatusEvent = "start_status"  // encryption/decryption has started (no args)
	StatusEnd    StatusEvent = "end_status"    // encryption/decryption has ended (no args)
	StatusUpdate StatusEvent = "update_status" // new status update (status string, percent float64)
	StatusError  StatusEvent = "error_status"  // encryption/decryption failed (error)
)

// return event key value as string
func (se StatusEvent) Str() string {
	return string(se)
}

func (sl *Safelock) updateStatus(status string, percent float64) {
	sl.StatusObs.Trigger(StatusUpdate.Str(), status, percent)
}

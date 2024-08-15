package safelock

// [safelock.Safelock.StatusObs] streaming event keys
const (
	EventStatusStart  = "start_status"  // encryption/decryption has started (no args)
	EventStatusEnd    = "end_status"    // encryption/decryption has ended (no args)
	EventStatusUpdate = "update_status" // new status update (status string, percent float64)
	EventStatusError  = "error_status"  // encryption/decryption failed (error)
)

func (sl *Safelock) updateStatus(status string, percent float64) {
	sl.StatusObs.Trigger(EventStatusUpdate, status, percent)
}

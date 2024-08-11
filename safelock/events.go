package safelock

// [safelock.Safelock.StatusObs] streaming event keys
const (
	EventStatusStart  = "start_status"  // encryption/decryption has started (no args)
	EventStatusEnd    = "end_status"    // encryption/decryption has ended (no args)
	EventStatusUpdate = "update_status" // new status update (status, percent string)
	EventStatusError  = "error_status"  // encryption/decryption failed (error)
)

func (sl *Safelock) updateStatus(status, percent string) {
	sl.StatusObs.Trigger(EventStatusUpdate, status, percent)
}

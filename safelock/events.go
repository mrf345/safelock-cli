package safelock

import "sync"

// [safelock.Safelock.StatusObs] streaming event keys type
type StatusEvent string

// [safelock.Safelock.StatusObs] streaming event keys
const (
	StatusStart  StatusEvent = "start_status"  // encryption/decryption has started
	StatusEnd    StatusEvent = "end_status"    // encryption/decryption has ended
	StatusUpdate StatusEvent = "update_status" // new status update
	StatusError  StatusEvent = "error_status"  // encryption/decryption failed
)

// return event key value as string
func (se StatusEvent) Str() string {
	return string(se)
}

// item used to communicate status changes
type StatusItem struct {
	// status change event key
	Event StatusEvent
	// completion percent
	Percent float64
	// optional status change text
	Msg string
	// optional status change error
	Err error
}

// observable like data structure used to stream status changes
type StatusObservable struct {
	mu      sync.RWMutex
	subs    map[int]func(StatusItem)
	counter int
}

// creates a new [safelock.StatusObservable] instance
func NewStatusObs() *StatusObservable {
	return &StatusObservable{
		subs: make(map[int]func(StatusItem)),
	}
}

// adds a new status change subscriber, and returns the unsubscribe function
func (obs *StatusObservable) Subscribe(callback func(StatusItem)) func() {
	obs.mu.Lock()
	id := obs.counter
	obs.subs[id] = callback
	obs.counter += 1
	obs.mu.Unlock()

	// returns unsubscribe function
	return func() {
		obs.mu.Lock()
		delete(obs.subs, id)
		obs.mu.Unlock()
	}
}

// clears all subscriptions
func (obs *StatusObservable) Unsubscribe() {
	obs.mu.Lock()
	clear(obs.subs)
	obs.mu.Unlock()
}

func (obs *StatusObservable) next(value StatusItem) {
	obs.mu.RLock()
	for _, callback := range obs.subs {
		go callback(value)
	}
	obs.mu.RUnlock()
}

func (sl *Safelock) updateStatus(status string, percent float64) {
	sl.StatusObs.next(StatusItem{
		Event:   StatusUpdate,
		Msg:     status,
		Percent: percent,
	})
}

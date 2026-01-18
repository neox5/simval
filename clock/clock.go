package clock

import "time"

// Publisher provides a subscription interface for typed values.
type Publisher[T any] interface {
	Subscribe() <-chan T
}

// ClockStats contains observable metrics for a Clock.
type ClockStats struct {
	TickCount uint64
	IsRunning bool
	Interval  time.Duration
}

// Clock provides timing signals for value updates.
type Clock interface {
	Publisher[struct{}]
	Start()
	Stop()
	Stats() ClockStats
}

package clock

import (
	"sync/atomic"
	"time"
)

// PeriodicClock generates ticks at fixed intervals.
type PeriodicClock struct {
	interval  time.Duration
	ticker    *time.Ticker
	tickChan  chan struct{}
	stop      chan struct{}
	tickCount atomic.Uint64
	running   atomic.Bool
}

// NewPeriodicClock creates a new clock that ticks at the specified interval.
func NewPeriodicClock(interval time.Duration) *PeriodicClock {
	return &PeriodicClock{
		interval: interval,
		tickChan: make(chan struct{}),
		stop:     make(chan struct{}),
	}
}

// Start begins generating ticks.
func (c *PeriodicClock) Start() {
	c.ticker = time.NewTicker(c.interval)
	c.running.Store(true)
	go c.run()
}

func (c *PeriodicClock) run() {
	for {
		select {
		case <-c.ticker.C:
			c.tickCount.Add(1)
			select {
			case c.tickChan <- struct{}{}:
			case <-c.stop:
				return
			}
		case <-c.stop:
			return
		}
	}
}

// Stop stops the clock and closes the tick channel.
func (c *PeriodicClock) Stop() {
	if c.ticker != nil {
		c.ticker.Stop()
	}
	c.running.Store(false)
	close(c.stop)
	close(c.tickChan)
}

// Subscribe returns the channel that receives tick events.
func (c *PeriodicClock) Subscribe() <-chan struct{} {
	return c.tickChan
}

// Stats returns current clock metrics.
func (c *PeriodicClock) Stats() ClockStats {
	return ClockStats{
		TickCount: c.tickCount.Load(),
		IsRunning: c.running.Load(),
		Interval:  c.interval,
	}
}

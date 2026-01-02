package clock

import "time"

// Clock provides timing signals for value updates.
type Clock interface {
	Start()
	Stop()
	Tick() <-chan time.Time
}

// PeriodicClock generates ticks at fixed intervals.
type PeriodicClock struct {
	interval time.Duration
	ticker   *time.Ticker
	tickChan chan time.Time
	stop     chan struct{}
}

// NewPeriodicClock creates a new clock that ticks at the specified interval.
func NewPeriodicClock(interval time.Duration) *PeriodicClock {
	return &PeriodicClock{
		interval: interval,
		tickChan: make(chan time.Time),
		stop:     make(chan struct{}),
	}
}

// Start begins generating ticks.
func (c *PeriodicClock) Start() {
	c.ticker = time.NewTicker(c.interval)
	go c.run()
}

func (c *PeriodicClock) run() {
	for {
		select {
		case t := <-c.ticker.C:
			select {
			case c.tickChan <- t:
			case <-c.stop:
				return
			}
		case <-c.stop:
			return
		}
	}
}

// Stop stops the clock.
func (c *PeriodicClock) Stop() {
	if c.ticker != nil {
		c.ticker.Stop()
	}
	close(c.stop)
	close(c.tickChan)
}

// Tick returns the channel that receives tick events.
func (c *PeriodicClock) Tick() <-chan time.Time {
	return c.tickChan
}

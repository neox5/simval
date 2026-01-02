package value

import (
	"sync"

	"github.com/neox5/simval/internal/clock"
	"github.com/neox5/simval/internal/source"
	"github.com/neox5/simval/internal/transform"
)

// Value represents a simulated value that changes over time.
type Value[T any] struct {
	mu      sync.RWMutex
	current T

	clock      clock.Clock
	source     source.NumberSource[T]
	transforms []transform.Transformation[T]
	stop       chan struct{}
}

// New creates a new Value with the given clock, source, and optional transforms.
// The Value automatically starts its internal goroutine.
func New[T any](
	clk clock.Clock,
	src source.NumberSource[T],
	transforms ...transform.Transformation[T],
) *Value[T] {
	v := &Value[T]{
		clock:      clk,
		source:     src,
		transforms: transforms,
		stop:       make(chan struct{}),
	}

	go v.run()
	return v
}

func (v *Value[T]) run() {
	for {
		select {
		case <-v.clock.Tick():
			// Generate base value
			value := v.source.Next()

			// Apply transforms
			for _, transform := range v.transforms {
				value = transform.Apply(value)
			}

			// Store result
			v.mu.Lock()
			v.current = value
			v.mu.Unlock()

		case <-v.stop:
			return
		}
	}
}

// Value returns the current value.
func (v *Value[T]) Value() T {
	v.mu.RLock()
	defer v.mu.RUnlock()
	return v.current
}

// Stop stops the value's internal goroutine.
func (v *Value[T]) Stop() {
	close(v.stop)
}

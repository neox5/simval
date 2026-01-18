package source

import (
	"sync"
	"sync/atomic"

	"github.com/neox5/simv/clock"
)

// ConstSource always returns the same value.
type ConstSource[T any] struct {
	clock clock.Clock
	value T

	initOnce        sync.Once
	clockChan       <-chan struct{}
	mu              sync.Mutex
	subscribers     []chan T
	generationCount atomic.Uint64
}

// NewConstSource creates a source that always returns the given value.
func NewConstSource[T any](clk clock.Clock, value T) *ConstSource[T] {
	return &ConstSource[T]{
		clock: clk,
		value: value,
	}
}

// Subscribe returns a channel that receives constant values on each clock tick.
func (s *ConstSource[T]) Subscribe() <-chan T {
	s.initOnce.Do(func() {
		s.clockChan = s.clock.Subscribe()
		go s.run()
	})

	s.mu.Lock()
	defer s.mu.Unlock()

	ch := make(chan T)
	s.subscribers = append(s.subscribers, ch)
	return ch
}

func (s *ConstSource[T]) run() {
	for range s.clockChan {
		value := s.value
		s.generationCount.Add(1)

		s.mu.Lock()
		subs := s.subscribers
		s.mu.Unlock()

		for _, subChan := range subs {
			subChan <- value
		}
	}

	// Clock closed, close all subscriber channels
	s.mu.Lock()
	for _, subChan := range s.subscribers {
		close(subChan)
	}
	s.mu.Unlock()
}

// Stats returns current source metrics.
func (s *ConstSource[T]) Stats() SourceStats {
	s.mu.Lock()
	subCount := len(s.subscribers)
	s.mu.Unlock()

	return SourceStats{
		GenerationCount: s.generationCount.Load(),
		SubscriberCount: subCount,
	}
}

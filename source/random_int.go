package source

import (
	"math/rand/v2"
	"sync"
	"sync/atomic"

	"github.com/neox5/simv/clock"
	"github.com/neox5/simv/seed"
)

// RandomIntSource generates random integers within a range [min, max].
type RandomIntSource struct {
	clock    clock.Clock
	min, max int
	rng      *rand.Rand

	initOnce        sync.Once
	clockChan       <-chan struct{}
	mu              sync.Mutex
	subscribers     []chan int
	generationCount atomic.Uint64
}

// NewRandomIntSource creates a source that generates random integers
// in the inclusive range [min, max].
// Uses the global seed registry for deterministic sequences when seeded.
func NewRandomIntSource(clk clock.Clock, min, max int) *RandomIntSource {
	seed1, seed2 := seed.Next()
	return &RandomIntSource{
		clock: clk,
		min:   min,
		max:   max,
		rng:   rand.New(rand.NewPCG(seed1, seed2)),
	}
}

// Subscribe returns a channel that receives random integers on each clock tick.
func (s *RandomIntSource) Subscribe() <-chan int {
	s.initOnce.Do(func() {
		s.clockChan = s.clock.Subscribe()
		go s.run()
	})

	s.mu.Lock()
	defer s.mu.Unlock()

	ch := make(chan int)
	s.subscribers = append(s.subscribers, ch)
	return ch
}

func (s *RandomIntSource) run() {
	for range s.clockChan {
		value := s.min + s.rng.IntN(s.max-s.min+1)
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
func (s *RandomIntSource) Stats() SourceStats {
	s.mu.Lock()
	subCount := len(s.subscribers)
	s.mu.Unlock()

	return SourceStats{
		GenerationCount: s.generationCount.Load(),
		SubscriberCount: subCount,
	}
}

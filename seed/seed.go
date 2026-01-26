package seed

import (
	"math/rand/v2"
	"sync"
	"time"
)

var (
	globalRegistry *registry
	registryOnce   sync.Once
)

// registry provides deterministic seed sequences.
type registry struct {
	mu         sync.Mutex
	masterSeed uint64
	nextStream uint64
	autoInit   bool
}

// Init initializes the global registry with a master seed for repeatable simulations.
// If not called, the registry auto-initializes with a time-based seed on first use.
func Init(masterSeed uint64) {
	registryOnce.Do(func() {
		globalRegistry = &registry{
			masterSeed: masterSeed,
			nextStream: 0,
			autoInit:   false,
		}
	})
}

// NewRand returns a new independent random number generator.
// Each call returns an RNG with seeds (masterSeed, streamN) where N increments.
// Auto-initializes with time-based seed if Init was not called.
func NewRand() *rand.Rand {
	registryOnce.Do(func() {
		now := uint64(time.Now().UnixNano())
		globalRegistry = &registry{
			masterSeed: now,
			nextStream: 0,
			autoInit:   true,
		}
	})

	return globalRegistry.newRand()
}

// Current returns the active seed state for logging and reproducibility.
// Returns (masterSeed, streamCounter, autoInitialized) where:
// - masterSeed: The seed value (from Init() or time-based if auto-initialized)
// - streamCounter: Current stream counter (number of NewRand() calls made)
// - autoInitialized: true if Init() was never called (time-based seed)
//
// For reproducibility, call Init(masterSeed) before creating any sources.
// Each source will receive RNGs with seeds (masterSeed, 0), (masterSeed, 1), etc.
func Current() (masterSeed, streamCounter uint64, autoInitialized bool) {
	if globalRegistry == nil {
		return 0, 0, false
	}

	globalRegistry.mu.Lock()
	defer globalRegistry.mu.Unlock()

	return globalRegistry.masterSeed, globalRegistry.nextStream, globalRegistry.autoInit
}

func (r *registry) newRand() *rand.Rand {
	r.mu.Lock()
	defer r.mu.Unlock()

	seed1 := r.masterSeed
	seed2 := r.nextStream
	r.nextStream++

	return rand.New(rand.NewPCG(seed1, seed2))
}

package seed

import (
	"math/rand/v2"
	"sync"
	"time"
)

var (
	globalRegistry *Registry
	registryOnce   sync.Once
)

// Registry provides deterministic seed sequences.
type Registry struct {
	mu     sync.Mutex
	source *rand.PCG
}

// Init initializes the global registry with a master seed for repeatable simulations.
// If not called, the registry auto-initializes with a time-based seed on first use.
func Init(masterSeed uint64) {
	registryOnce.Do(func() {
		globalRegistry = newRegistry(masterSeed, masterSeed)
	})
}

// Next returns the next seed from the global registry.
// Auto-initializes with time-based seed if Init was not called.
func Next() (uint64, uint64) {
	registryOnce.Do(func() {
		now := uint64(time.Now().UnixNano())
		globalRegistry = newRegistry(now, now^0xdeadbeef)
	})

	return globalRegistry.next()
}

func newRegistry(seed1, seed2 uint64) *Registry {
	return &Registry{
		source: rand.NewPCG(seed1, seed2),
	}
}

func (r *Registry) next() (uint64, uint64) {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.source.Uint64(), r.source.Uint64()
}

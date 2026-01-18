package value

import "github.com/neox5/simv/transform"

// Publisher provides a subscription interface for typed values.
type Publisher[T any] interface {
	Subscribe() <-chan T
}

// ValueStats contains observable metrics for a Value.
type ValueStats[T any] struct {
	UpdateCount    uint64
	CurrentValue   T
	TransformCount int
}

// Value represents a readable, clonable, and resettable simulated value.
type Value[T any] interface {
	Value() T
	Clone() Value[T]
	WithTransforms(transforms ...transform.Transformation[T]) Value[T]
	SetState(T)
	SetUpdateHook(UpdateHook[T])
	Stats() ValueStats[T]
}

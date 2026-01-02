package source

// NumberSource generates values on demand.
type NumberSource[T any] interface {
	Next() T
}

// ConstSource always returns the same value.
type ConstSource[T any] struct {
	value T
}

// NewConstSource creates a source that always returns the given value.
func NewConstSource[T any](value T) *ConstSource[T] {
	return &ConstSource[T]{value: value}
}

// Next returns the constant value.
func (s *ConstSource[T]) Next() T {
	return s.value
}

package transform

// Transformation modifies a value.
type Transformation[T any] interface {
	Apply(value T) T
}

// AccumulateTransform adds each value to a running total.
// Requires T to support the + operator (int, int64, float64, etc.).
type AccumulateTransform[T interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}] struct {
	accumulated T
}

// NewAccumulateTransform creates a transform that accumulates values.
func NewAccumulateTransform[T interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}]() *AccumulateTransform[T] {
	return &AccumulateTransform[T]{}
}

// Apply adds the value to the accumulated total and returns the new total.
func (t *AccumulateTransform[T]) Apply(value T) T {
	t.accumulated += value
	return t.accumulated
}

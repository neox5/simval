package source

// SourceStats contains observable metrics for a Source.
type SourceStats struct {
	GenerationCount uint64
	SubscriberCount int
}

// Publisher provides a subscription interface for typed values.
type Publisher[T any] interface {
	Subscribe() <-chan T
	Stats() SourceStats
}

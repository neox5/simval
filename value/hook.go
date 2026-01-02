package value

// UpdateHook receives notifications during value update cycles.
type UpdateHook[T any] interface {
	OnInput(input T, state T)
	OnTransform(name string, input T, output T, state T)
	AfterUpdate(finalState T)
}

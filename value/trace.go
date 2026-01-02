package value

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

// TraceEvent captures a complete value update cycle.
type TraceEvent[T any] struct {
	Timestamp  time.Time
	Input      T                   // zero value if direct SetState
	Transforms []TransformTrace[T] // empty if direct SetState
	FinalState T
}

// TransformTrace captures a single transform application.
type TransformTrace[T any] struct {
	Name   string
	Input  T
	Output T
	State  T
}

// TraceHook implements UpdateHook to capture trace events.
type TraceHook[T any] struct {
	callback func(TraceEvent[T])

	// Accumulates data during update cycle
	mu         sync.Mutex
	timestamp  time.Time
	input      T
	transforms []TransformTrace[T]
}

// NewTraceHook creates a hook that calls callback with complete trace events.
func NewTraceHook[T any](callback func(TraceEvent[T])) *TraceHook[T] {
	return &TraceHook[T]{
		callback: callback,
	}
}

// NewDefaultTraceHook creates a TraceHook that prints formatted lines to stdout.
func NewDefaultTraceHook[T any]() *TraceHook[T] {
	return NewTraceHook(func(evt TraceEvent[T]) {
		fmt.Println(FormatTraceLine(evt))
	})
}

func (h *TraceHook[T]) OnInput(input T, state T) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.timestamp = time.Now()
	h.input = input
	h.transforms = h.transforms[:0] // reset
}

func (h *TraceHook[T]) OnTransform(name string, input T, output T, state T) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.transforms = append(h.transforms, TransformTrace[T]{
		Name:   name,
		Input:  input,
		Output: output,
		State:  state,
	})
}

func (h *TraceHook[T]) AfterUpdate(finalState T) {
	h.mu.Lock()

	// Check if this is a direct SetState call (no OnInput)
	hasInput := !h.timestamp.IsZero()

	event := TraceEvent[T]{
		Timestamp:  h.timestamp,
		Input:      h.input,
		Transforms: append([]TransformTrace[T](nil), h.transforms...),
		FinalState: finalState,
	}

	if !hasInput {
		// Direct SetState, set timestamp now
		event.Timestamp = time.Now()
		// Clear input and transforms for SetState-only events
		event.Input = *new(T) // zero value
		event.Transforms = nil
	}

	// Reset for next cycle
	h.timestamp = time.Time{}
	h.input = *new(T)               // Clear input
	h.transforms = h.transforms[:0] // Clear transforms

	h.mu.Unlock()

	h.callback(event)
}

// FormatTraceLine formats a trace event as a single pipe-separated line.
func FormatTraceLine[T any](evt TraceEvent[T]) string {
	timestamp := evt.Timestamp.Format("15:04:05.000")

	if len(evt.Transforms) > 0 {
		// Normal update: input | Transform(s:state) | output | ...
		var parts []string
		parts = append(parts, fmt.Sprintf("%v", evt.Input))

		for _, tr := range evt.Transforms {
			parts = append(parts, fmt.Sprintf("%s(s:%v)", tr.Name, tr.State))
			parts = append(parts, fmt.Sprintf("%v", tr.Output))
		}

		return fmt.Sprintf("[%s] %s", timestamp, strings.Join(parts, " | "))
	} else {
		// Direct SetState
		return fmt.Sprintf("[%s] SetState | %v", timestamp, evt.FinalState)
	}
}

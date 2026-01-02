package main

import (
	"fmt"
	"time"

	"github.com/neox5/simval/clock"
	"github.com/neox5/simval/source"
	"github.com/neox5/simval/transform"
	"github.com/neox5/simval/value"
)

func main() {
	// Create clock
	clk := clock.NewPeriodicClock(100 * time.Millisecond)

	// Create random source
	randomSrc := source.NewRandomIntSource(clk, 1, 10)

	// Create accumulated value
	accumulated := value.New(randomSrc, transform.NewAccumulate[int]())

	// Create reset-on-read value (cloned from accumulated)
	resetOnRead := value.NewResetOnRead(accumulated.Clone(), 0)

	// Enable tracing with default formatter
	resetOnRead.SetUpdateHook(value.NewDefaultTraceHook[int]())

	// Start clock
	clk.Start()
	defer clk.Stop()

	// Read and print every 500ms
	for range 10 {
		fmt.Printf(">>> ResetOnRead Value: %d\n",
			resetOnRead.Value(),
		)

		time.Sleep(500 * time.Millisecond)
	}
}

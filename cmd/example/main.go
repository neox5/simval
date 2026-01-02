package main

import (
	"fmt"
	"time"

	"github.com/neox5/simval/internal/clock"
	"github.com/neox5/simval/internal/source"
	"github.com/neox5/simval/internal/transform"
	"github.com/neox5/simval/internal/value"
)

func main() {
	// Create clock
	clk := clock.NewPeriodicClock(100 * time.Millisecond)

	// Create counter: constant source of 1 + accumulate transform
	counter := value.New(
		clk,
		source.NewConstSource(1),
		transform.NewAccumulateTransform[int](),
	)
	defer counter.Stop()

	// Start clock
	clk.Start()
	defer clk.Stop()

	// Read and print every 500ms
	for range 10 {
		fmt.Printf("Counter value: %d\n", counter.Value())
		time.Sleep(500 * time.Millisecond)
	}
}

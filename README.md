# simV

A Go library for simulating time-varying values in testing environments.

## Overview

simV provides thread-safe entities representing different types of time-varying values, enabling developers to create realistic simulated metrics and time-series data for system validation.

**Use cases:**

- Testing OTEL/Prometheus metric pipelines
- Validating observability systems
- Simulating realistic load patterns
- Mocking external metric sources during development
- Creating observable test fixtures with temporal behavior

## Installation

```bash
go get github.com/neox5/simv@latest
```

**Requirements:** Go 1.25+

## Quick Start

```go
package main

import (
    "fmt"
    "time"

    "github.com/neox5/simv/clock"
    "github.com/neox5/simv/seed"
    "github.com/neox5/simv/source"
    "github.com/neox5/simv/transform"
    "github.com/neox5/simv/value"
)

func main() {
    // Optional: Initialize seed for repeatable simulations
    seed.Init(12345)

    // Create clock that ticks every 100ms
    clk := clock.NewPeriodicClock(100 * time.Millisecond)

    // Create source that generates random integers [1, 10]
    src := source.NewRandomIntSource(clk, 1, 10)

    // Create value that accumulates incoming values
    val := value.New(src, transform.NewAccumulate[int]())

    // Wrap with reset-on-read behavior
    resetVal := value.NewResetOnRead(val.Clone(), 0)

    // Enable trace output
    resetVal.SetUpdateHook(value.NewDefaultTraceHook[int]())

    // Start the clock
    clk.Start()
    defer clk.Stop()

    // Read current value
    current := resetVal.Value()
    fmt.Printf("Current value: %d\n", current)

    // Access metrics without side effects
    stats := resetVal.Stats()
    fmt.Printf("Updates: %d, Current: %d\n",
        stats.UpdateCount,
        stats.CurrentValue,
    )
}
```

## Architecture

simV uses a pipeline architecture:

**Clock** → **Source** → **Transform** → **Value**

- **Clock**: Generates timing signals at fixed intervals
- **Source**: Produces values on each clock tick
- **Transform**: Modifies values (accumulate, average, etc.)
- **Value**: Manages state and provides thread-safe access

## Core Concepts

### Seed

Control repeatability of random value generation.

```go
// Repeatable simulations - same seed produces identical sequences
seed.Init(12345)

// Non-repeatable (default) - auto-initializes with time-based seed
// No need to call Init if repeatability is not required
```

### Clock

Provides timing signals for value generation.

```go
clk := clock.NewPeriodicClock(100 * time.Millisecond)
clk.Start()
defer clk.Stop()

// Access metrics
stats := clk.Stats()
fmt.Printf("Ticks: %d, Running: %v\n", stats.TickCount, stats.IsRunning)
```

### Source

Generates values driven by clock ticks.

```go
// Constant value
constSrc := source.NewConstSource(clk, 42)

// Random integers
randomSrc := source.NewRandomIntSource(clk, 1, 100)

// Access metrics
stats := randomSrc.Stats()
fmt.Printf("Generated: %d, Subscribers: %d\n",
    stats.GenerationCount,
    stats.SubscriberCount,
)
```

### Transform

Applies operations to incoming values.

```go
// Running total
accumulated := value.New(src, transform.NewAccumulate[int]())
```

### Value

Thread-safe value management with optional behaviors.

```go
// Standard value
val := value.New(src)

// Reset on each read
resetVal := value.NewResetOnRead(val.Clone(), 0)

// Enable tracing
val.SetUpdateHook(value.NewDefaultTraceHook[int]())

// Access metrics without side effects
stats := val.Stats()
fmt.Printf("Updates: %d, Current: %d, Transforms: %d\n",
    stats.UpdateCount,
    stats.CurrentValue,
    stats.TransformCount,
)
```

## Observability

### Metrics

All components expose metrics via `Stats()` methods for integration with Prometheus/OTEL:

```go
// Clock metrics
clockStats := clk.Stats()
// - TickCount: total ticks generated
// - IsRunning: current operational state
// - Interval: tick rate

// Source metrics
sourceStats := src.Stats()
// - GenerationCount: total values produced
// - SubscriberCount: active subscriptions

// Value metrics
valueStats := val.Stats()
// - UpdateCount: total updates received
// - CurrentValue: current value without side effects
// - TransformCount: number of transforms in chain
```

**Prometheus example:**

```go
tickCounter := prometheus.NewCounter(...)
stats := clk.Stats()
tickCounter.Add(float64(stats.TickCount))
```

### Tracing

Enable trace output to observe value flow through the pipeline:

```go
val.SetUpdateHook(value.NewDefaultTraceHook[int]())
// Output: [15:04:05.000] 7 | Accumulate(s:42) | 49
```

## Features

- Generic type support
- Thread-safe concurrent access
- Subscription-based value distribution
- Composable transform pipeline
- Observable update cycles via hooks
- Metrics exposure for monitoring systems
- Repeatable simulations via seed control
- Zero external dependencies

## Examples

See `cmd/example/main.go` for complete working example.

## License

MIT License - see LICENSE file for details.

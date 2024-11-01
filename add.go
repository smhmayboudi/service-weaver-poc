package main

import (
	"context"

	"github.com/ServiceWeaver/weaver"
	"github.com/ServiceWeaver/weaver/metrics"
)

type (
	PortAdd interface {
		Add(context.Context, int, int) (int, error)
	}

	add struct {
		weaver.Implements[PortAdd]
	}
)

var (
	addCount = metrics.NewCounter(
		"add_count",
		"The number of times Add.Add has been called",
	)
	addConcurrent = metrics.NewGauge(
		"add_concurrent",
		"The number of concurrent Add.Add calls",
	)
	addSum = metrics.NewHistogram(
		"add_sum",
		"The sums returned by Add.Add",
		[]float64{1, 10, 100, 1000, 10000},
	)
)

func (a *add) Add(ctx context.Context, x, y int) (int, error) {
	logger := a.Logger(ctx).With("code.function", "Add")
	logger.Info("")

	addCount.Add(1.0)
	addConcurrent.Add(1.0)
	defer addConcurrent.Sub(1.0)

	out := x + y
	logger.Debug("out: %v", out)
	addSum.Put(float64(out))

	return out, nil
}

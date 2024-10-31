package main

import (
	"context"

	"github.com/ServiceWeaver/weaver"
	"github.com/ServiceWeaver/weaver/metrics"
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

type PortAdd interface {
	Add(context.Context, int, int) (int, error)
}

type add struct {
	weaver.Implements[PortAdd]
}

func (*add) Add(_ context.Context, x, y int) (int, error) {
	addCount.Add(1.0)
	addConcurrent.Add(1.0)
	defer addConcurrent.Sub(1.0)
	addSum.Put(float64(x + y))
	return x + y, nil
}

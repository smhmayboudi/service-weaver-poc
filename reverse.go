package main

import (
	"context"

	"github.com/ServiceWeaver/weaver"
	"github.com/ServiceWeaver/weaver/metadata"
	"github.com/ServiceWeaver/weaver/metrics"
)

type (
	// PortReverseOptions component.
	PortReverseOptions struct {
		Greeting string
	}

	// PortReverse component.
	PortReverse interface {
		Reverse(context.Context, string) (string, error)
	}

	// Implementation of the PortReverse component.
	reverse struct {
		weaver.Implements[PortReverse]
		weaver.WithConfig[PortReverseOptions]
	}

	labels struct {
		Parity string // "hello" or else
	}
)

var (
	halveCounts = metrics.NewCounterMap[labels](
		"reverse_label_count",
		"The number of values that have been reversed",
	)
	helloCount = halveCounts.Get(labels{"hello"})
	elseCount  = halveCounts.Get(labels{"else"})

	reverseCount = metrics.NewCounter(
		"reverse_count",
		"The number of times PortReverse.Reverse has been called",
	)
	reverseConcurrent = metrics.NewGauge(
		"reverse_concurrent",
		"The number of concurrent PortReverse.Reverse calls",
	)
	reverseSum = metrics.NewHistogram(
		"reverse_sum",
		"The sums returned by PortReverse.Reverse",
		[]float64{1, 10, 100, 1000, 10000},
	)
)

func (r *reverse) Reverse(ctx context.Context, s string) (string, error) {
	reverseCount.Add(1.0)
	reverseConcurrent.Add(1.0)
	defer reverseConcurrent.Sub(1.0)

	if s == "hello" {
		helloCount.Add(1)
	} else {
		elseCount.Add(1)
	}

	var defaultGreeting = ""
	meta, ok := metadata.FromContext(ctx)
	if ok {
		defaultGreeting = meta["default_greeting"]
	}

	logger := r.Logger(ctx).With("foo", "bar")
	logger.Debug("A debug log.")
	logger.Info("An info log.")
	logger.Error("An error log.")

	runes := []rune(s)
	n := len(runes)
	for i := 0; i < n/2; i++ {
		runes[i], runes[n-i-1] = runes[n-i-1], runes[i]
	}

	greeting := r.Config().Greeting
	if greeting == "" {
		greeting = defaultGreeting
	}

	return greeting + string(runes), nil
}

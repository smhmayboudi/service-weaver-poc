package main

import (
	"context"
	"strconv"

	"github.com/ServiceWeaver/weaver"
	"github.com/ServiceWeaver/weaver/metadata"
	"github.com/ServiceWeaver/weaver/metrics"
)

type (
	PortReverse interface {
		Reverse(context.Context, string) (string, error)
	}

	reverse struct {
		add weaver.Ref[PortAdd]
		weaver.Implements[PortReverse]
		weaver.WithConfig[reverseOption]
		word weaver.Ref[PortWord]
	}

	reverseOption struct {
		Greeting string
	}
)

var (
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

func (r *reverse) Reverse(ctx context.Context, str string) (string, error) {
	logger := r.Logger(ctx).With("code.function", "Reverse").With("str", str)
	logger.Info("")

	reverseCount.Add(1.0)
	reverseConcurrent.Add(1.0)
	defer reverseConcurrent.Sub(1.0)

	r.word.Get().Parse(ctx, str)

	var defaultGreeting = ""
	logger.Debug("", "defaultGreeting", defaultGreeting)
	meta, ok := metadata.FromContext(ctx)
	if ok {
		logger.Debug("", "meta", meta)
		defaultGreeting = meta["default_greeting"]
		logger.Debug("", "defaultGreeting", defaultGreeting)
	}

	greeting := r.Config().Greeting
	logger.Debug("", "greeting", greeting)
	if greeting == "" {
		logger.Debug("inside if")
		greeting = defaultGreeting
		logger.Debug("", "greeting", greeting)
	}

	runes := []rune(str)
	n := len(runes)
	for i := 0; i < n/2; i++ {
		runes[i], runes[n-i-1] = runes[n-i-1], runes[i]
	}
	res, _ := r.add.Get().Add(ctx, 1, 2)
	val := strconv.Itoa(res) + greeting + string(runes)
	logger.Debug("", "val", val)
	reverseSum.Put(1)

	return val, nil
}

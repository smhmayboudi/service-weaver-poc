package main

import (
	"context"
	"strconv"

	lru "github.com/hashicorp/golang-lru/v2"

	"github.com/ServiceWeaver/weaver"
	"github.com/ServiceWeaver/weaver/metadata"
	"github.com/ServiceWeaver/weaver/metrics"
)

// The size of a factorer's LRU cache.
const cacheSize = 100

type (
	// PortReverse component.
	PortReverse interface {
		Reverse(context.Context, string) (string, error)
	}

	// reverseOptions component.
	reverseOptions struct {
		Greeting string
	}

	// router component for cache.
	router struct{}

	// Implementation of the PortReverse component.
	reverse struct {
		add   weaver.Ref[PortAdd]
		cache *lru.Cache[string, string]
		weaver.Implements[PortReverse]
		weaver.WithConfig[reverseOptions]
		weaver.WithRouter[router]
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

func (r *reverse) Init(_ context.Context) error {
	cache, err := lru.New[string, string](cacheSize)
	r.cache = cache
	return err
}

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

	res, _ := r.add.Get().Add(ctx, 1, 2)

	out := strconv.Itoa(res) + greeting + string(runes)

	r.cache.Add(s, out)

	return out, nil
}

func (r *router) Reverse(ctx context.Context, s string) string {
	return s
}

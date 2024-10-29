package main

import (
	"context"

	"github.com/ServiceWeaver/weaver"
	"github.com/ServiceWeaver/weaver/metadata"
)

type PortReverseOptions struct {
	Greeting string
}

// PortReverse component.
type PortReverse interface {
	Reverse(context.Context, string) (string, error)
}

// Implementation of the Reverser component.
type reverse struct {
	weaver.Implements[PortReverse]
	weaver.WithConfig[PortReverseOptions]
}

func (r *reverse) Reverse(ctx context.Context, s string) (string, error) {
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

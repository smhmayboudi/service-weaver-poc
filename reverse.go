package main

import (
	"context"

	"github.com/ServiceWeaver/weaver"
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

func (r *reverse) Reverse(_ context.Context, s string) (string, error) {
	runes := []rune(s)
	n := len(runes)
	for i := 0; i < n/2; i++ {
		runes[i], runes[n-i-1] = runes[n-i-1], runes[i]
	}
	greeting := r.Config().Greeting
	if greeting == "" {
		greeting = "!!!"
	}

	return greeting + string(runes), nil
}

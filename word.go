package main

import (
	"context"

	"github.com/ServiceWeaver/weaver"
	"github.com/ServiceWeaver/weaver/metrics"
)

type (
	PortWord interface {
		Parse(context.Context, string) error
	}

	word struct {
		weaver.Implements[PortWord]
	}

	wordLabel struct {
		WordParity string // "hello" or "notHello"
	}
)

var (
	_         PortWord = (*word)(nil)
	wordCount          = metrics.NewCounterMap[wordLabel](
		"word_count",
		"The number of words that have been hello",
	)
	helloCount    = wordCount.Get(wordLabel{"hello"})
	notHelloCount = wordCount.Get(wordLabel{"notHello"})
)

func (w *word) Parse(ctx context.Context, value string) error {
	logger := w.Logger(ctx).With("code.function", "Parse").With("value", value)
	logger.Info("")

	if value == "hello" {
		logger.Debug("inside if")
		helloCount.Add(1)
	} else {
		logger.Debug("inside ifelse")
		notHelloCount.Add(1)
	}
	return nil
}

package main

import (
	"context"
	"sync"

	"github.com/ServiceWeaver/weaver"
)

// The size of a factorer's LRU cache.
// const cacheSize = 100

type (
	PortCache interface {
		Get(ctx context.Context, key string) (string, error)
		Put(ctx context.Context, key, value string) error
	}

	cache struct {
		// cache *lru.Cache[string, string]
		data map[string]string
		mu   sync.Mutex
		weaver.Implements[PortCache]
		weaver.WithRouter[cacheRouter]
	}
)

var (
	_ PortCache = (*cache)(nil)
	// _ weaver.NotRetriable = PortCache.Put
)

func (c *cache) Init(ctx context.Context) error {
	logger := c.Logger(ctx).With("code.function", "Init")
	logger.Info("")

	c.data = map[string]string{}
	// cache, err := lru.New[string, string](cacheSize)
	// r.cache = cache
	// return err
	return nil
}

func (c *cache) Get(ctx context.Context, key string) (string, error) {
	logger := c.Logger(ctx).With("code.function", "Get")
	logger.Info("")

	c.mu.Lock()
	defer c.mu.Unlock()
	return c.data[key], nil
}

func (c *cache) Put(ctx context.Context, key, value string) error {
	logger := c.Logger(ctx).With("code.function", "Put")
	logger.Info("")

	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = value
	return nil
}

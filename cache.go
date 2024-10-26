package main

import (
	"context"
	"sync"

	"github.com/ServiceWeaver/weaver"
)

type PortCache interface {
	Append(ctx context.Context, key, value string) error
	Get(ctx context.Context, key string) (string, error)
	Put(ctx context.Context, key, value string) error
}

type cache struct {
	weaver.Implements[PortCache]
	mu   sync.Mutex
	data map[string]string
}

var _ PortCache = (*cache)(nil)

var _ weaver.NotRetriable = PortCache.Append

func (c *cache) Append(_ context.Context, key, value string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = value
	return nil
}

func (c *cache) Get(_ context.Context, key string) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.data[key], nil
}

func (c *cache) Put(_ context.Context, key, value string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = value
	return nil
}

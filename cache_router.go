package main

import (
	"context"
	"fmt"
)

type cacheRouter struct{}

func (cacheRouter) Get(_ context.Context, key string) string {
	fmt.Printf("cacheRouter => Get")
	return key
}

func (cacheRouter) Put(_ context.Context, key, value string) string {
	fmt.Printf("cacheRouter => Put")
	return key
}

package main

import (
	"context"
	"fmt"
)

type cacheRouter struct{}

func (cacheRouter) Get(_ context.Context, key string) string {
	fmt.Println("cacheRouter => Get")
	return key
}

func (cacheRouter) Put(_ context.Context, key, value string) string {
	fmt.Println("cacheRouter => Put")
	return key
}

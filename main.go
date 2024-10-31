package main

import (
	"context"
	"flag"
	"log"

	"github.com/ServiceWeaver/weaver"
)

//go:generate weaver generate

func main() {
	flag.Parse()
	ctx := context.Background()
	if err := weaver.Run(ctx, serve); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"context"
	"strconv"
	"testing"

	"github.com/ServiceWeaver/weaver/weavertest"
)

type fakeAdd struct{}

var _ PortAdd = (*fakeAdd)(nil)

func (add *fakeAdd) Add(_ context.Context, a, b int) (int, error) {
	return 0, nil
}

func TestReverse(t *testing.T) {
	for _, runner := range weavertest.AllRunners() {
		fake := &fakeAdd{}
		runner.Fakes = append(runner.Fakes, weavertest.Fake[PortAdd](fake))
		runner.Test(t, func(t *testing.T, reverse PortReverse) {
			ctx := context.Background()
			got, err := reverse.Reverse(ctx, "12")
			if err != nil {
				t.Fatal(err)
			}
			if want := "021"; got != want {
				t.Fatalf("got %q, want %q", got, want)
			}
		})
	}
}

func TestReverse2(t *testing.T) {
	for _, runner := range weavertest.AllRunners() {
		runner.Test(t, func(t *testing.T, reverse PortReverse, add PortAdd) {
			ctx := context.Background()
			got, err := reverse.Reverse(ctx, "12")
			if err != nil {
				t.Fatal(err)
			}
			a, _ := add.Add(ctx, 1, 2)
			if want := strconv.Itoa(a) + "21"; got != want {
				t.Fatalf("got %q, want %q", got, want)
			}
		})
	}
}

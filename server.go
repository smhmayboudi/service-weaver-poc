package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/ServiceWeaver/weaver"
	"github.com/ServiceWeaver/weaver/metadata"
)

type server struct {
	weaver.Implements[weaver.Main]
	reverse weaver.Ref[PortReverse]
	cache   weaver.Ref[PortCache]
	hello   weaver.Listener
}

func serve(ctx context.Context, s *server) error {
	http.Handle("/", weaver.InstrumentHandlerFunc("/", s.handleFactors))
	s.Logger(ctx).Info("factors server running", "addr", s.hello)
	return http.Serve(s.hello, nil)
}

// handleFactors handles the /?x=<number> endpoint.
func (s *server) handleFactors(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "World"
	}

	// Call the cache.Get method.
	value, err := s.cache.Get().Get(r.Context(), "key")
	if errors.Is(err, weaver.RemoteCallError) {
		// cache.Get did not execute properly.
		fmt.Printf("1 Error %v\n", err)
	} else if err != nil {
		// cache.Get executed properly, but returned an error.
		fmt.Printf("2 Error %v\n", err)
	} else {
		fmt.Printf("3 OK\n")
		// cache.Get executed properly and did not return an error.
	}

	ctx := context.Background()
	ctx = metadata.NewContext(ctx, map[string]string{"default_greeting": "nothing"})

	reversed, err := s.reverse.Get().Reverse(ctx, name+value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Hello, %s!\n", reversed)

}

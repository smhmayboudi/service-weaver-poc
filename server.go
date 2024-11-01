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
	logger := s.Logger(ctx).With("code.function", "serve")
	logger.Info("")

	http.Handle("/", weaver.InstrumentHandlerFunc("/", s.handleFactors))
	logger.Debug("factors server running", "addr", s.hello)

	return http.Serve(s.hello, nil)
}

// handleFactors handles the /?x=<number> endpoint.
func (s *server) handleFactors(w http.ResponseWriter, r *http.Request) {
	logger := s.Logger(r.Context()).With("code.function", "handleFactors")
	logger.Info("")

	name := r.URL.Query().Get("name")
	if name == "" {
		name = "World"
	}

	// Call the cache.Get method.
	value, err := s.cache.Get().Get(r.Context(), name)
	if errors.Is(err, weaver.RemoteCallError) {
		// cache.Get did not execute properly.
		logger.Error("1 Error %v", err)
	} else if err != nil {
		// cache.Get executed properly, but returned an error.
		logger.Error("2 Error %v", err)
	} else {
		// cache.Get executed properly and did not return an error.
		logger.Debug("3 OK %v", value)
		if value == "" {
			logger.Debug("inside if")
			ctx := metadata.NewContext(r.Context(), map[string]string{"default_greeting": "nothing"})
			reversed, err := s.reverse.Get().Reverse(ctx, name+value)
			if err != nil {
				logger.Error("Error %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			value = reversed
		}
		fmt.Fprintf(w, "Hello, %s!\n", value)
	}
}

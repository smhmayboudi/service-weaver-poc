package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/ServiceWeaver/weaver"
	"github.com/ServiceWeaver/weaver/metadata"
)

func main() {
	if err := weaver.Run(context.Background(), serve); err != nil {
		log.Fatal(err)
	}
}

type app struct {
	weaver.Implements[weaver.Main]
	reverse weaver.Ref[PortReverse]
	cache   weaver.Ref[PortCache]
	hello   weaver.Listener
}

func serve(ctx context.Context, app *app) error {
	// The hello listener will listen on a random port chosen by the operating
	// system. This behavior can be changed in the config file.
	fmt.Printf("hello listener available on %v\n", app.hello)

	// Serve the /hello endpoint.
	http.Handle("/hello", weaver.InstrumentHandlerFunc("hello", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		if name == "" {
			name = "World"
		}

		// Call the cache.Get method.
		value, err := app.cache.Get().Get(ctx, "key")
		if errors.Is(err, weaver.RemoteCallError) {
			// cache.Get did not execute properly.
			fmt.Fprint(w, "1", err)
		} else if err != nil {
			// cache.Get executed properly, but returned an error.
			fmt.Fprint(w, "2", err)
		} else {
			fmt.Fprint(w, "3", "ok")
			// cache.Get executed properly and did not return an error.
		}

		ctx := context.Background()
		ctx = metadata.NewContext(ctx, map[string]string{"default_greeting": "nothing"})

		reversed, err := app.reverse.Get().Reverse(ctx, name+value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "Hello, %s!\n", reversed)
	}))

	return http.Serve(app.hello, nil)
}

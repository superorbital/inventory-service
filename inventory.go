// This is an example of implementing the Item Store from the OpenAPI documentation
// found at:
// https://github.com/superorbital/inventory-service/blob/master/inventory-service.yaml

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	middleware "github.com/deepmap/oapi-codegen/pkg/chi-middleware"
	"github.com/go-chi/chi/v5"
	api "github.com/superorbital/inventory-service/api"
)

func main() {
	var port = flag.Int("port", 8080, "Port for test HTTP server")
	flag.Parse()

	swagger, err := api.GetSwagger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading swagger spec\n: %s", err)
		os.Exit(1)
	}

	// Clear out the servers array in the swagger spec, that skips validating
	// that server names match. We don't know how this thing will be run.
	swagger.Servers = nil

	// Create an instance of our handler which satisfies the generated interface
	itemStore := api.NewInventory()

	// This is how you set up a basic chi router
	r := chi.NewRouter()

	// Use our validation middleware to check all requests against the
	// OpenAPI schema.
	r.Use(middleware.OapiRequestValidator(swagger))

	// We now register our itemStore above as the handler for the interface
	api.HandlerFromMux(itemStore, r)

	s := &http.Server{
		Handler:           r,
		Addr:              fmt.Sprintf("0.0.0.0:%d", *port),
		ReadHeaderTimeout: 60 * time.Second,
	}

	// And we serve HTTP until the world ends.
	log.Fatal(s.ListenAndServe())
}

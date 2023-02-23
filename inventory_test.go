package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	middleware "github.com/deepmap/oapi-codegen/pkg/chi-middleware"
	"github.com/deepmap/oapi-codegen/pkg/testutil"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	api "github.com/superorbital/inventory-service/api"
)

func doGet(t *testing.T, mux *chi.Mux, url string) *httptest.ResponseRecorder {
	response := testutil.NewRequest().Get(url).WithAcceptJson().GoWithHTTPHandler(t, mux)
	return response.Recorder
}

func TestInventory(t *testing.T) {
	var err error

	// Get the swagger description of our API
	swagger, err := api.GetSwagger()
	require.NoError(t, err)

	// Clear out the servers array in the swagger spec, that skips validating
	// that server names match. We don't know how this thing will be run.
	swagger.Servers = nil

	// This is how you set up a basic chi router
	r := chi.NewRouter()

	// Use our validation middleware to check all requests against the
	// OpenAPI schema.
	r.Use(middleware.OapiRequestValidator(swagger))

	store := api.NewInventory()
	api.HandlerFromMux(store, r)

	t.Run("Add item", func(t *testing.T) {
		tag := "TagOfSpot"
		newItem := api.NewItem{
			Name: "Spot",
			Tag:  &tag,
		}

		rr := testutil.NewRequest().Post("/items").WithJsonBody(newItem).GoWithHTTPHandler(t, r).Recorder
		assert.Equal(t, http.StatusCreated, rr.Code)

		var resultItem api.Item
		err = json.NewDecoder(rr.Body).Decode(&resultItem)
		assert.NoError(t, err, "error unmarshalling response")
		assert.Equal(t, newItem.Name, resultItem.Name)
		assert.Equal(t, *newItem.Tag, *resultItem.Tag)
	})

	t.Run("Find item by ID", func(t *testing.T) {
		item := api.Item{
			Id: 100,
		}

		store.Items[item.Id] = item
		rr := doGet(t, r, fmt.Sprintf("/items/%d", item.Id))

		var resultItem api.Item
		err = json.NewDecoder(rr.Body).Decode(&resultItem)
		assert.NoError(t, err, "error getting item")
		assert.Equal(t, item, resultItem)
	})

	t.Run("Item not found", func(t *testing.T) {
		rr := doGet(t, r, "/items/27179095781")
		assert.Equal(t, http.StatusNotFound, rr.Code)

		var itemError api.Error
		err = json.NewDecoder(rr.Body).Decode(&itemError)
		assert.NoError(t, err, "error getting response", err)
		assert.Equal(t, int32(http.StatusNotFound), itemError.Code)
	})

	t.Run("List all items", func(t *testing.T) {
		store.Items = map[int64]api.Item{
			1: {},
			2: {},
		}

		// Now, list all items, we should have two
		rr := doGet(t, r, "/items")
		assert.Equal(t, http.StatusOK, rr.Code)

		var itemList []api.Item
		err = json.NewDecoder(rr.Body).Decode(&itemList)
		assert.NoError(t, err, "error getting response", err)
		assert.Equal(t, 2, len(itemList))
	})

	t.Run("Filter items by tag", func(t *testing.T) {
		tag := "TagOfFido"

		store.Items = map[int64]api.Item{
			1: {
				Tag: &tag,
			},
			2: {},
		}

		// Filter items by tag, we should have 1
		rr := doGet(t, r, "/items?tags=TagOfFido")
		assert.Equal(t, http.StatusOK, rr.Code)

		var itemList []api.Item
		err = json.NewDecoder(rr.Body).Decode(&itemList)
		assert.NoError(t, err, "error getting response", err)
		assert.Equal(t, 1, len(itemList))
	})

	t.Run("Filter items by tag", func(t *testing.T) {
		store.Items = map[int64]api.Item{
			1: {},
			2: {},
		}

		// Filter items by non-existent tag, we should have 0
		rr := doGet(t, r, "/items?tags=NotExists")
		assert.Equal(t, http.StatusOK, rr.Code)

		var itemList []api.Item
		err = json.NewDecoder(rr.Body).Decode(&itemList)
		assert.NoError(t, err, "error getting response", err)
		assert.Equal(t, 0, len(itemList))
	})

	t.Run("Delete items", func(t *testing.T) {
		store.Items = map[int64]api.Item{
			1: {},
			2: {},
		}

		// Let's delete non-existent item
		rr := testutil.NewRequest().Delete("/items/7").GoWithHTTPHandler(t, r).Recorder
		assert.Equal(t, http.StatusNotFound, rr.Code)

		var itemError api.Error
		err = json.NewDecoder(rr.Body).Decode(&itemError)
		assert.NoError(t, err, "error unmarshalling ItemError")
		assert.Equal(t, int32(http.StatusNotFound), itemError.Code)

		// Now, delete both real items
		rr = testutil.NewRequest().Delete("/items/1").GoWithHTTPHandler(t, r).Recorder
		assert.Equal(t, http.StatusNoContent, rr.Code)

		rr = testutil.NewRequest().Delete("/items/2").GoWithHTTPHandler(t, r).Recorder
		assert.Equal(t, http.StatusNoContent, rr.Code)

		// Should have no items left.
		var itemList []api.Item
		rr = doGet(t, r, "/items")
		assert.Equal(t, http.StatusOK, rr.Code)
		err = json.NewDecoder(rr.Body).Decode(&itemList)
		assert.NoError(t, err, "error getting response", err)
		assert.Equal(t, 0, len(itemList))
	})
}

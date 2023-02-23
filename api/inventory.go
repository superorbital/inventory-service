//go:generate oapi-codegen --config=cfg.yaml ../inventory-openapi.yaml

package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type Inventory struct {
	Items  map[int64]Item
	NextId int64
	Lock   sync.Mutex
}

// Make sure we conform to ServerInterface

var _ ServerInterface = (*Inventory)(nil)

func NewInventory() *Inventory {
	return &Inventory{
		Items:  make(map[int64]Item),
		NextId: 1000,
	}
}

// This function wraps sending of an error in the Error format, and
// handling the failure to marshal that.
func sendInventoryError(w http.ResponseWriter, code int, message string) {
	itemErr := Error{
		Code:    int32(code),
		Message: message,
	}
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(itemErr)
}

// FindItems implements all the handlers in the ServerInterface
func (p *Inventory) FindItems(w http.ResponseWriter, r *http.Request, params FindItemsParams) {
	p.Lock.Lock()
	defer p.Lock.Unlock()

	var result []Item

	for _, item := range p.Items {
		if params.Tags != nil {
			// If we have tags,  filter items by tag
			for _, t := range *params.Tags {
				if item.Tag != nil && (*item.Tag == t) {
					result = append(result, item)
				}
			}
		} else {
			// Add all items if we're not filtering
			result = append(result, item)
		}

		if params.Limit != nil {
			l := int(*params.Limit)
			if len(result) >= l {
				// We're at the limit
				break
			}
		}
	}

	if result == nil {
		result = []Item{}
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func (p *Inventory) AddItem(w http.ResponseWriter, r *http.Request) {
	// We expect a NewItem object in the request body.
	var newItem NewItem
	if err := json.NewDecoder(r.Body).Decode(&newItem); err != nil {
		sendInventoryError(w, http.StatusBadRequest, "Invalid format for NewItem")
		return
	}

	// We now have a item, let's add it to our "database".

	// We're always asynchronous, so lock unsafe operations below
	p.Lock.Lock()
	defer p.Lock.Unlock()

	// We handle items, not NewItems, which have an additional ID field
	var item Item
	item.Name = newItem.Name
	item.Tag = newItem.Tag
	item.Id = p.NextId
	p.NextId = p.NextId + 1

	// Insert into map
	p.Items[item.Id] = item

	// Now, we have to return the NewItem
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}

func (p *Inventory) FindItemById(w http.ResponseWriter, r *http.Request, id int64) {
	p.Lock.Lock()
	defer p.Lock.Unlock()

	item, found := p.Items[id]
	if !found {
		sendInventoryError(w, http.StatusNotFound, fmt.Sprintf("Could not find item with ID %d", id))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(item)
}

func (p *Inventory) UpdateItem(w http.ResponseWriter, r *http.Request, id int64) {
	p.Lock.Lock()
	defer p.Lock.Unlock()

	_, found := p.Items[id]
	if !found {
		sendInventoryError(w, http.StatusNotFound, fmt.Sprintf("Could not find item with ID %d", id))
		return
	}

	var newItem NewItem
	if err := json.NewDecoder(r.Body).Decode(&newItem); err != nil {
		sendInventoryError(w, http.StatusBadRequest, "Invalid format for Item")
		return
	}

	item := Item{
		Id:   id,
		Name: newItem.Name,
		Tag:  newItem.Tag,
	}

	p.Items[id] = item

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(item)
}

func (p *Inventory) DeleteItem(w http.ResponseWriter, r *http.Request, id int64) {
	p.Lock.Lock()
	defer p.Lock.Unlock()

	_, found := p.Items[id]
	if !found {
		sendInventoryError(w, http.StatusNotFound, fmt.Sprintf("Could not find item with ID %d", id))
		return
	}
	delete(p.Items, id)

	w.WriteHeader(http.StatusNoContent)
}

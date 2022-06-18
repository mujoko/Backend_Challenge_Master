package handlers

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"go-inventory/errors"
	"go-inventory/objects"
	"go-inventory/store"

	"github.com/gorilla/mux"
)

// IStockHandler is implement all the handlers
type IStockHandler interface {
	Get(w http.ResponseWriter, r *http.Request)
	List(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	UpdateDetails(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	store store.IStockStore
}

// NewEventHandler return current IStockHandler implementation
func NewEventHandler(store store.IStockStore) IStockHandler {
	return &handler{store: store}
}

func (h *handler) Get(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		WriteError(w, errors.ErrValidStockIDIsRequired)
		return
	}
	fmt.Println(`id := `, id)

	evt, err := h.store.Get(r.Context(), &objects.GetRequest{ID: id})
	if err != nil {
		WriteError(w, err)
		return
	}
	WriteResponse(w, &objects.StockResponseWrapper{Stock: evt})
}

func (h *handler) List(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	// after
	after := values.Get("after")
	// name
	name := values.Get("name")
	// limit
	limit, err := IntFromString(w, values.Get("limit"))
	if err != nil {
		return
	}
	// list events
	list, err := h.store.List(r.Context(), &objects.ListRequest{
		Limit: limit,
		After: after,
		Name:  name,
	})
	if err != nil {
		WriteError(w, err)
		return
	}
	WriteResponse(w, &objects.StockResponseWrapper{Stocks: list})
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		WriteError(w, errors.ErrUnprocessableEntity)
		return
	}
	evt := &objects.Stock{}
	if Unmarshal(w, data, evt) != nil {
		return
	}
	if evt.Availability < 0 {
		WriteError(w, errors.ErrValidAvailibiltyIsRequired)
		return
	}
	if evt.Price <= 0 {
		WriteError(w, errors.ErrValidPriceIsRequired)
		return
	}
	if err = h.store.Create(r.Context(), &objects.CreateRequest{Stock: evt}); err != nil {
		WriteError(w, err)
		return
	}
	WriteResponse(w, &objects.StockResponseWrapper{Stock: evt})
}

func (h *handler) UpdateDetails(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		WriteError(w, errors.ErrUnprocessableEntity)
		return
	}
	req := &objects.UpdateDetailsRequest{}
	if Unmarshal(w, data, req) != nil {
		return
	}

	// check if event exist
	if _, err := h.store.Get(r.Context(), &objects.GetRequest{ID: req.ID}); err != nil {
		WriteError(w, err)
		return
	}

	if err = h.store.UpdateDetails(r.Context(), req); err != nil {
		WriteError(w, err)
		return
	}
	WriteResponse(w, &objects.StockResponseWrapper{})
}

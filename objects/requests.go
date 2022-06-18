package objects

import (
	"encoding/json"
	"net/http"
)

// MaxListLimit maximum listting
const MaxListLimit = 200

// GetRequest for retrieving single Stock
type GetRequest struct {
	ID string `json:"id"`
}

// ListRequest for retrieving list of Stocks
type ListRequest struct {
	Limit int    `json:"limit"`
	After string `json:"after"`
	// optional name matching
	Name string `json:"name"`
}

// CreateRequest for creating a new Stock
type CreateRequest struct {
	Stock *Stock `json:"Stock"`
}

// UpdateDetailsRequest to update existing Stock
type UpdateDetailsRequest struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	Price        float64 `json:"price"`
	Availability int     `json:"availability"`
	IsActive     bool    `json:"is_active"`
}

// CancelRequest to cancel an Stock
type CancelRequest struct {
	ID string `json:"id"`
}

// DeleteRequest to delete an Stock
type DeleteRequest struct {
	ID string `json:"id"`
}

// StockResponseWrapper reponse of any Stock request
type StockResponseWrapper struct {
	Stock  *Stock   `json:"Stock,omitempty"`
	Stocks []*Stock `json:"Stocks,omitempty"`
	Code   int      `json:"-"`
}

// JSON convert StockResponseWrapper in json
func (e *StockResponseWrapper) JSON() []byte {
	if e == nil {
		return []byte("{}")
	}
	res, _ := json.Marshal(e)
	return res
}

// StatusCode return status code
func (e *StockResponseWrapper) StatusCode() int {
	if e == nil || e.Code == 0 {
		return http.StatusOK
	}
	return e.Code
}

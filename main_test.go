package main

import (
	"testing"
)

func Test_stockFind(t *testing.T) {
	// 	r := mux.NewRouter()
	// 	r.Handle("/stock/{id}", controller.GetStock())
	// 	ts := httptest.NewServer(r)
	// 	defer ts.Close()
	// 	t.Run("not found", func(t *testing.T) {
	// 		res, err := http.Get(ts.URL + "/v1/bookmark/1")
	// 		if err != nil {
	// 			t.Errorf("Expected nil, received %s", err.Error())
	// 		}
	// 		if res.StatusCode != http.StatusNotFound {
	// 			t.Errorf("Expected %d, received %d", http.StatusNotFound, res.StatusCode)
	// 		}
	// 	})
	// 	t.Run("found", func(t *testing.T) {
	// 		res, err := http.Get(ts.URL + "/v1/bookmark/2")
	// 		if err != nil {
	// 			t.Errorf("Expected nil, received %s", err.Error())
	// 		}
	// 		if res.StatusCode != http.StatusOK {
	// 			t.Errorf("Expected %d, received %d", http.StatusOK, res.StatusCode)
	// 		}
	// 	})
}

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"go-inventory/errors"
	"go-inventory/handlers"
	"go-inventory/objects"
	"go-inventory/store"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	router    *mux.Router
	flushAll  func(t *testing.T)
	createOne func(t *testing.T, name string) *objects.Stock
	getOne    func(t *testing.T, id string, wantErr bool) *objects.Stock
)

func TestMain(t *testing.M) {
	log.Println("Registering")

	conn := "postgres://postgres:root@localhost/test?sslmode=disable"
	if c := os.Getenv("DB_CONN"); c != "" {
		conn = c
	}

	router = mux.NewRouter().PathPrefix("/api/v1/").Subrouter()
	st := store.NewPostgresStockStore(conn)
	hnd := handlers.NewEventHandler(st)
	RegisterAllRoutes(router, hnd)

	flushAll = func(t *testing.T) {
		db, err := gorm.Open(postgres.Open(conn), nil)
		if err != nil {
			t.Fatal(err)
		}
		db.Delete(&objects.Stock{}, "1=1")
	}

	createOne = func(t *testing.T, name string) *objects.Stock {
		evt := &objects.Stock{
			Name:  name,
			Price: 0,
		}
		err := st.Create(context.TODO(), &objects.CreateRequest{Stock: evt})
		if err != nil {
			t.Fatal(err)
		}
		return evt
	}
	getOne = func(t *testing.T, id string, wantErr bool) *objects.Stock {
		evt, err := st.Get(context.TODO(), &objects.GetRequest{ID: id})
		if err != nil && wantErr {
			t.Fatal(err)
		}
		return evt
	}

	log.Println("Starting")
	os.Exit(t.Run())
}

func Do(req *http.Request) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func TestUnknownEndpoints(t *testing.T) {
	tests := []struct {
		name  string
		setup func(t *testing.T) *http.Request
	}{
		{
			name: "root",
			setup: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodGet, "/", nil)
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
		},
		{
			name: "api-root",
			setup: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodGet, "/api/v1", nil)
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
		},
		{
			name: "random",
			setup: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodGet, "/random", nil)
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := Do(tt.setup(t))
			_ = assert.Equal(t, http.StatusNotFound, w.Code) &&
				assert.Equal(t, "404 page not found\n", string(w.Body.Bytes()))
		})
	}
}

func TestGetEndpoint(t *testing.T) {
	flushAll(t)
	tests := []struct {
		name  string
		code  int
		setup func(t *testing.T) *http.Request
	}{
		{
			name: "OK",
			setup: func(t *testing.T) *http.Request {
				evt := createOne(t, "Ok")
				req, err := http.NewRequest(http.MethodGet, "/api/v1/stock?id="+evt.ID, nil)
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			code: http.StatusOK,
		},
		{
			name: "NotFound",
			setup: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodGet, "/api/v1/stock?id=32", nil)
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			code: http.StatusNotFound,
		},
		{
			name: "WithoutParameter",
			setup: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodGet, "/api/v1/stock", nil)
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			code: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := Do(tt.setup(t))
			assert.Equal(t, tt.code, w.Code)
			got := &objects.StockResponseWrapper{}
			assert.Nil(t, json.Unmarshal(w.Body.Bytes(), got))
		})
	}
}

func TestListEndpoint(t *testing.T) {
	flushAll(t)
	tests := []struct {
		name    string
		code    int
		setup   func(t *testing.T) *http.Request
		listLen int
	}{
		{
			name: "Zero",
			setup: func(t *testing.T) *http.Request {
				flushAll(t)
				req, err := http.NewRequest(http.MethodGet, "/api/v1/stocks", nil)
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			code:    http.StatusOK,
			listLen: 0,
		},
		{
			name: "All",
			setup: func(t *testing.T) *http.Request {
				_ = createOne(t, "One")
				_ = createOne(t, "Two")
				req, err := http.NewRequest(http.MethodGet, "/api/v1/stocks", nil)
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			code:    http.StatusOK,
			listLen: 2,
		},
		{
			name: "Limited",
			setup: func(t *testing.T) *http.Request {
				_ = createOne(t, "Three")
				req, err := http.NewRequest(http.MethodGet, "/api/v1/stock?limit=2", nil)
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			code:    http.StatusOK,
			listLen: 2,
		},
		{
			name: "After",
			setup: func(t *testing.T) *http.Request {
				evt := createOne(t, "Four")
				_ = createOne(t, "Five")
				req, err := http.NewRequest(http.MethodGet, "/api/v1/stock?after="+evt.ID, nil)
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			code:    http.StatusOK,
			listLen: 1,
		},
		{
			name: "Name",
			setup: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodGet, "/api/v1/stock?name=e", nil)
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			code:    http.StatusOK,
			listLen: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := Do(tt.setup(t))
			got := &objects.StockResponseWrapper{}
			assert.Equal(t, tt.code, w.Code)
			assert.Nil(t, json.Unmarshal(w.Body.Bytes(), got))
			assert.Equal(t, len(got.Stocks), tt.listLen)
		})
	}
}

func TestCreateEndpoint(t *testing.T) {
	flushAll(t)
	tests := []struct {
		name    string
		message string
		code    int
		evt     *objects.Stock
	}{
		{
			name:    "Ok",
			message: "",
			code:    http.StatusOK,
			evt: &objects.Stock{
				Name: "Help Ok",
			},
		},

		{
			name:    "No input",
			message: errors.ErrObjectIsRequired.Message,
			code:    errors.ErrObjectIsRequired.Code,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := json.Marshal(tt.evt)
			if err != nil {
				t.Fatal(err)
			}
			req, err := http.NewRequest(http.MethodPost, "/api/v1/stock", bytes.NewReader(b))
			if err != nil {
				t.Fatal(err)
			}
			w := Do(req)
			got, gotErr := &objects.StockResponseWrapper{}, &errors.Error{}
			assert.Equal(t, tt.code, w.Code)
			assert.Nil(t, json.Unmarshal(w.Body.Bytes(), got))
			assert.Nil(t, json.Unmarshal(w.Body.Bytes(), gotErr))
			assert.Equal(t, tt.message, gotErr.Message)
			if tt.code == http.StatusOK {
				ok := assert.NotNil(t, got.Stock) &&
					assert.NotEmpty(t, got.Stock.ID) &&
					assert.NotEmpty(t, got.Stock.CreatedOn)
				if ok {
					tt.evt.ID = got.Stock.ID
					tt.evt.CreatedOn = got.Stock.CreatedOn
					assert.Equal(t, tt.evt, got.Stock)
				}
			}
		})
	}
}

func TestUpdateDetailsEndpoint(t *testing.T) {
	flushAll(t)
	reqFn := func(t *testing.T, evt *objects.Stock) (*http.Request, *objects.Stock) {
		var (
			b   []byte
			err error
		)
		if evt != nil {
			b, err = json.Marshal(&objects.UpdateDetailsRequest{
				ID:   evt.ID,
				Name: evt.Name,
			})
			if err != nil {
				t.Fatal(err)
			}
		}
		req, err := http.NewRequest(http.MethodPut, "/api/v1/stock/details", bytes.NewReader(b))
		if err != nil {
			t.Fatal(err)
		}
		return req, evt
	}
	tests := []struct {
		name    string
		code    int
		setup   func(t *testing.T) (*http.Request, *objects.Stock)
		message string
	}{
		{
			name: "OK",
			setup: func(t *testing.T) (*http.Request, *objects.Stock) {
				evt := createOne(t, "Ok")
				// evt.Address = "Fake"
				// evt.PhoneNumber = "+23765"
				return reqFn(t, evt)
			},
			code: http.StatusOK,
		},
		{
			name: "NotFound",
			setup: func(t *testing.T) (*http.Request, *objects.Stock) {
				evt := createOne(t, "Ok")
				evt.ID = ""
				return reqFn(t, evt)
			},
			message: errors.ErrStockNotFound.Message,
			code:    http.StatusNotFound,
		},
		{
			name: "No input",
			setup: func(t *testing.T) (*http.Request, *objects.Stock) {
				return reqFn(t, nil)
			},
			message: errors.ErrObjectIsRequired.Message,
			code:    http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, exp := tt.setup(t)
			w := Do(req)
			assert.Equal(t, tt.code, w.Code)
			if tt.message != "" {
				got := &errors.Error{}
				assert.Nil(t, json.Unmarshal(w.Body.Bytes(), got))
				assert.Equal(t, tt.message, got.Message)
			} else if exp != nil {
				evt := getOne(t, exp.ID, true)
				assert.Equal(t, exp.Name, evt.Name)
				// assert.Equal(t, exp.Description, evt.Description)
				// assert.Equal(t, exp.Website, evt.Website)
				// assert.Equal(t, exp.Address, evt.Address)
				// assert.Equal(t, exp.PhoneNumber, evt.PhoneNumber)
			}
		})
	}
}

func TestDeleteEndpoint(t *testing.T) {
	flushAll(t)
	reqFn := func(t *testing.T, in *objects.DeleteRequest) (*http.Request, string) {
		id := ""
		if in != nil {
			id = in.ID
		}
		req, err := http.NewRequest(http.MethodDelete, "/api/v1/stock?id="+id, nil)
		if err != nil {
			t.Fatal(err)
		}
		return req, id
	}
	tests := []struct {
		name    string
		code    int
		setup   func(t *testing.T) (*http.Request, string)
		message string
	}{
		{
			name: "OK",
			setup: func(t *testing.T) (*http.Request, string) {
				evt := createOne(t, "Ok")
				return reqFn(t, &objects.DeleteRequest{ID: evt.ID})
			},
			code: http.StatusOK,
		},
		{
			name: "No input",
			setup: func(t *testing.T) (*http.Request, string) {
				return reqFn(t, nil)
			},
			message: errors.ErrValidStockIDIsRequired.Message,
			code:    http.StatusBadRequest,
		},
		{
			name: "NotFound",
			setup: func(t *testing.T) (*http.Request, string) {
				return reqFn(t, &objects.DeleteRequest{ID: "fake"})
			},
			message: errors.ErrStockNotFound.Message,
			code:    http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, id := tt.setup(t)
			w := Do(req)
			assert.Equal(t, tt.code, w.Code)
			if tt.message != "" {
				got := &errors.Error{}
				assert.Nil(t, json.Unmarshal(w.Body.Bytes(), got))
				assert.Equal(t, tt.message, got.Message)
			} else if id != "" {
				assert.Nil(t, getOne(t, id, false))
			}
		})
	}
}

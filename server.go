package main

import (
	"log"
	"net/http"

	"go-inventory/handlers"
	"go-inventory/store"

	"github.com/gorilla/mux"
)

// Args args used to run the server
type Args struct {
	// postgres connection string, of the form,
	// e.g "postgres://user:password@localhost:5432/database?sslmode=disable"
	conn string
	// port for the server of the form,
	// e.g ":8080"
	port string
}

// Run run the server based on given args
func Run(args Args) error {
	// router
	router := mux.NewRouter().
		PathPrefix("/api/v1/"). // add prefix for v1 api `/api/v1/`
		Subrouter()

	st := store.NewPostgresStockStore(args.conn)
	hnd := handlers.NewEventHandler(st)
	RegisterAllRoutes(router, hnd)

	// start server
	log.Println("Starting server at port: ", args.port)
	return http.ListenAndServe(args.port, router)
}

// RegisterAllRoutes registers all routes of the api
func RegisterAllRoutes(router *mux.Router, hnd handlers.IStockHandler) {

	// set content type
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			next.ServeHTTP(w, r)
		})
	})

	// get stock
	router.HandleFunc("/stock/{id}", hnd.Get).Methods(http.MethodGet)
	// create stock
	router.HandleFunc("/stock", hnd.Create).Methods(http.MethodPost)

	// update stock details
	router.HandleFunc("/stock/details", hnd.UpdateDetails).Methods(http.MethodPut)

	// list stock
	router.HandleFunc("/stocks", hnd.List).Methods(http.MethodGet)
}

package router

import (
	"go-inventory/controller"

	"github.com/go-chi/chi"
)

// func Router() *mux.Router {

// 	router := mux.NewRouter()
// 	// r.Method("GET", "/books", requestlog.NewHandler(srvBook.HandleListBooks, l))

// 	router.HandleFunc("/stock/{id}", requestlog.NewHandler(controller.GetStock).Methods("GET", "OPTIONS"))
// 	// router.HandleFunc("/stock/{id}", controller.GetStock).Methods("GET", "OPTIONS")
// 	// router.HandleFunc("/stock", controller.AddStock).Methods("POST", "OPTIONS")
// 	// router.HandleFunc("/stock/{id}", controller.UpdateStock).Methods("PUT", "OPTIONS")

// 	return router
// }
func New() *chi.Mux {

	r := chi.NewRouter()
	r.MethodFunc("GET", "/stock/{id}", controller.GetStock)
	r.MethodFunc("POST", "/stock", controller.AddStock)
	r.MethodFunc("PUT", "/stock/{id}", controller.UpdateStock)

	return r
}

package products

import (
	"my-pp/modules/users"
	"net/http"

	"github.com/gorilla/mux"
)

func ListRoutes(r *mux.Router) {
	r.Handle("/list", users.ValidateJWT(http.HandlerFunc(createItemHandler))).Methods("POST")
	r.Handle("/list/{id}", users.ValidateJWT(http.HandlerFunc(getItemByIdHandler))).Methods("GET")
	r.Handle("/list", users.ValidateJWT(http.HandlerFunc(getItemsHandler))).Methods("GET")
	r.Handle("/list/{id}", users.ValidateJWT(http.HandlerFunc(updateItemHandler))).Methods("PATCH")
	r.Handle("/list/{id}", users.ValidateJWT(http.HandlerFunc(deleteItemHandler))).Methods("DELETE")
}

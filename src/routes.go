package src

import (
	"github.com/gorilla/mux"
)

func SetupRoutes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/list", createItemHandler).Methods("POST")
	r.HandleFunc("/list/{id}", getItemByIdHandler).Methods("GET")
	r.HandleFunc("/list", getItemsHandler).Methods("GET")
	r.HandleFunc("/list/{id}", updateItemHandler).Methods("PATCH")
	r.HandleFunc("/list/{id}", deleteItemHandler).Methods("DELETE")
	r.HandleFunc("/user/login", getUserLoginHandler).Methods("GET")

	return r
}

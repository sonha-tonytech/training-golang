package main

import (
	"my-pp/modules/products"
	"my-pp/modules/users"

	"github.com/gorilla/mux"
)

func SetupRoutes() *mux.Router {
	r := mux.NewRouter()

	products.ListRoutes(r)
	users.UserRoutes(r)

	return r
}

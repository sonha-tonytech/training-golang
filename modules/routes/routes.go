package routes

import (
	"my-pp/modules/products"
	servermiddlewares "my-pp/modules/server-middlewares"
	"my-pp/modules/users"

	"github.com/gorilla/mux"
)

func SetupRoutes() *mux.Router {
	r := mux.NewRouter()

	products.ListRoutes(r)
	users.UserRoutes(r)

	r.Use(servermiddlewares.RecoverMiddleware)
	r.Use(servermiddlewares.LoggingMiddleware)

	return r
}

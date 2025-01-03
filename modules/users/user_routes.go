package users

import (
	"net/http"

	"github.com/gorilla/mux"
)

func UserRoutes(r *mux.Router) {
	r.Handle("/user/login", VadidateLoginUser(http.HandlerFunc(getUserLoginHandler))).Methods("GET")
	r.Handle("/user/register", VadidateRegisterUser(http.HandlerFunc(registerUserHandler))).Methods("POST")
	r.Handle("/user/{id}", ValidateJWT(http.HandlerFunc(updateUserHandler))).Methods("PATCH")
}

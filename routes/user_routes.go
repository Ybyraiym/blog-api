package routes

import (
	"blog-api/handlers"

	"github.com/gorilla/mux"
)

func RegisterUserRoutes(r *mux.Router) {
	r.HandleFunc("/register", handlers.RegisterUser).Methods("POST")
	r.HandleFunc("/authenticate", handlers.AuthenticateUser).Methods("POST")
}

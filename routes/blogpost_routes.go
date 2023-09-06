package routes

import (
	"blog-api/handlers"

	"github.com/gorilla/mux"
)

func RegisterBlogPostRoutes(r *mux.Router) {
	r.HandleFunc("/blogposts", handlers.CreateBlogPost).Methods("POST")
	r.HandleFunc("/blogposts", handlers.GetBlogPosts).Methods("GET")
	r.HandleFunc("/blogposts/{id}", handlers.GetBlogPost).Methods("GET")
	r.HandleFunc("/blogposts/{id}", handlers.UpdateBlogPost).Methods("PUT")
	r.HandleFunc("/blogposts/{id}", handlers.DeleteBlogPost).Methods("DELETE")
}

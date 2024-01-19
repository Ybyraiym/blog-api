package routes

import (
	"blog-api/handlers"
	"blog-api/middleware"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func SetupBlogPostRoutes(router *mux.Router) {
	router.Use(middleware.AuthMiddleware)

	router.HandleFunc("/api/blogposts", logRequest(handlers.GetBlogPosts)).Methods("GET")
	router.HandleFunc("/api/blogposts/{id}", logRequest(handlers.GetBlogPost)).Methods("GET")
	router.HandleFunc("/api/blogposts", logRequest(handlers.CreateBlogPost)).Methods("POST")
	router.HandleFunc("/api/blogposts/{id}", logRequest(handlers.UpdateBlogPost)).Methods("PUT")
	router.HandleFunc("/api/blogposts/{id}", logRequest(handlers.DeleteBlogPost)).Methods("DELETE")
}

func logRequest(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logrus.Info("Received request:", r.Method, r.URL)
		handler(w, r)
	}
}

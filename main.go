package main

import (
	"blog-api/db"
	blog "blog-api/pkg"
	"blog-api/routes"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// Подключение к базе данных
	db, err := db.InitDatabase()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Запуск сервера
	srv := new(blog.Server)
	if err := srv.Run("8080"); err != nil {
		log.Fatalf("error occured while http server: %s", err.Error())
	}

	// Регистрация маршрутов
	r := mux.NewRouter()
	routes.RegisterUserRoutes(r)
	routes.RegisterBlogPostRoutes(r)

	http.Handle("/", r)

	// Запуск веб-сервера
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

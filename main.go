package main

import (
	"blog-api/db"
	blog "blog-api/pkg"
	"blog-api/pkg/logging"
	"blog-api/routes"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// Инициализация логгера
	logger := logging.NewLogger().GetLogger()

	// Пример использования логгера
	logger.Info("Starting the application")

	// Подключение к базе данных
	db, err := db.InitDatabase()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Регистрация маршрутов
	r := mux.NewRouter()
	routes.RegisterUserRoutes(r)
	routes.SetupBlogPostRoutes(r)

	http.Handle("/", r)

	// Запуск сервера
	srv := new(blog.Server)
	if err := srv.Run("8080"); err != nil {
		log.Fatalf("error occured while http server: %s", err.Error())
	}

}

package handlers

import (
	"blog-api/models"
	"encoding/json"
	"net/http"

	"github.com/jinzhu/gorm"
)

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var newUser models.User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Устанавливаем соединение с базой данных
	db, err := connectToDatabase()
	if err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Проверка наличия обязательных полей
	if newUser.Login == "" || newUser.Username == "" || newUser.Password == "" {
		http.Error(w, "Login, Username and password are required", http.StatusBadRequest)
		return
	}

	// Сохранение нового пользователя в базе данных
	if err := db.Create(&newUser).Error; err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Вернуть успешный статус
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)
}

// Функция для подключения к базе данных
func connectToDatabase() (*gorm.DB, error) {
	// Здесь укажите параметры подключения к вашей PostgreSQL базе данных
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=blogdb sslmode=disable password=postgres")
	if err != nil {
		return nil, err
	}
	return db, nil
}

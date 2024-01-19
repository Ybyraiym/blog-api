package handlers

import (
	"blog-api/models"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"blog-api/authentication"

	"github.com/golang-jwt/jwt"
	"github.com/jinzhu/gorm"
)

// Вспомогательная функция для получения имени пользователя из токена
func GetUsernameFromToken(tokenString string) (string, error) {
	// Проверяем, что токен не пустой
	if tokenString == "" {
		return "", errors.New("токен отсутствует")
	}

	// Разбиваем токен на части
	tokenParts := strings.Split(tokenString, ".")
	if len(tokenParts) != 3 {
		return "", errors.New("неверный формат токена")
	}

	// Проверяем токен и извлекаем из него имя пользователя
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Проверка метода подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неподдерживаемый метод подписи: %v", token.Header["alg"])
		}

		// Возвращаем секретный ключ
		return authentication.TokenSecret, nil
	})
	if err != nil {
		return "", err
	}

	// Проверяем валидность токена
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if username, ok := claims["username"].(string); ok {
			return username, nil
		}
	}

	return "", errors.New("неверный токен")
}

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
	if newUser.Login == "" && newUser.Username == "" && newUser.Password == "" {
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

func AuthenticateUser(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&requestBody); err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	login := requestBody.Login
	password := requestBody.Password

	// Устанавливаем соединение с базой данных
	db, err := connectToDatabase()
	if err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Аутентификация
	authenticator := &authentication.PostgreSQLAuthenticator{DB: db.DB()}
	isAuthenticated, err := authenticator.Authenticate(login, password)

	if err != nil {
		http.Error(w, "Ошибка аутентификации", http.StatusInternalServerError)
		return
	}

	if isAuthenticated {
		// Аутентификация успешна. Генерация токена.
		tokenString, err := authentication.GenerateToken(login)
		if err != nil {
			http.Error(w, "Ошибка генерации токена", http.StatusInternalServerError)
			return
		}

		// Отправляем токен в ответе
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
	} else {
		// Неверный логин или пароль.
		http.Error(w, "Неверные учетные данные", http.StatusUnauthorized)
	}
}

// Функция для подключения к базе данных
func connectToDatabase() (*gorm.DB, error) {
	// Здесь параметры подключения к PostgreSQL базе данных
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=blogdb sslmode=disable password=postgres")
	if err != nil {
		return nil, err
	}
	return db, nil
}

package handlers

import (
	"blog-api/authentication"
	"blog-api/db"
	"blog-api/models"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

// Вспомогательная функция для извлечения токена из запроса
func getTokenFromRequest(r *http.Request) (string, error) {
	// Получаем токен из заголовка Authorization
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		return "", errors.New("токен отсутствует в запросе")
	}

	// Возвращаем токен без префикса "Bearer "
	return strings.TrimPrefix(tokenString, "Bearer "), nil
}

func CreateBlogPost(w http.ResponseWriter, r *http.Request) {
	// Извлекаем токен из запроса
	token, err := getTokenFromRequest(r)
	if err != nil {
		HandleError(w, "Ошибка получения токена", err, http.StatusUnauthorized)
		return
	}

	// Извлекаем информацию о пользователе из токена
	username, err := authentication.GetUsernameFromToken(token)
	if err != nil {
		HandleError(w, "Ошибка аутентификации", err, http.StatusUnauthorized)
		return
	}

	// Добавляем токен в заголовок запроса на создание блогпоста
	r.Header.Set("Authorization", "Bearer "+token)

	var newPost models.BlogPost
	err = json.NewDecoder(r.Body).Decode(&newPost)
	if err != nil {
		HandleError(w, "Invalid JSON", err, http.StatusBadRequest)
		return
	}

	newPost.PublishedAt = time.Now()
	newPost.Username = username

	// Сохраняем запись в базе данных
	if err := db.GlobalDB.Create(&newPost).Error; err != nil {
		HandleError(w, "Database error", err, http.StatusInternalServerError)
		return
	}

	Logger.Info("Блогпост успешно создан")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newPost)
}

// Вспомогательная функция для получения имени пользователя из токена
func getUsernameFromToken(r *http.Request) (string, error) {
	return "", nil
}

func GetBlogPosts(w http.ResponseWriter, r *http.Request) {
	var posts []models.BlogPost
	if err := db.GlobalDB.Find(&posts).Error; err != nil {
		HandleError(w, "Database error", err, http.StatusInternalServerError)
		return
	}

	Logger.Info("Список блогпостов успешно получен")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func GetBlogPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	postID := params["id"]

	var post models.BlogPost
	if err := db.GlobalDB.Where("id = ?", postID).First(&post).Error; err != nil {
		HandleError(w, "Блогпост не найден", err, http.StatusNotFound)
		return
	}

	Logger.Info("Блогпост успешно получен")
	json.NewEncoder(w).Encode(post)
}

func UpdateBlogPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	postID := params["id"]
	var updatedPost models.BlogPost
	err := json.NewDecoder(r.Body).Decode(&updatedPost)
	if err != nil {
		HandleError(w, "Invalid JSON", err, http.StatusBadRequest)
		return
	}
	var existingPost models.BlogPost
	if err := db.GlobalDB.Where("id = ?", postID).First(&existingPost).Error; err != nil {
		HandleError(w, "Блогпост не найден", err, http.StatusNotFound)
		return
	}
	updatedPost.ID = existingPost.ID
	updatedPost.PublishedAt = existingPost.PublishedAt

	// Обновление записи в базе данных
	if err := db.GlobalDB.Save(&updatedPost).Error; err != nil {
		HandleError(w, "Database error", err, http.StatusInternalServerError)
		return
	}

	Logger.Info("Блогпост успешно обновлен")
	json.NewEncoder(w).Encode(updatedPost)
}

func DeleteBlogPost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	postID := params["id"]

	// Удаление записи из базы данных
	if err := db.GlobalDB.Where("id = ?", postID).Delete(&models.BlogPost{}).Error; err != nil {
		HandleError(w, "Database error", err, http.StatusInternalServerError)
		return
	}

	Logger.Info("Блогпост успешно удален")
	w.WriteHeader(http.StatusNoContent)
}

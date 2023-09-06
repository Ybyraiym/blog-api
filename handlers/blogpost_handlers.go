package handlers

import (
	"blog-api/db"
	"blog-api/models"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func CreateBlogPost(w http.ResponseWriter, r *http.Request) {
	var newPost models.BlogPost
	err := json.NewDecoder(r.Body).Decode(&newPost)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	newPost.PublishedAt = time.Now()

	// Сохранение записи в базе данных
	if err := db.GlobalDB.Create(&newPost).Error; err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newPost)
}

func GetBlogPosts(w http.ResponseWriter, r *http.Request) {
	var posts []models.BlogPost
	if err := db.GlobalDB.Find(&posts).Error; err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func GetBlogPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	postID := params["id"]

	var post models.BlogPost
	if err := db.GlobalDB.Where("id = ?", postID).First(&post).Error; err != nil {
		http.NotFound(w, r)
		return
	}

	json.NewEncoder(w).Encode(post)
}

func UpdateBlogPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	postID := params["id"]
	var updatedPost models.BlogPost
	err := json.NewDecoder(r.Body).Decode(&updatedPost)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	var existingPost models.BlogPost
	if err := db.GlobalDB.Where("id = ?", postID).First(&existingPost).Error; err != nil {
		http.NotFound(w, r)
		return
	}
	updatedPost.ID = existingPost.ID
	updatedPost.PublishedAt = existingPost.PublishedAt

	// Обновление записи в базе данных
	if err := db.GlobalDB.Save(&updatedPost).Error; err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(updatedPost)
}

func DeleteBlogPost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	postID := params["id"]

	// Удаление записи из базы данных
	if err := db.GlobalDB.Where("id = ?", postID).Delete(&models.BlogPost{}).Error; err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

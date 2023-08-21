package main

import (
	"encoding/json"
	// "fmt"
	"net/http"
	"time"

	"log"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// Структура данных
type BlogPost struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	PublishedAt time.Time `json:"publishedAt"`
}

// Глобальная переменная для хранения подключения к базе данных
var db *gorm.DB

func main() {
	// Подключение к базе данных
	var err error
	db, err = gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=blogdb sslmode=disable password=postgres")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Автомиграция для создания таблицы, если её нет
	db.AutoMigrate(&BlogPost{})

	r := mux.NewRouter()
	r.HandleFunc("/blogposts", CreateBlogPost).Methods("POST")
	r.HandleFunc("/blogposts", GetBlogPosts).Methods("GET")
	r.HandleFunc("/blogposts/{id}", GetBlogPost).Methods("GET")
	r.HandleFunc("/blogposts/{id}", UpdateBlogPost).Methods("PUT")
	r.HandleFunc("/blogposts/{id}", DeleteBlogPost).Methods("DELETE")
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}

// Метод для создания новой записи блога
func CreateBlogPost(w http.ResponseWriter, r *http.Request) {
	var newPost BlogPost
	err := json.NewDecoder(r.Body).Decode(&newPost)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	newPost.PublishedAt = time.Now()

	// Сохранение записи в базе данных
	if err := db.Create(&newPost).Error; err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newPost)
}

// Метод для получения списка всех записей блога из базы данных
func GetBlogPosts(w http.ResponseWriter, r *http.Request) {
	var posts []BlogPost
	if err := db.Find(&posts).Error; err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

// Метод для получения конкретной записи блога по ID из базы данных
func GetBlogPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	postID := params["id"]

	var post BlogPost
	if err := db.Where("id = ?", postID).First(&post).Error; err != nil {
		http.NotFound(w, r)
		return
	}

	json.NewEncoder(w).Encode(post)
}

// Метод для обновления записи блога по ID
func UpdateBlogPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	postID := params["id"]
	var updatedPost BlogPost
	err := json.NewDecoder(r.Body).Decode(&updatedPost)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	var existingPost BlogPost
	if err := db.Where("id = ?", postID).First(&existingPost).Error; err != nil {
		http.NotFound(w, r)
		return
	}
	updatedPost.ID = existingPost.ID
	updatedPost.PublishedAt = existingPost.PublishedAt

	// Обновление записи в базе данных
	if err := db.Save(&updatedPost).Error; err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(updatedPost)
}

// Метод для удаления записи блога по ID
func DeleteBlogPost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	postID := params["id"]

	// Удаление записи из базы данных
	if err := db.Where("id = ?", postID).Delete(&BlogPost{}).Error; err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

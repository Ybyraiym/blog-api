package main

import (
	"encoding/json"
	"fmt"
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

// Объявление переменных "Пост" "ID счетчик"
var blogPosts []BlogPost
var idCounter uint = 1

// Метод для создания новой записи блога
func CreateBlogPost(w http.ResponseWriter, r *http.Request) {
	var newPost BlogPost
	err := json.NewDecoder(r.Body).Decode(&newPost)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	newPost.ID = idCounter
	idCounter++
	newPost.PublishedAt = time.Now()
	blogPosts = append(blogPosts, newPost)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newPost)
}

// Метод для получения списка всех записей блога
func GetBlogPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(blogPosts)
}

// Метод для получения конкретной записи блога по ID
func GetBlogPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	postID := params["id"]
	for _, post := range blogPosts {
		if fmt.Sprint(post.ID) == postID {
			json.NewEncoder(w).Encode(post)
			return
		}
	}
	http.NotFound(w, r)
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
	for index, post := range blogPosts {
		if fmt.Sprint(post.ID) == postID {
			updatedPost.ID = post.ID
			updatedPost.PublishedAt = post.PublishedAt
			blogPosts[index] = updatedPost
			json.NewEncoder(w).Encode(updatedPost)
			return
		}
	}
	http.NotFound(w, r)
}

// Метод для удаления записи блога по ID
func DeleteBlogPost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	postID := params["id"]
	for index, post := range blogPosts {
		if fmt.Sprint(post.ID) == postID {
			blogPosts = append(blogPosts[:index], blogPosts[index+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	http.NotFound(w, r)
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

	// Ваш существующий код
	r := mux.NewRouter()
	r.HandleFunc("/blogposts", CreateBlogPost).Methods("POST")
	r.HandleFunc("/blogposts", GetBlogPosts).Methods("GET")
	r.HandleFunc("/blogposts/{id}", GetBlogPost).Methods("GET")
	r.HandleFunc("/blogposts/{id}", UpdateBlogPost).Methods("PUT")
	r.HandleFunc("/blogposts/{id}", DeleteBlogPost).Methods("DELETE")
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}

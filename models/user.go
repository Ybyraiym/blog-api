package models

type User struct {
	ID       uint   `json:"id"`
	Login    string `json:"login"`
	Username string `json:"username"`
	Password string `json:"password"`
}

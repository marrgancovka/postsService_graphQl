package models

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type UserData struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

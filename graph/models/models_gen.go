// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package models

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type NewPost struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

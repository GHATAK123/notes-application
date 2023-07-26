package main

type User struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Note struct {
	ID   uint32 `json:"id"`
	Note string `json:"note"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type AuthResponse struct {
	SID string `json:"sid"`
}

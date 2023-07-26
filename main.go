package main

import (
	"fmt"
	"net/http"
)

func main() {
	router := http.NewServeMux()

	router.HandleFunc("/signup", handleSignup)
	router.HandleFunc("/login", handleLogin)
	router.HandleFunc("/notes", handleGetNotes)
	router.HandleFunc("/notes/create", handleCreateNote)
	router.HandleFunc("/notes/delete", handleDeleteNote)

	fmt.Println("Starting server on port 8080...")
	http.ListenAndServe(":8080", router)
}

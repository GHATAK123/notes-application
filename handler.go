package main

import (
	"encoding/json"
	"net/http"
)

var users = make(map[string]User)
var notes = make(map[string][]Note)

func handleSignup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var newUser User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil && (newUser.Name == "" || newUser.Email == "" || newUser.Password == "") {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	// Assuming email is unique, check if the user already exists
	if _, ok := users[newUser.Email]; ok {
		http.Error(w, "User already exists", http.StatusBadRequest)
		return
	}

	users[newUser.Email] = newUser
	w.WriteHeader(http.StatusOK)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var credentials User
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	user, ok := users[credentials.Email]
	if !ok || user.Password != credentials.Password {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// session_id := uuid.New().String()
	authResponse := AuthResponse{SID: credentials.Email} // In this example, we're using the email as the session ID
	json.NewEncoder(w).Encode(authResponse)
}

func handleGetNotes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var sid AuthResponse
	err := json.NewDecoder(r.Body).Decode(&sid)
	if err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	_, ok := users[sid.SID]
	if !ok {
		http.Error(w, "Invalid session ID", http.StatusUnauthorized)
		return
	}

	userNotes, ok := notes[sid.SID]
	if !ok {
		userNotes = []Note{}
	}

	response := map[string][]Note{"notes": userNotes}
	json.NewEncoder(w).Encode(response)
}

func handleCreateNote(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var noteData struct {
		SID  string `json:"sid"`
		Note string `json:"note"`
	}

	err := json.NewDecoder(r.Body).Decode(&noteData)
	if err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	_, ok := users[noteData.SID]
	if !ok {
		http.Error(w, "Invalid session ID", http.StatusUnauthorized)
		return
	}

	noteID := uint32(len(notes[noteData.SID]) + 1)
	newNote := Note{ID: noteID, Note: noteData.Note}
	notes[noteData.SID] = append(notes[noteData.SID], newNote)

	response := map[string]uint32{"id": noteID}
	json.NewEncoder(w).Encode(response)
}

func handleDeleteNote(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var noteData struct {
		SID string `json:"sid"`
		ID  uint32 `json:"id"`
	}

	err := json.NewDecoder(r.Body).Decode(&noteData)
	if err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	_, ok := users[noteData.SID]
	if !ok {
		http.Error(w, "Invalid session ID", http.StatusUnauthorized)
		return
	}

	userNotes, ok := notes[noteData.SID]
	if !ok {
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	}

	var updatedNotes []Note
	for _, note := range userNotes {
		if note.ID != noteData.ID {
			updatedNotes = append(updatedNotes, note)
		}
	}

	notes[noteData.SID] = updatedNotes
	w.WriteHeader(http.StatusOK)
}

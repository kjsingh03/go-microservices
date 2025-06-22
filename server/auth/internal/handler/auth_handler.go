package handler

import (
	"authentication/internal/model"
	"authentication/utils"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var signUpPayload struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		FirstName string `json:"first_name,omitempty"`
		LastName  string `json:"last_name,omitempty"`
	}

	err := json.NewDecoder(r.Body).Decode(&signUpPayload)
	if err != nil {
		utils.ErrorJSON(w, errors.New("invalid request payload"), http.StatusBadRequest)
		return
	}

	newUser := model.User{
		Email:     signUpPayload.Email,
		Password:  signUpPayload.Password,
		FirstName: signUpPayload.FirstName,
		LastName:  signUpPayload.LastName,
		Active:    true,
	}

	id, err := newUser.Insert(newUser)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]any{
		"message": "user registered successfully",
		"user_id": id,
	})

	log.Printf("User registered: %s", signUpPayload.Email)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.ErrorJSON(w, errors.New("invalid request"), http.StatusBadRequest)
		return
	}

	// Create a User instance to access methods
	var u model.User
	
	// Get user by email
	user, err := u.GetByEmail(payload.Email)

	log.Println(user)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusUnauthorized)
		return
	}

	// Check password
	match, err := user.PasswordMatches(payload.Password)
	if err != nil || !match {
		utils.ErrorJSON(w, err, http.StatusUnauthorized)
		return
	}

	// TODO: Generate JWT token here

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"message": "login successful",
		"user_id": user.ID,
	})
}
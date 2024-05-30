package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/fmdunlap/unhash/internal/uerr"
)

func (app *application) createUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}

	err := app.readJSON(r, &input)
	if err != nil {
		http.Error(w, fmt.Sprintf("error reading JSON: %v", err), http.StatusBadRequest)
		return
	}

	existingUser, err := app.userService.GetUserByEmail(input.Email)
	if err != nil {
		if !errors.Is(err, &uerr.ErrorNotFound{}) {
			http.Error(w, fmt.Sprintf("error getting user when checking for existing user: %v", err), http.StatusInternalServerError)
			return
		}
	}

	if existingUser != nil {
		http.Error(w, fmt.Sprintf("user with email `%s` already exists", input.Email), http.StatusConflict)
		return
	}

	userId := app.generateID()
	err = app.userService.CreateUser(userId, input.Username, input.Email)
	if err != nil {
		http.Error(w, fmt.Sprintf("error inserting user: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%s", userId)
}

func (app *application) getUserQueryHandler(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "missing email query parameter", http.StatusBadRequest)
		return
	}

	log.Println("email:", email)

	user, err := app.userService.GetUserByEmail(email)
	if err != nil {
		http.Error(w, fmt.Sprintf("error getting user: %v", err), http.StatusInternalServerError)
		return
	}

	err = app.writeJSON(w, http.StatusOK, user, nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("error writing JSON: %v", err), http.StatusInternalServerError)
		return
	}
}

func (app *application) getUserByIdHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := app.userService.GetUser(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("error getting user: %v", err), http.StatusInternalServerError)
		return
	}

	err = app.writeJSON(w, http.StatusOK, user, nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("error writing JSON: %v", err), http.StatusInternalServerError)
		return
	}
}

func (app *application) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = app.userService.DeleteUser(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("error deleting user: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "User with id `%s` deleted", id)
}

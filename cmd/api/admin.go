package main

import (
	"net/http"
)

func (app *application) adminHandler(w http.ResponseWriter, r *http.Request) {
	users, err := app.userService.ListAllUsers()
	if err != nil {
		http.Error(w, "error listing users", http.StatusInternalServerError)
		return
	}

	err = app.writeJSON(w, http.StatusOK, users, nil)
	if err != nil {
		http.Error(w, "error writing JSON", http.StatusInternalServerError)
		return
	}
}

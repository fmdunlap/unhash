package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/fmdunlap/unhash/internal/uerr"
)

func (app *application) createHashJobHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		OwnerId string   `json:"ownerId"`
		Hashes  []string `json:"hashes"`
	}

	err := app.readJSON(r, &input)
	if err != nil {
		http.Error(w, fmt.Sprintf("error reading JSON: %v", err), http.StatusBadRequest)
		return
	}

	owner, err := app.userService.GetUser(input.OwnerId)
	if err != nil {
		if errors.Is(err, &uerr.ErrorNotFound{}) {
			http.Error(w, fmt.Sprintf("user with id `%v` not found", input.OwnerId), http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("error getting user: %v", err), http.StatusInternalServerError)
		return
	}

	hashjobId, err := app.hashJobService.CreateHashJob(input.Hashes, owner)
	if err != nil {
		http.Error(w, fmt.Sprintf("error creating hash job: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%s", hashjobId)
}

func (app *application) getHashJobHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hj, err := app.hashJobService.GetHashJob(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = app.writeJSON(w, http.StatusOK, hj, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (app *application) deleteHashJobHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = app.hashJobService.DeleteHashJob(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "hash job deleted")
}

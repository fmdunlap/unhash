package main

import (
	"fmt"
	"net/http"
)

func (app *application) createHashJobHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "create a new hash job")
}

func (app *application) getHashJobHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "get hash job with id %d\n", id)
}

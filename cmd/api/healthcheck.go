package main

import (
	"fmt"
	"net/http"
	"time"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "status: available")
	fmt.Fprintf(w, "environment: %s\n", app.config.env)
	fmt.Fprintf(w, "version: %s\n", version)
	fmt.Fprintf(w, "time: %v", time.Now().UTC().Format(time.RFC3339))
}

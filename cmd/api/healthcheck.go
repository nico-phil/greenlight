package main

import (
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {

	data := map[string]string{
		"status":      "available/n",
		"environment": app.config.env,
		"version":     version,
	}

	err := app.writeJson(w, http.StatusOK, data, nil)
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "The server encoutered a problem and could not process your request", http.StatusInternalServerError)
	}

}

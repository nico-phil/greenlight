package main

import (
	"net/http"
)

func (app *application) createUserHandler(w http.ResponseWriter, r *http.Request) {

	err := app.writeJson(w, http.StatusOK, envelope{"hello": "world"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

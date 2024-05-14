package main

import (
	"errors"
	"net/http"

	"github.com/Nico2220/greenlight/internal/data"
	"github.com/Nico2220/greenlight/internal/validator"
)

func (app *application) createUserHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &data.User{
		Name:      input.Name,
		Email:     input.Email,
		Activated: false,
	}

	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	v := validator.New()

	if data.ValidateUser(v, user); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Users.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	//lauch a goroutine to send the email concurrently.
	// if you do not use a goroutine, the request will take around 2 secs, witch is a lot for a http request.
	//on the other hand, if you use a go routine, the request will take around 0.268... second, witch great

	app.background(func() {
		err = app.mailer.Send(user.Email, "user_welcome.tmpl", user) // this does not work because them smtp provider return a timeout error. later you have to change the smtp provider
		if err != nil {
			app.logger.Error(err.Error())
			return
		}
	})

	err = app.writeJson(w, http.StatusCreated, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

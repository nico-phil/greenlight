package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)


func TestListMoviesHandler(t *testing.T){
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	app := application{}
	app.listMoviesHandler(w, req)

	res := w.Result()
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("there is an error in the request"))
	}

	fmt.Println(string(body))


}
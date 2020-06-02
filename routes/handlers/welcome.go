package handlers

import (
	"net/http"

	"adiachenko/go-scaffold/routes/responses"

	"github.com/go-chi/render"
)

func Welcome(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, responses.NewWelcomeResponse("Welcome"))
}

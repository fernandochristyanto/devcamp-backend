package handler

import (
	"database/sql"
	"github.com/fernandochristyanto/devcamp-backend/internal"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// Handler object used to handle the HTTP API
type Handler struct {

	// DB object that'll be used
	DB *sql.DB
}

// Index is the home page handler
func (h *Handler) Index(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	internal.RenderJSON(w, []byte(`
		{
			"module": "search",
			"version": "1.0.0", 
			"tagline": "You know, for search"
		}
	`), http.StatusOK)
}

// ServeHTTP is used for 404 page
func (h *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	internal.RenderJSON(w, []byte(`
		{
			"message": "There's nothing here"
		}
	`), http.StatusNotFound)
}

package http

import (
	"musiclib/internal/song"

	"github.com/gorilla/mux"
)

// Map songs routes
func MapSongRoutes(newsGroup *mux.Router, h song.Handlers) {
	newsGroup.HandleFunc("/list", h.GetList).Methods("GET")
	newsGroup.HandleFunc("/text", h.GetText).Methods("GET")
	newsGroup.HandleFunc("/", h.Delete).Methods("DELETE")
	newsGroup.HandleFunc("/", h.Update).Methods("PUT")
	newsGroup.HandleFunc("/", h.Add).Methods("POST")
}

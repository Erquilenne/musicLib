package server

import (
	"github.com/gorilla/mux"

	songHttp "musiclib/internal/song/delivery/http"
)

// MapHandlers Map Server Handlers
func (s *Server) MapHandlers(router *mux.Router) error {
	// Init handlers
	songHandlers := songHttp.NewSongHandlers(s.cfg, s.logger)

	// API endpoints
	apiRouter := router.PathPrefix("/api/v1").Subrouter()

	// Songs routes
	songsGroup := apiRouter.PathPrefix("/songs").Subrouter()
	songHttp.MapSongRoutes(songsGroup, songHandlers)

	return nil
}

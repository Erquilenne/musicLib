package server

import (
	"github.com/gorilla/mux"

	songHttp "musiclib/internal/song/delivery/http"
	"musiclib/internal/song/repository"
)

// MapHandlers Map Server Handlers
func (s *Server) MapHandlers(router *mux.Router) error {
	// Init repository
	songRepo := repository.NewSongRepository(s.db)

	// Init handlers
	songHandlers := songHttp.NewSongHandlers(s.cfg, s.logger, songRepo)

	// API endpoints
	apiRouter := router.PathPrefix("/api/v1").Subrouter()

	// Songs routes
	songsGroup := apiRouter.PathPrefix("/songs").Subrouter()
	songHttp.MapSongRoutes(songsGroup, songHandlers)

	return nil
}

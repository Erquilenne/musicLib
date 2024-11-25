package server

import (
	"github.com/gorilla/mux"

	songHttp "musiclib/internal/song/delivery/http"
	"musiclib/internal/song/repository"
)

// MapHandlers Map Server Handlers
func (s *Server) MapHandlers(router *mux.Router) error {
	songRepo := repository.NewSongRepository(s.db, s.logger)

	songHandlers := songHttp.NewSongHandlers(s.cfg, s.logger, songRepo)

	apiRouter := router.PathPrefix("/api/v1").Subrouter()

	songsGroup := apiRouter.PathPrefix("/songs").Subrouter()
	songHttp.MapSongRoutes(songsGroup, songHandlers)

	return nil
}

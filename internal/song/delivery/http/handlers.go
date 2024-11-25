package http

import (
	"encoding/json"
	"fmt"
	"musiclib/config"
	"musiclib/internal/models"
	"musiclib/internal/song"
	"musiclib/pkg/logger"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const defaultLimit = "10"
const defaultOffset = "0"
const defaultSortBy = "id"
const defaultSortOrder = "asc"

// Song handlers
type songHandlers struct {
	cfg      *config.Config
	songRepo song.Repository
	logger   logger.Logger
}

// NewSongHandlers Song handlers constructor
func NewSongHandlers(cfg *config.Config, logger logger.Logger, repo song.Repository) *songHandlers {
	return &songHandlers{cfg: cfg, logger: logger, songRepo: repo}
}

// @Summary     List songs
// @Description Get paginated and sorted list of songs
// @Tags        songs
// @Accept      json
// @Produce     json
// @Param       sort_by query string false "Field to sort by (default: id)"
// @Param       sort_order query string false "Sort order: asc or desc (default: asc)"
// @Param       limit query int false "Number of items to return (default: 10)"
// @Param       offset query int false "Number of items to skip (default: 0)"
// @Success     200 {array} models.Song
// @Failure     400 {object} models.ErrorResponse
// @Failure     500 {object} models.ErrorResponse
// @Router      /songs/list [get]
func (h *songHandlers) GetList(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("Starting GetList handler")
	// Get query parameters for sorting and pagination
	sortBy := r.URL.Query().Get("sort_by")
	sortOrder := r.URL.Query().Get("sort_order")
	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")

	h.logger.Debug("Raw query parameters",
		"sortBy", sortBy,
		"sortOrder", sortOrder,
		"limit", limit,
		"offset", offset,
	)

	// Validate and set default values for query parameters
	if sortBy == "" {
		sortBy = defaultSortBy
	}
	if sortOrder == "" {
		sortOrder = defaultSortOrder
	}
	if limit == "" {
		limit = defaultLimit
	}
	if offset == "" {
		offset = defaultOffset
	}

	h.logger.Debug("Normalized query parameters",
		"sortBy", sortBy,
		"sortOrder", sortOrder,
		"limit", limit,
		"offset", offset,
	)

	// Convert query parameters to integers
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		h.logger.Error("Invalid limit value", err)
		http.Error(w, "Invalid limit value", http.StatusBadRequest)
		return
	}
	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		h.logger.Error("Invalid offset value", err)
		http.Error(w, "Invalid offset value", http.StatusBadRequest)
		return
	}

	h.logger.Debug("Converted parameters",
		"limitInt", limitInt,
		"offsetInt", offsetInt,
	)

	// Get the list of songs from the repository
	songs, err := h.songRepo.GetList(sortBy, sortOrder, limitInt, offsetInt)
	if err != nil {
		h.logger.Error("Error getting songs from repository", err)
		http.Error(w, "Error getting songs", http.StatusInternalServerError)
		return
	}

	h.logger.Debug("Retrieved songs from repository",
		"count", len(songs),
	)

	// Marshal the songs to JSON
	songsJSON, err := json.Marshal(songs)
	if err != nil {
		h.logger.Error("Error marshaling songs to JSON", err)
		http.Error(w, "Error marshaling songs to JSON", http.StatusInternalServerError)
		return
	}

	h.logger.Debug("Successfully marshaled songs to JSON",
		"bytesLength", len(songsJSON),
	)

	// Write the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.Write(songsJSON)
}

// @Summary     Get song text
// @Description Get the text of a song by ID with pagination
// @Tags        songs
// @Accept      json
// @Produce     json
// @Param       id query int true "Song ID"
// @Param       limit query int false "Number of verses to return (default: 10)"
// @Param       offset query int false "Number of verses to skip (default: 0)"
// @Success     200 {object} models.Song
// @Failure     400 {object} models.ErrorResponse
// @Failure     404 {object} models.ErrorResponse
// @Router      /songs/text [get]
func (h *songHandlers) GetText(w http.ResponseWriter, r *http.Request) {
	// Get song ID from query parameters
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Song ID is required", http.StatusBadRequest)
		return
	}

	// Convert ID to integer
	songID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid song ID", http.StatusBadRequest)
		return
	}

	// Get pagination parameters
	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")

	if limit == "" {
		limit = defaultLimit
	}
	if offset == "" {
		offset = defaultOffset
	}

	// Convert pagination parameters to integers
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		http.Error(w, "Invalid limit value", http.StatusBadRequest)
		return
	}
	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		http.Error(w, "Invalid offset value", http.StatusBadRequest)
		return
	}

	// Get song text from repository
	text, err := h.songRepo.GetText(songID)
	if err != nil {
		h.logger.Error("Error getting song text", err)
		http.Error(w, "Song not found", http.StatusNotFound)
		return
	}

	if text == "" {
		http.Error(w, "Song text is empty", http.StatusNotFound)
		return
	}

	// Replace escaped newlines with actual newlines
	text = strings.ReplaceAll(text, "\\n", "\n")

	// Split the song text into lines
	lines := strings.Split(text, "\n")

	// Process lines to remove empty ones and verse markers
	var verses []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			verses = append(verses, line)
		}
	}

	// Validate offset
	if offsetInt >= len(verses) {
		http.Error(w, "Offset is out of range", http.StatusBadRequest)
		return
	}

	// Calculate end index for pagination
	endIndex := offsetInt + limitInt
	if endIndex > len(verses) {
		endIndex = len(verses)
	}

	// Paginate the verses
	paginatedVerses := verses[offsetInt:endIndex]

	// Create response structure
	response := struct {
		Total   int      `json:"total"`
		Verses  []string `json:"verses"`
		HasMore bool     `json:"has_more"`
	}{
		Total:   len(verses),
		Verses:  paginatedVerses,
		HasMore: endIndex < len(verses),
	}

	// Marshal the response to JSON
	responseJSON, err := json.Marshal(response)
	if err != nil {
		h.logger.Error("Error marshaling response", err)
		http.Error(w, "Error preparing response", http.StatusInternalServerError)
		return
	}

	// Write the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}

// @Summary     Delete song
// @Description Delete a song by ID
// @Tags        songs
// @Accept      json
// @Produce     plain
// @Param       id query int true "Song ID"
// @Success     200 {string} string "Song deleted successfully"
// @Failure     400 {object} models.ErrorResponse
// @Failure     404 {object} models.ErrorResponse
// @Router      /songs/ [delete]
func (h *songHandlers) Delete(w http.ResponseWriter, r *http.Request) {
	// Get the song ID from the URL parameters
	id := r.URL.Query().Get("id")

	// Validate the song ID
	if id == "" {
		http.Error(w, "Song ID is required", http.StatusBadRequest)
		return
	}

	// Convert the song ID to an integer
	songID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid song ID", http.StatusBadRequest)
		return
	}

	// Delete the song from the repository
	err = h.songRepo.Delete(songID)
	if err != nil {
		http.Error(w, "Error deleting song", http.StatusInternalServerError)
		return
	}

	// Return a success response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Song deleted successfully"))
}

// @Summary     Update song
// @Description Update song details by ID
// @Tags        songs
// @Accept      json
// @Produce     json
// @Param       id query int true "Song ID"
// @Param       request body models.UpdateSongRequest true "Song update request"
// @Success     200 {object} models.Song
// @Failure     400 {object} models.ErrorResponse
// @Failure     404 {object} models.ErrorResponse
// @Router      /songs/ [put]
func (h *songHandlers) Update(w http.ResponseWriter, r *http.Request) {
	// Get the song ID from the URL parameters
	id := r.URL.Query().Get("id")

	// Validate the song ID
	if id == "" {
		http.Error(w, "Song ID is required", http.StatusBadRequest)
		return
	}

	// Convert the song ID to an integer
	songID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid song ID", http.StatusBadRequest)
		return
	}

	// Get the updated song data from the request body
	var song models.Song
	err = json.NewDecoder(r.Body).Decode(&song)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Set the ID from the URL
	song.ID = songID

	// Update the song in the repository
	err = h.songRepo.Update(&song)
	if err != nil {
		h.logger.Error("Failed to update song", err)
		http.Error(w, fmt.Sprintf("Error updating song: %v", err), http.StatusInternalServerError)
		return
	}

	// Return a success response
	response := map[string]string{
		"message": "Song updated successfully",
		"id":      id,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// @Summary     Add new song
// @Description Add a new song to the library
// @Tags        songs
// @Accept      json
// @Produce     json
// @Param       request body models.AddSongRequest true "Song request"
// @Success     201 {object} models.Song
// @Failure     400 {object} models.ErrorResponse
// @Failure     500 {object} models.ErrorResponse
// @Router      /songs/ [post]
func (h *songHandlers) Add(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("Starting Add handler")

	var songRequest models.AddSongRequest
	if err := json.NewDecoder(r.Body).Decode(&songRequest); err != nil {
		h.logger.Error("Error decoding request body", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	h.logger.Debug("Received song request",
		"group", songRequest.Group,
		"song", songRequest.Song,
	)

	// Validate the request
	if songRequest.Group == "" {
		h.logger.Debug("Validation failed: empty Group")
		http.Error(w, "Group is required", http.StatusBadRequest)
		return
	}
	if songRequest.Song == "" {
		h.logger.Debug("Validation failed: empty Song")
		http.Error(w, "Song is required", http.StatusBadRequest)
		return
	}

	h.logger.Debug("Request validation passed")

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	// Prepare request to external API
	apiURL := fmt.Sprintf("%s/info?group=%s&song=%s",
		h.cfg.MusicApi.URL,
		url.QueryEscape(songRequest.Group),
		url.QueryEscape(songRequest.Song),
	)

	externalReq, err := http.NewRequestWithContext(r.Context(), "GET", apiURL, nil)
	if err != nil {
		h.logger.Error("Failed to create external API request", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Make request to external API
	resp, err := client.Do(externalReq)
	if err != nil {
		h.logger.Error("Failed to make external API request", "error", err, "url", apiURL)
		http.Error(w, "Failed to fetch song details", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode == http.StatusBadRequest {
		http.Error(w, "Invalid song or group name", http.StatusBadRequest)
		return
	}
	if resp.StatusCode != http.StatusOK {
		h.logger.Error("External API returned non-200 status",
			"status", resp.StatusCode,
			"url", apiURL,
		)
		http.Error(w, "Failed to fetch song details", http.StatusInternalServerError)
		return
	}

	// Parse external API response
	var songDetail models.SongDetail
	if err := json.NewDecoder(resp.Body).Decode(&songDetail); err != nil {
		h.logger.Error("Failed to decode external API response", "error", err)
		http.Error(w, "Failed to parse song details", http.StatusInternalServerError)
		return
	}

	// Validate required fields from external API
	if songDetail.ReleaseDate == "" || songDetail.Text == "" || songDetail.Link == "" {
		h.logger.Error("External API returned incomplete data",
			"releaseDate", songDetail.ReleaseDate,
			"hasText", songDetail.Text != "",
			"hasLink", songDetail.Link != "",
		)
		http.Error(w, "Incomplete song details received", http.StatusInternalServerError)
		return
	}

	// Create song entity with combined data
	song := &models.Song{
		Group:       songRequest.Group,
		Song:        songRequest.Song,
		ReleaseDate: songDetail.ReleaseDate,
		Text:        songDetail.Text,
		Link:        songDetail.Link,
	}

	// Save to database
	createdSong, err := h.songRepo.Create(r.Context(), song)
	if err != nil {
		h.logger.Error("Failed to create song", "error", err)
		http.Error(w, "Failed to create song", http.StatusInternalServerError)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(createdSong); err != nil {
		h.logger.Error("Failed to encode response", "error", err)
		return
	}
}

package http

import (
	"encoding/json"
	"musiclib/config"
	"musiclib/internal/models"
	"musiclib/internal/song"
	"musiclib/pkg/logger"
	"fmt"
	"net/http"
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

// Create godoc
// @Summary Get list
// @Description Get sorted list of songs
// @Tags Songs
// @Accept json
// @Produce json
// @Success 201 {array} models.Song
// @Router /list/ [get]
func (h *songHandlers) GetList(w http.ResponseWriter, r *http.Request) {
	// Get query parameters for sorting and pagination
	sortBy := r.URL.Query().Get("sort_by")
	sortOrder := r.URL.Query().Get("sort_order")
	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")

	h.logger.Info("GetList request parameters",
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

	// Get the list of songs from the repository
	songs, err := h.songRepo.GetList(sortBy, sortOrder, limitInt, offsetInt)
	if err != nil {
		h.logger.Error("Error getting songs from repository", err)
		http.Error(w, "Error getting songs", http.StatusInternalServerError)
		return
	}

	// Marshal the songs to JSON
	songsJSON, err := json.Marshal(songs)
	if err != nil {
		h.logger.Error("Error marshaling songs to JSON", err)
		http.Error(w, "Error marshaling songs to JSON", http.StatusInternalServerError)
		return
	}

	// Write the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.Write(songsJSON)
}

// Create godoc
// @Summary Get song text
// @Description Get the text of a song by song ID
// @Tags songs
// @Accept  json
// @Produce  json
// @Param id query string true "Song ID"
// @Param limit query string false "Number of lyrics to return"
// @Param offset query string false "Offset from which to start returning lyrics"
// @Success 201 {array} models.Lyric
// @Router /text/ [get]
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

// DeleteSongHandler handles the deletion of a song
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

// UpdateSongHandler handles the updating of a song
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

// Add godoc
// @Summary Add new song
// @Description Add a new song to the library with details from external API
// @Tags Songs
// @Accept json
// @Produce json
// @Param request body models.AddSongRequest true "Song request"
// @Success 201 {object} models.Song
// @Router / [post]
func (h *songHandlers) Add(w http.ResponseWriter, r *http.Request) {
	var request models.AddSongRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.logger.Error("failed to decode request body", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if request.Group == "" || request.Song == "" {
		http.Error(w, "Group and Song are required fields", http.StatusBadRequest)
		return
	}

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	// Prepare request to external API
	externalReq, err := http.NewRequestWithContext(
		r.Context(),
		"GET",
		h.cfg.MusicApi+"/info",
		nil,
	)
	if err != nil {
		h.logger.Error("failed to create external API request", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Add query parameters
	q := externalReq.URL.Query()
	q.Add("group", request.Group)
	q.Add("song", request.Song)
	externalReq.URL.RawQuery = q.Encode()

	// Make request to external API
	resp, err := client.Do(externalReq)
	if err != nil {
		h.logger.Error("failed to make external API request", "error", err)
		http.Error(w, "Failed to fetch song details", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		h.logger.Error("external API returned non-200 status", "status", resp.StatusCode)
		http.Error(w, "Failed to fetch song details", http.StatusInternalServerError)
		return
	}

	// Parse external API response
	var songDetail models.SongDetail
	if err := json.NewDecoder(resp.Body).Decode(&songDetail); err != nil {
		h.logger.Error("failed to decode external API response", "error", err)
		http.Error(w, "Failed to parse song details", http.StatusInternalServerError)
		return
	}

	// Create song entity with combined data
	song := &models.Song{
		Group:       request.Group,
		Song:        request.Song,
		ReleaseDate: songDetail.ReleaseDate,
		Text:        songDetail.Text,
		Link:        songDetail.Link,
	}

	// Save to database
	createdSong, err := h.songRepo.Create(r.Context(), song)
	if err != nil {
		h.logger.Error("failed to create song", "error", err)
		http.Error(w, "Failed to create song", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(createdSong); err != nil {
		h.logger.Error("failed to encode response", "error", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

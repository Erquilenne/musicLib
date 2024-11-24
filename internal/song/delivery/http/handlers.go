package http

import (
	"encoding/json"
	"musiclib/config"
	"musiclib/internal/models"
	"musiclib/internal/song"
	"musiclib/pkg/logger"
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
func NewSongHandlers(cfg *config.Config, logger logger.Logger) *songHandlers {
	return &songHandlers{cfg: cfg, logger: logger}
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
		http.Error(w, "Invalid limit value", http.StatusBadRequest)
		return
	}
	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		http.Error(w, "Invalid offset value", http.StatusBadRequest)
		return
	}

	// Get the list of songs from the repository
	songs, err := h.songRepo.GetList(sortBy, sortOrder, limitInt, offsetInt)
	if err != nil {
		http.Error(w, "Error getting songs", http.StatusInternalServerError)
		return
	}

	// Marshal the songs to JSON
	songsJSON, err := json.Marshal(songs)
	if err != nil {
		http.Error(w, "Error marshaling songs to JSON", http.StatusInternalServerError)
		return
	}

	// Write the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.Write(songsJSON)
}

// Create godoc
// @Summary Get song text
// @Description Get the text of a song by song name and group name
// @Tags songs
// @Accept  json
// @Produce  json
// @Param song query string true "Song name"
// @Param group query string true "Group name"
// @Param limit query string false "Number of lyrics to return"
// @Param offset query string false "Offset from which to start returning lyrics"
// @Success 201 {array} models.Lyric
// @Router /text/ [get]
func (h *songHandlers) GetText(w http.ResponseWriter, r *http.Request) {
	// Get query parameters for filtering and pagination
	song := r.URL.Query().Get("song")
	group := r.URL.Query().Get("group")
	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")

	// Validate required parameters
	if song == "" {
		http.Error(w, "Song name is required", http.StatusBadRequest)
		return
	}
	if group == "" {
		http.Error(w, "Group name is required", http.StatusBadRequest)
		return
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
		http.Error(w, "Invalid limit value", http.StatusBadRequest)
		return
	}
	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		http.Error(w, "Invalid offset value", http.StatusBadRequest)
		return
	}

	// Get song text from repository
	text, err := h.songRepo.GetText(song, group)
	if err != nil {
		http.Error(w, "Error getting song", http.StatusInternalServerError)
		return
	}

	// Split the song text into verses
	verses := strings.Split(text, "\n\n")

	// Paginate the verses
	paginatedVerses := verses[offsetInt : offsetInt+limitInt]

	// Marshal the paginated verses to JSON
	versesJSON, err := json.Marshal(paginatedVerses)
	if err != nil {
		http.Error(w, "Error marshaling verses to JSON", http.StatusInternalServerError)
		return
	}

	// Write the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.Write(versesJSON)
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
	song.ID = songID
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Update the song in the repository
	err = h.songRepo.Update(&song)
	if err != nil {
		http.Error(w, "Error updating song", http.StatusInternalServerError)
		return
	}

	// Return a success response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Song updated successfully"))
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
		h.cfg.Server.ExternalMusicAPI+"/info",
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

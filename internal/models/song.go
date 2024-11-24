package models

// Song represents a song in the database
type Song struct {
	ID          int    `json:"id" db:"id" example:"1"`
	Group       string `json:"group" db:"group_name" example:"Beatles"`
	Song        string `json:"song" db:"song" example:"Yesterday"`
	ReleaseDate string `json:"releaseDate" db:"release_date" example:"1965-09-13"`
	Text        string `json:"text" db:"text" example:"Yesterday all my troubles seemed so far away..."`
	Link        string `json:"link" db:"link" example:"https://example.com/song"`
}

// AddSongRequest represents the incoming request structure for adding a new song
type AddSongRequest struct {
	Group string `json:"group" example:"Beatles"`
	Song  string `json:"song" example:"Yesterday"`
}

// UpdateSongRequest represents the incoming request structure for updating a song
type UpdateSongRequest struct {
	Group       string `json:"group,omitempty" example:"Beatles"`
	Song        string `json:"song,omitempty" example:"Yesterday"`
	ReleaseDate string `json:"releaseDate,omitempty" example:"1965-09-13"`
	Text        string `json:"text,omitempty" example:"Yesterday all my troubles seemed so far away..."`
	Link        string `json:"link,omitempty" example:"https://example.com/song"`
}

// SongDetail represents the external API response structure
type SongDetail struct {
	ReleaseDate string `json:"releaseDate" example:"1965-09-13"`
	Text        string `json:"text" example:"Yesterday all my troubles seemed so far away..."`
	Link        string `json:"link" example:"https://example.com/song"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Message string `json:"message" example:"error message"`
}

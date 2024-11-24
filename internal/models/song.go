package models

type Song struct {
	ID          int    `json:"id" db:"id"`
	Group       string `json:"group" db:"group_name"`
	Song        string `json:"song" db:"song"`
	ReleaseDate string `json:"releaseDate" db:"release_date"`
	Text        string `json:"text" db:"text"`
	Link        string `json:"link" db:"link"`
}

// AddSongRequest represents the incoming request structure for adding a new song
type AddSongRequest struct {
	Group string `json:"group"`
	Song  string `json:"song"`
}

// SongDetail represents the external API response structure
type SongDetail struct {
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

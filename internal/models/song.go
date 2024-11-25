package models

type Song struct {
	ID          int    `json:"id" db:"id" example:"1"`
	Group       string `json:"group" db:"group_name" example:"Beatles"`
	Song        string `json:"song" db:"song" example:"Yesterday"`
	ReleaseDate string `json:"releaseDate" db:"release_date" example:"1965-09-13"`
	Text        string `json:"text" db:"text" example:"Yesterday all my troubles seemed so far away..."`
	Link        string `json:"link" db:"link" example:"https://example.com/song"`
}

type AddSongRequest struct {
	Group string `json:"group" example:"Beatles"`
	Song  string `json:"song" example:"Yesterday"`
}

type UpdateSongRequest struct {
	Group       string `json:"group,omitempty" example:"Beatles"`
	Song        string `json:"song,omitempty" example:"Yesterday"`
	ReleaseDate string `json:"releaseDate,omitempty" example:"1965-09-13"`
	Text        string `json:"text,omitempty" example:"Yesterday all my troubles seemed so far away..."`
	Link        string `json:"link,omitempty" example:"https://example.com/song"`
}

type SongDetail struct {
	ReleaseDate string `json:"releaseDate" example:"1965-09-13"`
	Text        string `json:"text" example:"Yesterday all my troubles seemed so far away..."`
	Link        string `json:"link" example:"https://example.com/song"`
}

type ErrorResponse struct {
	Message string `json:"message" example:"error message"`
}

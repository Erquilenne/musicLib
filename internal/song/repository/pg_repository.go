package repository

import (
	"context"
	"fmt"
	"musiclib/internal/models"

	"github.com/jmoiron/sqlx"
)

type songRepository struct {
	db *sqlx.DB
}

func NewPGRepository(db *sqlx.DB) *songRepository {
	return &songRepository{db: db}
}

func (r *songRepository) GetList(sortBy string, sortOrder string, limit int, offset int) ([]models.Song, error) {
	query := fmt.Sprintf(getList, sortBy, sortOrder, limit, offset)

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	songs := make([]models.Song, 0)

	for rows.Next() {
		var song models.Song
		err := rows.Scan(&song.ID, &song.Group, &song.Song, &song.Text, &song.Link)
		if err != nil {
			return nil, err
		}
		songs = append(songs, song)
	}

	return songs, nil
}

func (r *songRepository) GetText(song string, group string) (string, error) {

	// Execute the query and retrieve the song
	row := r.db.QueryRow(getText, song, group)
	var text string
	err := row.Scan(&text)
	if err != nil {
		return "", err
	}

	return text, nil
}

func (r *songRepository) Delete(id int) error {
	_, err := r.db.Exec(deleteSong, id)
	return err
}

func (r *songRepository) Update(song *models.Song) error {
	_, err := r.db.Exec(updateSong, song.Group, song.Song, song.Text, song.Link, song.ID)
	return err
}

func (r *songRepository) Create(ctx context.Context, song *models.Song) (*models.Song, error) {
	var id int
	err := r.db.QueryRowContext(ctx, createSong, song.Group, song.Song, song.ReleaseDate, song.Text, song.Link).Scan(&id)
	if err != nil {
		return nil, err
	}

	song.ID = id
	return song, nil
}

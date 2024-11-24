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

func NewSongRepository(db *sqlx.DB) *songRepository {
	return &songRepository{db: db}
}

func (r *songRepository) GetList(sortBy string, sortOrder string, limit int, offset int) ([]models.Song, error) {
	// Формируем полный запрос
	query := getList

	// Добавляем сортировку
	orderBy := " ORDER BY "
	if sortBy == "" {
		orderBy += "id"
	} else {
		orderBy += sortBy
	}

	// Добавляем направление сортировки
	if sortOrder == "desc" {
		orderBy += " DESC"
	} else {
		orderBy += " ASC"
	}

	// Добавляем пагинацию
	query += orderBy + " LIMIT $1 OFFSET $2"

	// Выполняем запрос
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get songs list: %v", err)
	}
	defer rows.Close()

	songs := make([]models.Song, 0)

	for rows.Next() {
		var song models.Song
		err := rows.Scan(&song.ID, &song.Group, &song.Song, &song.Text, &song.Link)
		if err != nil {
			return nil, fmt.Errorf("failed to scan song: %v", err)
		}
		songs = append(songs, song)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating songs: %v", err)
	}

	return songs, nil
}

func (r *songRepository) GetText(id int) (string, error) {
	row := r.db.QueryRow(getText, id)
	var text string
	err := row.Scan(&text)
	if err != nil {
		return "", fmt.Errorf("failed to get song text: %v", err)
	}
	return text, nil
}

func (r *songRepository) Delete(id int) error {
	_, err := r.db.Exec(deleteSong, id)
	return err
}

func (r *songRepository) Update(song *models.Song) error {
	var id int
	err := r.db.QueryRow(updateSong, song.Group, song.Song, song.Text, song.Link, song.ReleaseDate, song.ID).Scan(&id)
	if err != nil {
		return fmt.Errorf("failed to update song: %v", err)
	}
	return nil
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

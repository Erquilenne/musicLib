package repository

import (
	"context"
	"fmt"
	"musiclib/internal/models"
	"musiclib/pkg/logger"

	"github.com/jmoiron/sqlx"
)

type songRepository struct {
	db     *sqlx.DB
	logger logger.Logger
}

func NewSongRepository(db *sqlx.DB, logger logger.Logger) *songRepository {
	return &songRepository{
		db:     db,
		logger: logger,
	}
}

func (r *songRepository) GetList(sortBy string, sortOrder string, limit int, offset int) ([]models.Song, error) {
	r.logger.Debug("Starting GetList in repository",
		"sortBy", sortBy,
		"sortOrder", sortOrder,
		"limit", limit,
		"offset", offset,
	)

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

	r.logger.Debug("Executing SQL query", "query", query)

	// Выполняем запрос
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		r.logger.Debug("Failed to execute query", "error", err)
		return nil, fmt.Errorf("failed to get songs list: %v", err)
	}
	defer rows.Close()

	songs := make([]models.Song, 0)

	for rows.Next() {
		var song models.Song
		err := rows.Scan(&song.ID, &song.Group, &song.Song, &song.Text, &song.Link)
		if err != nil {
			r.logger.Debug("Failed to scan row", "error", err)
			return nil, fmt.Errorf("failed to scan song: %v", err)
		}
		songs = append(songs, song)
	}

	if err = rows.Err(); err != nil {
		r.logger.Debug("Error iterating rows", "error", err)
		return nil, fmt.Errorf("error iterating songs: %v", err)
	}

	r.logger.Debug("Successfully retrieved songs", "count", len(songs))
	return songs, nil
}

func (r *songRepository) GetText(id int) (string, error) {
	r.logger.Debug("Starting GetText in repository", "id", id)

	row := r.db.QueryRow(getText, id)
	var text string
	err := row.Scan(&text)
	if err != nil {
		r.logger.Debug("Failed to get song text", "error", err, "id", id)
		return "", fmt.Errorf("failed to get song text: %v", err)
	}

	r.logger.Debug("Successfully retrieved song text", 
		"id", id,
		"textLength", len(text),
	)
	return text, nil
}

func (r *songRepository) Delete(id int) error {
	r.logger.Debug("Starting Delete in repository", "id", id)

	result, err := r.db.Exec(deleteSong, id)
	if err != nil {
		r.logger.Debug("Failed to delete song", "error", err, "id", id)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.Debug("Failed to get rows affected", "error", err)
		return err
	}

	r.logger.Debug("Successfully deleted song", 
		"id", id,
		"rowsAffected", rowsAffected,
	)
	return nil
}

func (r *songRepository) Update(song *models.Song) error {
	r.logger.Debug("Starting Update in repository", 
		"id", song.ID,
		"group", song.Group,
		"song", song.Song,
	)

	var id int
	err := r.db.QueryRow(updateSong, 
		song.Group, 
		song.Song, 
		song.Text, 
		song.Link, 
		song.ReleaseDate, 
		song.ID,
	).Scan(&id)

	if err != nil {
		r.logger.Debug("Failed to update song", 
			"error", err,
			"id", song.ID,
		)
		return fmt.Errorf("failed to update song: %v", err)
	}

	r.logger.Debug("Successfully updated song", "id", id)
	return nil
}

func (r *songRepository) Create(ctx context.Context, song *models.Song) (*models.Song, error) {
	r.logger.Debug("Starting Create in repository", 
		"group", song.Group,
		"song", song.Song,
	)

	var id int
	err := r.db.QueryRowContext(ctx, createSong,
		song.Group,
		song.Song,
		song.ReleaseDate,
		song.Text,
		song.Link,
	).Scan(&id)

	if err != nil {
		r.logger.Debug("Failed to create song", "error", err)
		return nil, err
	}

	song.ID = id
	r.logger.Debug("Successfully created song", "id", id)
	return song, nil
}

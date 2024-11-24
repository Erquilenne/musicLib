package song

import (
	"context"
	"musiclib/internal/models"
)

// Repository interface
type Repository interface {
	GetList(sortBy string, sortOrder string, limit int, offset int) ([]models.Song, error)
	GetText(id int) (string, error)
	Delete(id int) error
	Update(song *models.Song) error
	Create(ctx context.Context, song *models.Song) (*models.Song, error)
}

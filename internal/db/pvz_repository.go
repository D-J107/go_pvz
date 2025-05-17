package db

import (
	"context"
	"my_pvz/internal/domain/models"
	"time"
)

type PvzRepository interface {
	Create(ctx context.Context, id, registrationDate, city string) (*models.Pvz, error)
	GetAllWithFilter(ctx context.Context, startDate, endDate time.Time, page, limit int) ([]models.Pvz, error)
	GetAll(ctx context.Context) ([]models.Pvz, error)
}

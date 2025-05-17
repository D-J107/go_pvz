package db

import (
	"context"
	"my_pvz/internal/domain/models"
)

type ReceptionRepository interface {
	Create(ctx context.Context, pvzId, status string) (*models.Reception, error)
	GetAllByPvzID(ctx context.Context, pvzId string) ([]*models.Reception, error)
	GetLastByPvzID(ctx context.Context, pvzId string) (*models.Reception, error)
	UpdateReceptionStatus(ctx context.Context, receptionId, newStatus string) error
}

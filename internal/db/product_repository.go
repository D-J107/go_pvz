package db

import (
	"context"
	"my_pvz/internal/domain/models"
)

type ProductRepository interface {
	Create(ctx context.Context, productType, receptionId string) (*models.Product, error)
	GetAllByReceptionId(ctx context.Context, receptionId string) ([]*models.Product, error)
	DeleteLastProductInReception(ctx context.Context, receptionId string) error
}

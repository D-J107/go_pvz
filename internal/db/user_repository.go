package db

import (
	"context"
	"my_pvz/internal/domain/models"
)

type UserRepository interface {
	Create(ctx context.Context, email, hashedPwd, role string) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
}

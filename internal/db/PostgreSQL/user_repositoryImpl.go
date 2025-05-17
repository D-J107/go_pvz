package postgresql

import (
	"context"
	"my_pvz/internal/db"
	"my_pvz/internal/db/sq"
	"my_pvz/internal/domain/models"

	"github.com/Masterminds/squirrel"
)

type UserRepositoryImpl struct {
	db *db.DB
}

func NewPostgesUserRepositoryImpl(db *db.DB) *UserRepositoryImpl {
	return &UserRepositoryImpl{db: db}
}

func (rep *UserRepositoryImpl) Create(ctx context.Context, email, hashedPwd, role string) (*models.User, error) {
	sql, args, err := sq.Psql.
		Insert("users").
		Columns("email", "password", "role").
		Values(email, hashedPwd, role).
		Suffix("RETURNING id, email, password, role").
		ToSql()
	if err != nil {
		return nil, err
	}
	user := &models.User{}
	if err = rep.db.Pool.QueryRow(ctx, sql, args...).Scan(&user.ID, &user.Email, &user.Password, &user.Role); err != nil {
		return nil, err
	}
	return user, nil
}

func (rep *UserRepositoryImpl) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	sql, args, err := sq.Psql.
		Select("id, email, password, role").
		From("users").
		Where(squirrel.Eq{"email": email}).
		ToSql()
	if err != nil {
		return nil, err
	}
	user := &models.User{}
	rep.db.Pool.QueryRow(ctx, sql, args...).Scan(&user.ID, &user.Email, &user.Password, &user.Role)
	return user, nil
}

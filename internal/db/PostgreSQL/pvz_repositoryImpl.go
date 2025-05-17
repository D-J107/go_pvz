package postgresql

import (
	"context"
	"fmt"
	"my_pvz/internal/db"
	"my_pvz/internal/db/sq"
	"my_pvz/internal/domain/models"
	"time"

	"github.com/Masterminds/squirrel"
)

type PostgresPvzRepositoryImpl struct {
	db *db.DB
}

func NewPostgresPvzRepositoryImpl(db *db.DB) *PostgresPvzRepositoryImpl {
	return &PostgresPvzRepositoryImpl{db: db}
}

func (rep *PostgresPvzRepositoryImpl) Create(ctx context.Context, id, registrationDate, city string) (*models.Pvz, error) {
	sql, args, err := sq.Psql.
		Insert("pvz").
		Columns("id", "registration_date_time", "city").
		Values(id, registrationDate, city).
		Suffix("RETURNING id, registration_date_time, city").
		ToSql()
	if err != nil {
		return nil, err
	}
	pvz := &models.Pvz{}
	if err := rep.db.Pool.QueryRow(ctx, sql, args...).Scan(&pvz.ID, &pvz.RegistrationDate, &pvz.City); err != nil {
		return nil, err
	}
	return pvz, nil
}

func (rep *PostgresPvzRepositoryImpl) GetAllWithFilter(ctx context.Context, startDate, endDate time.Time, page, limit int) ([]models.Pvz, error) {
	offset := page - 1
	sql, args, err := sq.Psql.
		Select("id", "registration_date_time", "city").
		From("pvz").
		Where(squirrel.GtOrEq{"registration_date_time": startDate}).
		Where(squirrel.LtOrEq{"registration_date_time": endDate}).
		Limit(uint64(limit)).
		Offset(uint64(offset)).
		ToSql()
	if err != nil {
		return nil, err
	}
	fmt.Println("sql:", sql)
	fmt.Println("args:", args)
	rows, err := rep.db.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Pvz
	for rows.Next() {
		var p models.Pvz
		if err := rows.Scan(&p.ID, &p.RegistrationDate, &p.City); err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}

func (rep *PostgresPvzRepositoryImpl) GetAll(ctx context.Context) ([]models.Pvz, error) {
	sql, args, err := sq.Psql.
		Select("id", "registration_date_time", "city").
		From("pvz").
		ToSql()
	if err != nil {
		return nil, err
	}
	rows, err := rep.db.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	asnwer := make([]models.Pvz, 0)
	for rows.Next() {
		pvz := models.Pvz{}
		if err = rows.Scan(&pvz.ID, &pvz.RegistrationDate, &pvz.City); err != nil {
			return nil, err
		}
		asnwer = append(asnwer, pvz)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return asnwer, nil
}

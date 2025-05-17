package postgresql

import (
	"context"
	"my_pvz/internal/db"
	"my_pvz/internal/db/sq"
	"my_pvz/internal/domain/models"

	"github.com/Masterminds/squirrel"
)

type ReceptionRepositoryImpl struct {
	db *db.DB
}

func NewPostgesReceptionRepositoryImpl(db *db.DB) *ReceptionRepositoryImpl {
	return &ReceptionRepositoryImpl{db: db}
}

func (r *ReceptionRepositoryImpl) Create(ctx context.Context, pvzId, status string) (*models.Reception, error) {
	sql, args, err := sq.Psql.
		Insert("receptions").
		Columns("pvz_id", "status").
		Values(pvzId, status).
		Suffix("RETURNING id, date_time, pvz_id, status").
		ToSql()
	if err != nil {
		return nil, err
	}

	rec := &models.Reception{}
	if err = r.db.Pool.QueryRow(ctx, sql, args...).Scan(&rec.ID, &rec.DateTime, &rec.PvzId, &rec.Status); err != nil {
		return nil, err
	}
	return rec, nil

}

func (r *ReceptionRepositoryImpl) GetAllByPvzID(ctx context.Context, pvzId string) ([]*models.Reception, error) {
	sql, args, err := sq.Psql.
		Select("id", "date_time", "status").
		From("receptions").
		Where(squirrel.Eq{"pvz_id": pvzId}).
		ToSql()
	if err != nil {
		return nil, err
	}
	rows, err := r.db.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recs []*models.Reception
	for rows.Next() {
		rec := &models.Reception{}
		if err := rows.Scan(&rec.ID, &rec.DateTime, &rec.Status); err != nil {
			return nil, err
		}
		rec.PvzId = pvzId
		recs = append(recs, rec)
	}
	return recs, rows.Err()
}

func (r *ReceptionRepositoryImpl) GetLastByPvzID(ctx context.Context, pvzId string) (*models.Reception, error) {
	sql, args, err := sq.Psql.
		Select("id", "date_time", "status").
		From("receptions").
		Where(squirrel.Eq{"pvz_id": pvzId}).
		OrderBy("date_time DESC").
		Limit(1).
		ToSql()
	if err != nil {
		return nil, err
	}

	rec := &models.Reception{}
	if err = r.db.Pool.QueryRow(ctx, sql, args...).Scan(&rec.ID, &rec.DateTime, &rec.Status); err != nil {
		return nil, err
	}
	rec.PvzId = pvzId
	return rec, nil
}

func (r *ReceptionRepositoryImpl) UpdateReceptionStatus(ctx context.Context, receptionId, newStatus string) error {
	sql, args, err := sq.Psql.
		Update("receptions").
		Set("status", newStatus).
		Where(squirrel.Eq{"id": receptionId}).
		ToSql()
	if err != nil {
		return err
	}
	_, err = r.db.Pool.Exec(ctx, sql, args...)
	return err
}

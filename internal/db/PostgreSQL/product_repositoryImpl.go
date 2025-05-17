package postgresql

import (
	"context"
	"errors"
	"my_pvz/internal/db"
	"my_pvz/internal/db/sq"
	"my_pvz/internal/domain/models"

	"github.com/Masterminds/squirrel"
)

type ProductRepositoryImpl struct {
	db *db.DB
}

func NewPosgresProductRepositoryImpl(db *db.DB) *ProductRepositoryImpl {
	return &ProductRepositoryImpl{db: db}
}

func (rep *ProductRepositoryImpl) Create(ctx context.Context, productType, receptionId string) (*models.Product, error) {
	sql, args, err := sq.Psql.
		Insert("products").
		Columns("type, reception_id").
		Values(productType, receptionId).
		Suffix("RETURNING id, date_time, type").
		ToSql()
	if err != nil {
		return nil, err
	}

	product := &models.Product{}
	if err := rep.db.Pool.QueryRow(ctx, sql, args...).Scan(&product.ID, &product.DateTime, &product.Type); err != nil {
		return nil, err
	}
	product.ReceptionId = receptionId

	return product, nil
}

func (rep *ProductRepositoryImpl) GetAllByReceptionId(ctx context.Context, receptionId string) ([]*models.Product, error) {
	sql, args, err := sq.Psql.
		Select("id", "date_time", "type", "reception_id").
		From("products").
		Where(squirrel.Eq{"reception_id": receptionId}).
		ToSql()

	if err != nil {
		return nil, err
	}
	rows, err := rep.db.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	output := make([]*models.Product, 0)
	for rows.Next() {
		p := &models.Product{}
		if err = rows.Scan(&p.ID, &p.DateTime, &p.Type, &p.ReceptionId); err != nil {
			return nil, err
		}
		output = append(output, p)
	}
	return output, nil
}

func (rep *ProductRepositoryImpl) DeleteLastProductInReception(ctx context.Context, receptionId string) error {
	query := `DELETE FROM products WHERE id = (SELECT id FROM products WHERE reception_id = $1 ORDER BY date_time DESC LIMIT 1)`
	cmdTag, err := rep.db.Pool.Exec(ctx, query, receptionId)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		// ничего не удалили <==> в текущей приемке не было товаров
		return errors.New("active reception run out of products to be deleted!")
	}
	return nil
}

package models

import "time"

type Product struct {
	ID          string    `db:"id"`
	DateTime    time.Time `db:"date_time"`
	Type        string    `db:"type"`
	ReceptionId string    `db:"reception_id"`
}

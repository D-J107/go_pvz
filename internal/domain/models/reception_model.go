package models

import "time"

type Reception struct {
	ID       string    `db:"id"`
	DateTime time.Time `db:"date_time"`
	PvzId    string    `db:"pvz_id"` // foreign key to pvz table
	Status   string    `db:"status"`
}

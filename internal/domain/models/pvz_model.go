package models

import "time"

type Pvz struct {
	ID               string    `db:"db"`
	RegistrationDate time.Time `db:"registration_date_time"`
	City             string    `db:"city"`
}

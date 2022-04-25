package models

import "time"

type GormModel struct {
	id        int `gorm:"primary_key,em"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

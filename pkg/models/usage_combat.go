package models

import "time"

type UsageCombat struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	Key       string `gorm:"primaryKey"`
	Version   int
	Results   string
}

package models

import "time"

type UsageCombat struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	Type      string `gorm:"primaryKey"`
	Key       string `gorm:"primaryKey"`
	Version   int
	Results   string
}

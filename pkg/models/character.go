package models

import (
	"time"

	"gorm.io/gorm"
)

type Character struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Name      string         `gorm:"primaryKey"`
	Data      string
}

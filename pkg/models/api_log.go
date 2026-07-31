package models

import "gorm.io/gorm"

type APILog struct {
	gorm.Model
	Method   string `gorm:"index"`
	Path     string `gorm:"index"`
	Body     string
	Response string
	Status   int `gorm:"index"`
	Cooldown int
}

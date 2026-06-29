package models

import "gorm.io/gorm"

type APILog struct {
	gorm.Model
	Method   string
	Path     string
	Body     string
	Response string
	Status   int
	Cooldown int
}

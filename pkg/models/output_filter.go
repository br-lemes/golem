package models

type OutputFilter struct {
	Command string `gorm:"primaryKey"`
	Kind    string `gorm:"primaryKey"`
	Pattern string `gorm:"primaryKey"`
}

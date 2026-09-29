package models

type OutputFilter struct {
	Command string `gorm:"primaryKey" json:"command"`
	Kind    string `gorm:"primaryKey" json:"kind"`
	Pattern string `gorm:"primaryKey" json:"pattern"`
}

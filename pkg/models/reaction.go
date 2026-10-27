package models

type Reaction struct {
	Name     string `gorm:"primaryKey;not null;index:,unique,composite:priority" json:"name"`
	Code     string `gorm:"primaryKey;not null" json:"code"`
	Priority int    `gorm:"index:,unique,composite:priority" json:"priority"`
	Enabled  bool   `json:"enabled"`
}

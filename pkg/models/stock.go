package models

type StockMode string

const (
	StockModeExact  StockMode = "exact"
	StockModeSafety StockMode = "safety"
)

type Stock struct {
	Code     string    `gorm:"primaryKey" json:"code"`
	Quantity int       `json:"quantity"`
	Mode     StockMode `json:"mode"`
	Enabled  bool      `json:"enabled"`
	Discard  bool      `json:"discard"`
}

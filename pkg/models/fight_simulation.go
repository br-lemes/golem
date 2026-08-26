package models

import "time"

type FightSimulation struct {
	CreatedAt             time.Time
	UpdatedAt             time.Time
	Name                  string `gorm:"primaryKey"`
	Version               int
	Winrate               float32
	AverageTurns          float32
	AverageFinalHP        float32
	AverageFightCooldown  float32
	EstimatedRestCooldown float32
	CycleCost             float32
	DamageSurplus         float32
	SurvivalSurplus       float32
	XP                    int
	XPPerSecond           float32
	GoldPerSecond         float32
	ProspectingEfficiency float32
	Safe                  bool
}

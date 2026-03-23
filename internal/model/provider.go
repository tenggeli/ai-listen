package model

import (
	"time"

	"gorm.io/gorm"
)

type Provider struct {
	ID             uint64         `gorm:"primaryKey;autoIncrement"`
	UserID         uint64         `gorm:"not null;uniqueIndex:uk_user_id"`
	ProviderNo     string         `gorm:"size:32;not null;uniqueIndex:uk_provider_no"`
	Zodiac         string         `gorm:"size:16"`
	Constellation  string         `gorm:"size:16"`
	Level          string         `gorm:"size:16"`
	Score          float64        `gorm:"type:decimal(4,2);default:0"`
	TotalOrders    int            `gorm:"default:0"`
	TotalIncome    float64        `gorm:"type:decimal(12,2);default:0"`
	ServiceStatus  uint8          `gorm:"type:tinyint;default:3"`
	OnlineStatus   uint8          `gorm:"type:tinyint;default:0"`
	CityID         uint64         `gorm:"index:idx_city_id"`
	Intro          string         `gorm:"size:500"`
	Tags           string         `gorm:"type:json"`
	CommissionRate float64        `gorm:"type:decimal(5,2);default:0"`
	ComplaintCount int            `gorm:"default:0"`
	CancelCount    int            `gorm:"default:0"`
	CreatedAt      time.Time      `gorm:"type:datetime;not null"`
	UpdatedAt      time.Time      `gorm:"type:datetime;not null"`
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

func (Provider) TableName() string {
	return "providers"
}

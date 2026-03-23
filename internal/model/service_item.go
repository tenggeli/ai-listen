package model

import (
	"time"

	"gorm.io/gorm"
)

type ServiceItem struct {
	ID          uint64         `gorm:"primaryKey;autoIncrement"`
	ProviderID  uint64         `gorm:"not null;index:idx_provider_id_status,priority:1;index:idx_provider_id"`
	CategoryID  uint64         `gorm:"index:idx_category_id"`
	Title       string         `gorm:"size:128;not null"`
	Description string         `gorm:"size:500"`
	UnitPrice   float64        `gorm:"type:decimal(10,2);default:0"`
	BillingType uint8          `gorm:"type:tinyint;default:1"`
	MinHours    int            `gorm:"default:0"`
	MaxHours    int            `gorm:"default:0"`
	UnitName    string         `gorm:"size:20"`
	Status      uint8          `gorm:"type:tinyint;default:1;index:idx_provider_id_status,priority:2;index:idx_status"`
	CreatedAt   time.Time      `gorm:"type:datetime;not null"`
	UpdatedAt   time.Time      `gorm:"type:datetime;not null"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (ServiceItem) TableName() string {
	return "service_items"
}

package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           uint64         `gorm:"primaryKey;autoIncrement"`
	Mobile       string         `gorm:"size:20;not null;uniqueIndex:idx_mobile"`
	WechatOpenID string         `gorm:"size:64;index:idx_wechat_open_id"`
	UnionID      string         `gorm:"size:64"`
	Nickname     string         `gorm:"size:64"`
	Avatar       string         `gorm:"size:255"`
	Gender       uint8          `gorm:"type:tinyint;default:0"`
	Birthday     *time.Time     `gorm:"type:date"`
	CityID       uint64         `gorm:"index:idx_city_id"`
	Bio          string         `gorm:"size:255"`
	VipLevel     uint8          `gorm:"type:tinyint;default:0"`
	UserStatus   uint8          `gorm:"type:tinyint;default:1"`
	LastLoginAt  *time.Time     `gorm:"type:datetime"`
	CreatedAt    time.Time      `gorm:"type:datetime;not null"`
	UpdatedAt    time.Time      `gorm:"type:datetime;not null"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

func (User) TableName() string {
	return "users"
}

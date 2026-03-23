package model

import "time"

type ProviderApplication struct {
	ID               uint64     `gorm:"primaryKey;autoIncrement"`
	UserID           uint64     `gorm:"not null;index:idx_user_id"`
	RealName         string     `gorm:"size:64;not null"`
	IDCardNo         string     `gorm:"size:128;not null"`
	IDCardFront      string     `gorm:"size:255;not null"`
	IDCardBack       string     `gorm:"size:255;not null"`
	FaceVerifyStatus uint8      `gorm:"type:tinyint;default:0"`
	AgreementSigned  uint8      `gorm:"type:tinyint;default:0"`
	CityID           uint64     `gorm:"index:idx_city_id"`
	Intro            string     `gorm:"size:500"`
	Photos           string     `gorm:"type:json"`
	ServiceDesc      string     `gorm:"size:500"`
	AuditStatus      uint8      `gorm:"type:tinyint;default:0"`
	RejectReason     string     `gorm:"size:255"`
	SubmittedAt      *time.Time `gorm:"type:datetime"`
	AuditedAt        *time.Time `gorm:"type:datetime"`
	CreatedAt        time.Time  `gorm:"type:datetime;not null"`
	UpdatedAt        time.Time  `gorm:"type:datetime;not null"`
}

func (ProviderApplication) TableName() string {
	return "provider_applications"
}

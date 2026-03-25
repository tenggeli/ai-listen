package model

import "time"

type ProviderServiceItem struct {
	ServiceItemID uint64 `json:"serviceItemId"`
	PriceAmount   int64  `json:"priceAmount"`
	PriceUnit     string `json:"priceUnit"`
}

type Provider struct {
	ID           uint64                `json:"id"`
	UserID       uint64                `json:"userId"`
	ProviderNo   string                `json:"providerNo"`
	RealName     string                `json:"realName"`
	IDCardNo     string                `json:"idCardNo"`
	AuditStatus  int                   `json:"auditStatus"`
	AuditRemark  string                `json:"auditRemark"`
	WorkStatus   int                   `json:"workStatus"`
	DisplayName  string                `json:"displayName"`
	Intro        string                `json:"intro"`
	Tags         []string              `json:"tags"`
	ServiceItems []ProviderServiceItem `json:"serviceItems"`
	CreatedAt    time.Time             `json:"createdAt"`
	UpdatedAt    time.Time             `json:"updatedAt"`
}

func cloneProvider(provider *Provider) *Provider {
	if provider == nil {
		return nil
	}
	copyProvider := *provider
	copyProvider.Tags = append([]string(nil), provider.Tags...)
	copyProvider.ServiceItems = append([]ProviderServiceItem(nil), provider.ServiceItems...)
	return &copyProvider
}

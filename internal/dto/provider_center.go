package dto

type ApplyProviderRequest struct {
	RealName         string   `json:"realName" binding:"required"`
	IDCardNo         string   `json:"idCardNo" binding:"required"`
	IDCardFront      string   `json:"idCardFront" binding:"required"`
	IDCardBack       string   `json:"idCardBack" binding:"required"`
	FaceVerifyStatus uint8    `json:"faceVerifyStatus"`
	AgreementSigned  bool     `json:"agreementSigned"`
	CityID           uint64   `json:"cityId" binding:"required"`
	Intro            string   `json:"intro"`
	Photos           []string `json:"photos"`
	ServiceDesc      string   `json:"serviceDesc"`
}

type ApplyProviderResponse struct {
	ApplicationID uint64 `json:"applicationId"`
	AuditStatus   uint8  `json:"auditStatus"`
}

type AuditStatusResponse struct {
	HasApplication bool   `json:"hasApplication"`
	AuditStatus    int    `json:"auditStatus"`
	RejectReason   string `json:"rejectReason"`
	SubmittedAt    string `json:"submittedAt,omitempty"`
	AuditedAt      string `json:"auditedAt,omitempty"`
	ProviderExists bool   `json:"providerExists"`
}

type UpdateProviderProfileRequest struct {
	Zodiac        *string   `json:"zodiac"`
	Constellation *string   `json:"constellation"`
	Level         *string   `json:"level"`
	ServiceStatus *uint8    `json:"serviceStatus"`
	OnlineStatus  *uint8    `json:"onlineStatus"`
	CityID        *uint64   `json:"cityId"`
	Intro         *string   `json:"intro"`
	Tags          *[]string `json:"tags"`
}

type ProviderProfileResponse struct {
	ID            uint64   `json:"id"`
	UserID        uint64   `json:"userId"`
	ProviderNo    string   `json:"providerNo"`
	Zodiac        string   `json:"zodiac"`
	Constellation string   `json:"constellation"`
	Level         string   `json:"level"`
	ServiceStatus uint8    `json:"serviceStatus"`
	OnlineStatus  uint8    `json:"onlineStatus"`
	CityID        uint64   `json:"cityId"`
	Intro         string   `json:"intro"`
	Tags          []string `json:"tags"`
}

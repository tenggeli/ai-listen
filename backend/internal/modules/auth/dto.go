package auth

type SendSMSRequest struct {
	Mobile string `json:"mobile"`
	Scene  string `json:"scene"`
}

type SMSLoginRequest struct {
	Mobile string `json:"mobile"`
	Code   string `json:"code"`
}

type WechatLoginRequest struct {
	Code string `json:"code"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken"`
}

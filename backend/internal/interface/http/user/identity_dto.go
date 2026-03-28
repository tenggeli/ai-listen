package user

type LoginBySMSRequestDTO struct {
	Phone             string `json:"phone"`
	VerifyCode        string `json:"verify_code"`
	AgreementAccepted bool   `json:"agreement_accepted"`
}

type LoginByWechatRequestDTO struct {
	AuthCode          string `json:"auth_code"`
	AgreementAccepted bool   `json:"agreement_accepted"`
}

type UserIdentityResponseDTO struct {
	UserID           string `json:"user_id"`
	LoginChannel     string `json:"login_channel"`
	AccessToken      string `json:"access_token"`
	RefreshToken     string `json:"refresh_token"`
	ExpiresInSeconds int64  `json:"expires_in_seconds"`
	DisplayName      string `json:"display_name"`
	AvatarURL        string `json:"avatar_url"`
	IsNewUser        bool   `json:"is_new_user"`
	ProfileCompleted bool   `json:"profile_completed"`
}

type SaveProfileRequestDTO struct {
	Nickname              string `json:"nickname"`
	AvatarURL             string `json:"avatar_url"`
	Gender                string `json:"gender"`
	AgeRange              string `json:"age_range"`
	City                  string `json:"city"`
	Bio                   string `json:"bio"`
	GenderChangeConfirmed bool   `json:"gender_change_confirmed"`
}

type UserMeResponseDTO struct {
	UserID               string   `json:"user_id"`
	Phone                string   `json:"phone"`
	Nickname             string   `json:"nickname"`
	AvatarURL            string   `json:"avatar_url"`
	RegisterSource       string   `json:"register_source"`
	Status               string   `json:"status"`
	ProfileCompleted     bool     `json:"profile_completed"`
	Gender               string   `json:"gender"`
	AgeRange             string   `json:"age_range"`
	City                 string   `json:"city"`
	Bio                  string   `json:"bio"`
	PersonalityCompleted bool     `json:"personality_completed"`
	MBTI                 string   `json:"mbti"`
	InterestTags         []string `json:"interest_tags"`
	PersonalitySkipped   bool     `json:"personality_skipped"`
}

type SavePersonalityRequestDTO struct {
	MBTI         string   `json:"mbti"`
	InterestTags []string `json:"interest_tags"`
}

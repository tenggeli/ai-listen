package dto

type UserProfileResponse struct {
	ID         uint64 `json:"id"`
	Mobile     string `json:"mobile"`
	Nickname   string `json:"nickname"`
	Avatar     string `json:"avatar"`
	Gender     uint8  `json:"gender"`
	CityID     uint64 `json:"cityId"`
	Bio        string `json:"bio"`
	VipLevel   uint8  `json:"vipLevel"`
	UserStatus uint8  `json:"userStatus"`
}

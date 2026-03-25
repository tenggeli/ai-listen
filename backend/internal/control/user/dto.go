package user

type UpdateMeRequest struct {
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Gender   int    `json:"gender"`
	Birthday string `json:"birthday"`
	CityCode string `json:"cityCode"`
}

type UpdateProfileRequest struct {
	Zodiac        string   `json:"zodiac"`
	ChineseZodiac string   `json:"chineseZodiac"`
	Bio           string   `json:"bio"`
	Interests     []string `json:"interests"`
}

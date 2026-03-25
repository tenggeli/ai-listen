package model

import "time"

type User struct {
	ID        uint64    `json:"id"`
	Mobile    string    `json:"mobile"`
	Nickname  string    `json:"nickname"`
	Avatar    string    `json:"avatar"`
	Gender    int       `json:"gender"`
	Birthday  string    `json:"birthday"`
	CityCode  string    `json:"cityCode"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func cloneUser(user *User) *User {
	if user == nil {
		return nil
	}
	copyUser := *user
	return &copyUser
}

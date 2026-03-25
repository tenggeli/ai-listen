package model

import (
	"fmt"
	"time"
)

func (s *MemoryStore) IssueSMSCode(mobile string) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.smsCodes[mobile] = "123456"
	return "123456"
}

func (s *MemoryStore) LoginBySMS(mobile, code string) (*User, string, string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.smsCodes[mobile] != code {
		return nil, "", "", ErrInvalidSMSCode
	}

	user, ok := s.usersByMobile[mobile]
	if !ok {
		now := time.Now()
		user = &User{
			ID:        s.nextUserID,
			Mobile:    mobile,
			Nickname:  fmt.Sprintf("listen用户%04d", s.nextUserID),
			CityCode:  "310100",
			CreatedAt: now,
			UpdatedAt: now,
		}
		s.nextUserID++
		s.users[user.ID] = user
		s.usersByMobile[mobile] = user
	}

	token := fmt.Sprintf("token-%d", user.ID)
	refreshToken := fmt.Sprintf("refresh-%d", user.ID)
	s.tokens[token] = user.ID
	s.refreshTokens[refreshToken] = user.ID
	return cloneUser(user), token, refreshToken, nil
}

func (s *MemoryStore) RefreshToken(refreshToken string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	userID, ok := s.refreshTokens[refreshToken]
	if !ok {
		return "", ErrUnauthorized
	}
	token := fmt.Sprintf("token-%d-%d", userID, time.Now().UnixNano())
	s.tokens[token] = userID
	return token, nil
}

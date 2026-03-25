package model

import (
	"strings"
	"time"
)

func (s *MemoryStore) UserByToken(raw string) (*User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	token := strings.TrimSpace(strings.TrimPrefix(raw, "Bearer "))
	userID, ok := s.tokens[token]
	if !ok {
		return nil, ErrUnauthorized
	}
	user, ok := s.users[userID]
	if !ok {
		return nil, ErrUserNotFound
	}
	return cloneUser(user), nil
}

func (s *MemoryStore) GetUser(userID uint64) (*User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, ok := s.users[userID]
	if !ok {
		return nil, ErrUserNotFound
	}
	return cloneUser(user), nil
}

func (s *MemoryStore) UpdateUser(userID uint64, updater func(*User)) (*User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	user, ok := s.users[userID]
	if !ok {
		return nil, ErrUserNotFound
	}
	updater(user)
	user.UpdatedAt = time.Now()
	return cloneUser(user), nil
}

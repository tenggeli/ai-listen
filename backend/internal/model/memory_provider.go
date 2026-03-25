package model

import (
	"fmt"
	"time"
)

func (s *MemoryStore) ApplyProvider(userID uint64, realName, idCardNo string) (*Provider, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if provider, ok := s.providerByUID[userID]; ok {
		provider.RealName = realName
		provider.IDCardNo = idCardNo
		provider.AuditStatus = 1
		provider.AuditRemark = ""
		provider.UpdatedAt = time.Now()
		return cloneProvider(provider), nil
	}

	now := time.Now()
	provider := &Provider{
		ID:          s.nextProviderID,
		UserID:      userID,
		ProviderNo:  fmt.Sprintf("P%08d", s.nextProviderID),
		RealName:    realName,
		IDCardNo:    idCardNo,
		AuditStatus: 1,
		WorkStatus:  1,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	s.nextProviderID++
	s.providers[provider.ID] = provider
	s.providerByUID[userID] = provider
	return cloneProvider(provider), nil
}

func (s *MemoryStore) ProviderByUserID(userID uint64) (*Provider, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	provider, ok := s.providerByUID[userID]
	if !ok {
		return nil, ErrProviderNotFound
	}
	return cloneProvider(provider), nil
}

func (s *MemoryStore) ProviderByID(providerID uint64) (*Provider, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	provider, ok := s.providers[providerID]
	if !ok {
		return nil, ErrProviderNotFound
	}
	return cloneProvider(provider), nil
}

func (s *MemoryStore) Providers() []*Provider {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]*Provider, 0, len(s.providers))
	for _, provider := range s.providers {
		result = append(result, cloneProvider(provider))
	}
	return result
}

func (s *MemoryStore) UpdateProviderProfile(userID uint64, displayName, intro string, tags []string) (*Provider, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	provider, ok := s.providerByUID[userID]
	if !ok {
		return nil, ErrProviderNotFound
	}
	provider.DisplayName = displayName
	provider.Intro = intro
	provider.Tags = append([]string(nil), tags...)
	provider.UpdatedAt = time.Now()
	return cloneProvider(provider), nil
}

func (s *MemoryStore) UpdateProviderServiceItems(userID uint64, items []ProviderServiceItem) (*Provider, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	provider, ok := s.providerByUID[userID]
	if !ok {
		return nil, ErrProviderNotFound
	}
	provider.ServiceItems = append([]ProviderServiceItem(nil), items...)
	provider.UpdatedAt = time.Now()
	return cloneProvider(provider), nil
}

func (s *MemoryStore) UpdateProviderWorkStatus(userID uint64, workStatus int) (*Provider, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	provider, ok := s.providerByUID[userID]
	if !ok {
		return nil, ErrProviderNotFound
	}
	provider.WorkStatus = workStatus
	provider.UpdatedAt = time.Now()
	return cloneProvider(provider), nil
}

func (s *MemoryStore) ApproveProvider(providerID uint64, remark string) (*Provider, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	provider, ok := s.providers[providerID]
	if !ok {
		return nil, ErrProviderNotFound
	}
	provider.AuditStatus = 2
	provider.AuditRemark = remark
	provider.UpdatedAt = time.Now()
	return cloneProvider(provider), nil
}

func (s *MemoryStore) RejectProvider(providerID uint64, remark string) (*Provider, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	provider, ok := s.providers[providerID]
	if !ok {
		return nil, ErrProviderNotFound
	}
	provider.AuditStatus = 3
	provider.AuditRemark = remark
	provider.UpdatedAt = time.Now()
	return cloneProvider(provider), nil
}

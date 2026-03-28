package memory

import (
	"context"
	"sync"

	domain "listen/backend/internal/domain/identity"
)

type IdentityRepository struct {
	mu           sync.RWMutex
	byPhone      map[string]string
	byWechat     map[string]string
	accountsByID map[string]domain.UserAccount
}

func NewIdentityRepository() *IdentityRepository {
	return &IdentityRepository{
		byPhone:      make(map[string]string),
		byWechat:     make(map[string]string),
		accountsByID: make(map[string]domain.UserAccount),
	}
}

func (r *IdentityRepository) GetByPhone(_ context.Context, phone string) (domain.UserAccount, bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	userID, ok := r.byPhone[phone]
	if !ok {
		return domain.UserAccount{}, false, nil
	}
	account, ok := r.accountsByID[userID]
	if !ok {
		return domain.UserAccount{}, false, nil
	}
	return account, true, nil
}

func (r *IdentityRepository) GetByID(_ context.Context, userID string) (domain.UserAccount, bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	account, ok := r.accountsByID[userID]
	if !ok {
		return domain.UserAccount{}, false, nil
	}
	return account, true, nil
}

func (r *IdentityRepository) GetByWechatOpenID(_ context.Context, openID string) (domain.UserAccount, bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	userID, ok := r.byWechat[openID]
	if !ok {
		return domain.UserAccount{}, false, nil
	}
	account, ok := r.accountsByID[userID]
	if !ok {
		return domain.UserAccount{}, false, nil
	}
	return account, true, nil
}

func (r *IdentityRepository) Save(_ context.Context, account domain.UserAccount) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.accountsByID[account.UserID] = account
	if account.Phone != "" {
		r.byPhone[account.Phone] = account.UserID
	}
	if account.WechatOpenID != "" {
		r.byWechat[account.WechatOpenID] = account.UserID
	}
	return nil
}

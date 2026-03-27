package memory

import (
	"context"
	"sync"

	domain "listen/backend/internal/domain/ai"
)

type SessionRepository struct {
	mu   sync.RWMutex
	data map[string]domain.Session
}

func NewSessionRepository() *SessionRepository {
	return &SessionRepository{data: make(map[string]domain.Session)}
}

func (r *SessionRepository) Create(_ context.Context, session domain.Session) error {
	r.mu.Lock()
	r.data[session.ID] = session
	r.mu.Unlock()
	return nil
}

func (r *SessionRepository) GetByID(_ context.Context, id string) (domain.Session, error) {
	r.mu.RLock()
	session, ok := r.data[id]
	r.mu.RUnlock()
	if !ok {
		return domain.Session{}, domain.ErrSessionNotFound
	}
	return session, nil
}

func (r *SessionRepository) Save(_ context.Context, session domain.Session) error {
	r.mu.Lock()
	r.data[session.ID] = session
	r.mu.Unlock()
	return nil
}

package ai

import (
	"errors"
	"time"
)

const DailyMatchLimit = 5

var (
	ErrDailyLimitReached = errors.New("daily match limit reached")
	ErrInvalidInput      = errors.New("invalid input")
	ErrSessionNotFound   = errors.New("session not found")
)

type MatchCandidate struct {
	ProviderID  string
	DisplayName string
	ReasonText  string
	Score       float64
}

type MatchResult struct {
	Candidates []MatchCandidate
}

func (r MatchResult) IsEmpty() bool {
	return len(r.Candidates) == 0
}

type DailyQuota struct {
	UserID string
	Date   string
	Used   int
}

func NewDailyQuota(userID string, date string) DailyQuota {
	return DailyQuota{UserID: userID, Date: date, Used: 0}
}

func (q DailyQuota) Remaining() int {
	left := DailyMatchLimit - q.Used
	if left < 0 {
		return 0
	}
	return left
}

func (q *DailyQuota) Consume() error {
	if q.Remaining() == 0 {
		return ErrDailyLimitReached
	}
	q.Used++
	return nil
}

type Session struct {
	ID            string
	UserID        string
	SceneType     string
	Summary       string
	Status        string
	LastMessageAt time.Time
	Messages      []Message
}

type Message struct {
	SenderType string
	Content    string
	CreatedAt  time.Time
}

func NewSession(id string, userID string, sceneType string) (Session, error) {
	if id == "" || userID == "" {
		return Session{}, ErrInvalidInput
	}
	return Session{
		ID:        id,
		UserID:    userID,
		SceneType: sceneType,
		Status:    "active",
		Messages:  make([]Message, 0),
	}, nil
}

func (s *Session) AppendMessage(senderType string, content string, now time.Time) error {
	if senderType == "" || content == "" {
		return ErrInvalidInput
	}
	message := Message{SenderType: senderType, Content: content, CreatedAt: now}
	s.Messages = append(s.Messages, message)
	s.LastMessageAt = now
	return nil
}

package ai

import (
	"errors"
	"strings"
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

type HomeQuickAction struct {
	Key    string
	Label  string
	Icon   string
	Route  string
	Prompt string
}

type HomeOverview struct {
	GreetingPeriod string
	GreetingText   string
	GreetingSub    string
	MoodEmoji      string
	MoodText       string
	WeatherText    string
	CompanionDays  int
	OnlineCount    int
	WaitingCount   int
	Remaining      int
	QuickActions   []HomeQuickAction
}

func (o HomeOverview) IsEmpty() bool {
	return o.GreetingText == "" && len(o.QuickActions) == 0
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
	trimmed := strings.TrimSpace(content)
	if senderType == "" || trimmed == "" {
		return ErrInvalidInput
	}
	message := Message{SenderType: senderType, Content: trimmed, CreatedAt: now}
	s.Messages = append(s.Messages, message)
	s.LastMessageAt = now
	return nil
}

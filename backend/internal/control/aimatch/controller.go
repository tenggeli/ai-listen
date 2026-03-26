package aimatch

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"ai-listen/backend/internal/model"
	"ai-listen/backend/pkg/ecode"
	"ai-listen/backend/pkg/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Controller struct {
	logger *zap.Logger
}

type MatchRecommendation struct {
	ProviderID   uint64   `json:"providerId"`
	DisplayName  string   `json:"displayName"`
	MatchScore   int      `json:"matchScore"`
	MatchedTags  []string `json:"matchedTags"`
	Reason       string   `json:"reason"`
	WorkStatus   int      `json:"workStatus"`
	ServiceCount int      `json:"serviceCount"`
}

type MatchSession struct {
	SessionID       string                `json:"sessionId"`
	InputType       string                `json:"inputType"`
	Content         string                `json:"content"`
	CityCode        string                `json:"cityCode"`
	Recommendations []MatchRecommendation `json:"recommendations"`
	CreatedAt       time.Time             `json:"createdAt"`
	ExpiresAt       time.Time             `json:"expiresAt"`
}

type matchState struct {
	mu       sync.RWMutex
	seq      uint64
	sessions map[string]MatchSession
}

var aiMatchStore = &matchState{
	sessions: map[string]MatchSession{},
}

func (s *matchState) createSession(req MatchRequest, recommendations []MatchRecommendation) MatchSession {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.seq++
	now := time.Now()
	sessionID := fmt.Sprintf("M%08d", s.seq)
	session := MatchSession{
		SessionID:       sessionID,
		InputType:       strings.TrimSpace(req.InputType),
		Content:         strings.TrimSpace(req.Content),
		CityCode:        strings.TrimSpace(req.CityCode),
		Recommendations: recommendations,
		CreatedAt:       now,
		ExpiresAt:       now.Add(30 * time.Minute),
	}
	s.sessions[sessionID] = session
	return session
}

func (s *matchState) getSession(sessionID string) (MatchSession, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	session, ok := s.sessions[sessionID]
	if !ok {
		return MatchSession{}, false
	}
	if time.Now().After(session.ExpiresAt) {
		return MatchSession{}, false
	}
	return session, true
}

func (s *matchState) gcExpired() {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	for key, session := range s.sessions {
		if now.After(session.ExpiresAt) {
			delete(s.sessions, key)
		}
	}
}

func buildRecommendations(req MatchRequest) []MatchRecommendation {
	providers := model.Default().Providers()
	if len(providers) == 0 {
		return []MatchRecommendation{}
	}
	keywords := extractKeywords(req.Content)
	result := make([]MatchRecommendation, 0, len(providers))
	for _, provider := range providers {
		if provider == nil {
			continue
		}
		if provider.AuditStatus != 2 {
			continue
		}
		score, matchedTags := matchScore(provider, keywords)
		result = append(result, MatchRecommendation{
			ProviderID:   provider.ID,
			DisplayName:  provider.DisplayName,
			MatchScore:   score,
			MatchedTags:  matchedTags,
			Reason:       matchReason(matchedTags),
			WorkStatus:   provider.WorkStatus,
			ServiceCount: len(provider.ServiceItems),
		})
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].MatchScore == result[j].MatchScore {
			return result[i].ProviderID < result[j].ProviderID
		}
		return result[i].MatchScore > result[j].MatchScore
	})
	if len(result) > 5 {
		result = result[:5]
	}
	return result
}

func extractKeywords(content string) []string {
	words := strings.Fields(strings.ToLower(content))
	ext := make([]string, 0, len(words)+6)
	ext = append(ext, words...)
	lower := strings.ToLower(content)
	keywordMap := map[string]string{
		"失眠": "助眠",
		"焦虑": "减压",
		"压力": "减压",
		"聊天": "陪伴",
		"陪伴": "陪伴",
		"专注": "专注",
		"放松": "放松",
	}
	for source, normalized := range keywordMap {
		if strings.Contains(lower, source) {
			ext = append(ext, normalized)
		}
	}
	uniq := map[string]struct{}{}
	result := make([]string, 0, len(ext))
	for _, word := range ext {
		word = strings.TrimSpace(word)
		if word == "" {
			continue
		}
		if _, ok := uniq[word]; ok {
			continue
		}
		uniq[word] = struct{}{}
		result = append(result, word)
	}
	return result
}

func matchScore(provider *model.Provider, keywords []string) (int, []string) {
	score := 60
	matched := []string{}
	if provider.WorkStatus == 1 {
		score += 12
	}
	if provider.DisplayName != "" {
		score += 4
	}
	if len(provider.ServiceItems) > 0 {
		score += 6
	}
	for _, tag := range provider.Tags {
		lowerTag := strings.ToLower(tag)
		for _, kw := range keywords {
			if strings.Contains(lowerTag, kw) || strings.Contains(kw, lowerTag) {
				score += 8
				matched = append(matched, tag)
				break
			}
		}
	}
	if len(matched) == 0 && strings.TrimSpace(provider.Intro) != "" {
		for _, kw := range keywords {
			if strings.Contains(strings.ToLower(provider.Intro), kw) {
				score += 5
				matched = append(matched, kw)
			}
		}
	}
	if score > 98 {
		score = 98
	}
	return score, dedupeStrings(matched)
}

func dedupeStrings(values []string) []string {
	if len(values) == 0 {
		return values
	}
	result := make([]string, 0, len(values))
	seen := map[string]struct{}{}
	for _, value := range values {
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	return result
}

func matchReason(tags []string) string {
	if len(tags) > 0 {
		return "标签匹配度高"
	}
	return "基础资料完整，适合推荐"
}

func (h *Controller) Match(c *gin.Context) {
	var req MatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, ecode.BadRequest, gin.H{"reason": "invalid request body"})
		return
	}
	if strings.TrimSpace(req.Content) == "" {
		response.Fail(c, http.StatusBadRequest, ecode.BadRequest, gin.H{"reason": "content is required"})
		return
	}
	aiMatchStore.gcExpired()
	recommendations := buildRecommendations(req)
	session := aiMatchStore.createSession(req, recommendations)
	response.Success(c, gin.H{
		"module":          "ai",
		"action":          "match",
		"sessionId":       session.SessionID,
		"request":         req,
		"recommendations": recommendations,
		"expiresAt":       session.ExpiresAt,
	})
}

func (h *Controller) GetMatch(c *gin.Context) {
	sessionID := c.Param("sessionId")
	session, ok := aiMatchStore.getSession(sessionID)
	if !ok {
		response.Fail(c, http.StatusNotFound, ecode.NotFound, gin.H{"reason": "match session not found"})
		return
	}
	response.Success(c, gin.H{
		"module":          "ai",
		"action":          "get_match",
		"sessionId":       session.SessionID,
		"request":         gin.H{"inputType": session.InputType, "content": session.Content, "cityCode": session.CityCode},
		"recommendations": session.Recommendations,
		"createdAt":       session.CreatedAt,
		"expiresAt":       session.ExpiresAt,
	})
}

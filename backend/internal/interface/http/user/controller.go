package user

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	app "listen/backend/internal/application/ai"
	domain "listen/backend/internal/domain/ai"
)

type AIController struct {
	getHome      app.GetAiHomeUseCase
	getRemaining app.GetRemainingMatchUseCase
	submitMatch  app.SubmitMatchUseCase
	createSess   app.CreateAiSessionUseCase
	getSess      app.GetAiSessionUseCase
	appendMsg    app.AppendAiMessageUseCase
}

func NewAIController(
	getHome app.GetAiHomeUseCase,
	getRemaining app.GetRemainingMatchUseCase,
	submitMatch app.SubmitMatchUseCase,
	createSess app.CreateAiSessionUseCase,
	getSess app.GetAiSessionUseCase,
	appendMsg app.AppendAiMessageUseCase,
) AIController {
	return AIController{
		getHome:      getHome,
		getRemaining: getRemaining,
		submitMatch:  submitMatch,
		createSess:   createSess,
		getSess:      getSess,
		appendMsg:    appendMsg,
	}
}

func (c AIController) HandleGetHome(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	output, err := c.getHome.Execute(r.Context(), app.GetHomeOverviewInput{UserID: userID})
	if err != nil {
		writeDomainError(w, err)
		return
	}

	actions := make([]map[string]any, 0, len(output.Overview.QuickActions))
	for _, item := range output.Overview.QuickActions {
		actions = append(actions, map[string]any{
			"key":    item.Key,
			"label":  item.Label,
			"icon":   item.Icon,
			"route":  item.Route,
			"prompt": item.Prompt,
		})
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"greeting_period":   output.Overview.GreetingPeriod,
		"greeting_text":     output.Overview.GreetingText,
		"greeting_sub_text": output.Overview.GreetingSub,
		"mood_emoji":        output.Overview.MoodEmoji,
		"mood_text":         output.Overview.MoodText,
		"weather_text":      output.Overview.WeatherText,
		"companion_days":    output.Overview.CompanionDays,
		"online_count":      output.Overview.OnlineCount,
		"waiting_count":     output.Overview.WaitingCount,
		"remaining":         output.Overview.Remaining,
		"quick_actions":     actions,
	})
}

func (c AIController) HandleGetRemaining(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	output, err := c.getRemaining.Execute(r.Context(), app.GetRemainingMatchInput{UserID: userID})
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"remaining": output.Remaining})
}

func (c AIController) HandleMatch(w http.ResponseWriter, r *http.Request) {
	var body MatchRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	output, err := c.submitMatch.Execute(r.Context(), app.SubmitMatchInput{UserID: body.UserID, InputText: body.InputText})
	if err != nil {
		writeDomainError(w, err)
		return
	}

	items := make([]map[string]any, 0, len(output.Candidates))
	for _, c := range output.Candidates {
		items = append(items, map[string]any{
			"provider_id":  c.ProviderID,
			"display_name": c.DisplayName,
			"reason_text":  c.ReasonText,
			"score":        c.Score,
		})
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"remaining":  output.Remaining,
		"candidates": items,
	})
}

func (c AIController) HandleCreateSession(w http.ResponseWriter, r *http.Request) {
	var body CreateSessionRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	output, err := c.createSess.Execute(r.Context(), app.CreateSessionInput{UserID: body.UserID, SceneType: body.SceneType})
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"session_id": output.SessionID})
}

func (c AIController) HandleSessionDetail(w http.ResponseWriter, r *http.Request) {
	sessionID := sessionIDFromPath(r.URL.Path)
	if sessionID == "" {
		writeJSONError(w, http.StatusBadRequest, "invalid session id")
		return
	}

	output, err := c.getSess.Execute(r.Context(), app.GetSessionInput{SessionID: sessionID})
	if err != nil {
		writeDomainError(w, err)
		return
	}

	messages := make([]map[string]any, 0, len(output.Session.Messages))
	for _, message := range output.Session.Messages {
		item := map[string]any{
			"sender_type": message.SenderType,
			"content":     message.Content,
			"created_at":  message.CreatedAt.Format(time.RFC3339),
		}
		if message.SafetyLevel != "" {
			item["safety_level"] = message.SafetyLevel
		}
		if message.ActionCard != nil {
			item["action_card"] = map[string]any{
				"action":      message.ActionCard.Action,
				"title":       message.ActionCard.Title,
				"description": message.ActionCard.Description,
				"route":       message.ActionCard.Route,
				"button_text": message.ActionCard.ButtonText,
			}
		}
		messages = append(messages, item)
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"id":              output.Session.ID,
		"user_id":         output.Session.UserID,
		"scene_type":      output.Session.SceneType,
		"status":          output.Session.Status,
		"last_message_at": output.Session.LastMessageAt.Format(time.RFC3339),
		"messages":        messages,
	})
}

func (c AIController) HandleAppendMessage(w http.ResponseWriter, r *http.Request) {
	sessionID := sessionIDFromPath(strings.TrimSuffix(r.URL.Path, "/messages"))
	if sessionID == "" {
		writeJSONError(w, http.StatusBadRequest, "invalid session id")
		return
	}

	var body AppendMessageRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	_, err := c.appendMsg.Execute(r.Context(), app.AppendMessageInput{
		SessionID:   sessionID,
		SenderType:  body.SenderType,
		ContentText: body.Content,
	})
	if err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"session_id": sessionID})
}

func writeDomainError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrInvalidInput):
		writeJSONError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, domain.ErrDailyLimitReached):
		writeJSONError(w, http.StatusTooManyRequests, err.Error())
	case errors.Is(err, domain.ErrSessionNotFound):
		writeJSONError(w, http.StatusNotFound, err.Error())
	default:
		writeJSONError(w, http.StatusInternalServerError, "internal error")
	}
}

func sessionIDFromPath(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) < 4 {
		return ""
	}
	return parts[3]
}

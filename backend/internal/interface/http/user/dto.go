package user

type MatchRequestDTO struct {
	UserID    string `json:"user_id"`
	InputText string `json:"input_text"`
}

type CreateSessionRequestDTO struct {
	UserID    string `json:"user_id"`
	SceneType string `json:"scene_type"`
}

type AppendMessageRequestDTO struct {
	SenderType string `json:"sender_type"`
	Content    string `json:"content"`
}

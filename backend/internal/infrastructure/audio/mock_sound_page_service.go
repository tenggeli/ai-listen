package audio

import (
	"context"

	domain "listen/backend/internal/domain/audio"
)

type MockSoundPageService struct{}

func NewMockSoundPageService() MockSoundPageService {
	return MockSoundPageService{}
}

func (MockSoundPageService) GetHomePage(_ context.Context, _ string) (domain.HomePage, error) {
	categories := []domain.Category{
		{Key: "all", Label: "全部"},
		{Key: "nature", Label: "自然白噪音"},
		{Key: "sleep", Label: "睡眠引导"},
		{Key: "meditation", Label: "正念冥想"},
		{Key: "story", Label: "治愈故事"},
		{Key: "breath", Label: "呼吸练习"},
	}

	tracks := []domain.Track{
		{ID: "track-rain", Title: "深夜雨声 · 安眠版", Category: "自然白噪音", PlayCountText: "1,248 次播放", DurationText: "35:00", Emoji: "🌧", Author: "listen 治愈声音库"},
		{ID: "track-wave", Title: "海浪呼吸引导 · 4-7-8 法", Category: "呼吸练习", PlayCountText: "856 次播放", DurationText: "12:00", Emoji: "🌊", Author: "listen 治愈声音库"},
		{ID: "track-forest", Title: "森林清晨 · 入睡前冥想", Category: "正念冥想", PlayCountText: "2,034 次播放", DurationText: "20:30", Emoji: "🍃", Author: "listen 治愈声音库"},
		{ID: "track-radio", Title: "今晚陪你说话 · 深夜电台 Vol.12", Category: "治愈故事", PlayCountText: "3,671 次播放", DurationText: "28:15", Emoji: "🌙", Author: "listen 深夜电台"},
		{ID: "track-fire", Title: "壁炉暖意 · 冬夜陪伴", Category: "自然白噪音", PlayCountText: "987 次播放", DurationText: "60:00", Emoji: "🔥", Author: "listen 治愈声音库"},
		{ID: "track-star", Title: "星空冥想 · 放下今天的重量", Category: "正念冥想", PlayCountText: "4,122 次播放", DurationText: "18:00", Emoji: "⭐", Author: "listen 冥想室"},
	}

	return domain.HomePage{
		Title:               "声音",
		Subtitle:            "用声音抚慰此刻的你",
		Categories:          categories,
		RecommendedTracks:   tracks,
		CurrentTrackID:      "track-rain",
		CurrentProgressText: "12:34",
		TotalDurationText:   "35:00",
		IsPlaying:           true,
	}, nil
}

var _ domain.HomePageService = MockSoundPageService{}

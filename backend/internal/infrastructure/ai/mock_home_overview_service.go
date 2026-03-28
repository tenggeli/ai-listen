package ai

import (
	"context"
	"time"

	domain "listen/backend/internal/domain/ai"
)

type MockHomeOverviewService struct{}

func NewMockHomeOverviewService() MockHomeOverviewService {
	return MockHomeOverviewService{}
}

func (MockHomeOverviewService) BuildOverview(_ context.Context, input domain.HomeOverviewServiceInput) (domain.HomeOverview, error) {
	return domain.HomeOverview{
		GreetingPeriod: mockGreetingPeriod(input.Now),
		GreetingText:   "今晚，遇见懂你的人",
		GreetingSub:    "有 1,247 位搭子正在等待陪伴",
		MoodEmoji:      "🌙",
		MoodText:       "平静 · 有点想聊聊",
		WeatherText:    "上海 21°C 微风",
		CompanionDays:  28,
		OnlineCount:    1247,
		WaitingCount:   312,
		Remaining:      input.Remaining,
		QuickActions: []domain.HomeQuickAction{
			{Key: "quick-join", Label: "快速加入", Icon: "join", Route: "/chat", Prompt: "想快速加入一个轻松的聊天陪伴场景"},
			{Key: "square", Label: "热门广场", Icon: "square", Route: "/home", Prompt: "想看看大家最近都在聊什么"},
			{Key: "voice", Label: "治愈声音", Icon: "voice", Route: "/sound", Prompt: "我想听一点能让我放松下来的声音"},
			{Key: "topic", Label: "今日话题", Icon: "topic", Route: "/home", Prompt: "给我一个今晚适合开启聊天的话题"},
		},
	}, nil
}

func mockGreetingPeriod(now time.Time) string {
	period := "晚上好"
	switch hour := now.Hour(); {
	case hour < 11:
		period = "早上好"
	case hour < 14:
		period = "中午好"
	case hour < 18:
		period = "下午好"
	}

	weekdays := []string{"周日", "周一", "周二", "周三", "周四", "周五", "周六"}
	return period + " · " + weekdays[int(now.Weekday())]
}

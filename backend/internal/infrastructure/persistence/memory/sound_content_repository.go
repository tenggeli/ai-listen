package memory

import (
	"context"
	"sort"
	"strings"
	"sync"

	domain "listen/backend/internal/domain/audio"
)

type SoundContentRepository struct {
	mu         sync.RWMutex
	categories []domain.Category
	tracks     []memorySoundTrack
}

type memorySoundTrack struct {
	categoryKey string
	track       domain.Track
	sortOrder   int
	status      string
}

func NewSoundContentRepository() *SoundContentRepository {
	return &SoundContentRepository{
		categories: []domain.Category{
			{Key: "all", Label: "全部"},
			{Key: "nature", Label: "自然白噪音"},
			{Key: "sleep", Label: "睡眠引导"},
			{Key: "meditation", Label: "正念冥想"},
			{Key: "story", Label: "治愈故事"},
			{Key: "breath", Label: "呼吸练习"},
		},
		tracks: []memorySoundTrack{
			{
				categoryKey: "nature",
				sortOrder:   1,
				status:      domain.SoundStatusActive,
				track: domain.Track{
					ID:            "track-rain",
					Title:         "深夜雨声 · 安眠版",
					Category:      "自然白噪音",
					PlayCountText: "1,248 次播放",
					DurationText:  "35:00",
					Emoji:         "🌧",
					Author:        "listen 治愈声音库",
				},
			},
			{
				categoryKey: "breath",
				sortOrder:   2,
				status:      domain.SoundStatusActive,
				track: domain.Track{
					ID:            "track-wave",
					Title:         "海浪呼吸引导 · 4-7-8 法",
					Category:      "呼吸练习",
					PlayCountText: "856 次播放",
					DurationText:  "12:00",
					Emoji:         "🌊",
					Author:        "listen 治愈声音库",
				},
			},
			{
				categoryKey: "meditation",
				sortOrder:   3,
				status:      domain.SoundStatusActive,
				track: domain.Track{
					ID:            "track-forest",
					Title:         "森林清晨 · 入睡前冥想",
					Category:      "正念冥想",
					PlayCountText: "2,034 次播放",
					DurationText:  "20:30",
					Emoji:         "🍃",
					Author:        "listen 治愈声音库",
				},
			},
			{
				categoryKey: "story",
				sortOrder:   4,
				status:      domain.SoundStatusActive,
				track: domain.Track{
					ID:            "track-radio",
					Title:         "今晚陪你说话 · 深夜电台 Vol.12",
					Category:      "治愈故事",
					PlayCountText: "3,671 次播放",
					DurationText:  "28:15",
					Emoji:         "🌙",
					Author:        "listen 深夜电台",
				},
			},
			{
				categoryKey: "nature",
				sortOrder:   5,
				status:      domain.SoundStatusActive,
				track: domain.Track{
					ID:            "track-fire",
					Title:         "壁炉暖意 · 冬夜陪伴",
					Category:      "自然白噪音",
					PlayCountText: "987 次播放",
					DurationText:  "60:00",
					Emoji:         "🔥",
					Author:        "listen 治愈声音库",
				},
			},
			{
				categoryKey: "meditation",
				sortOrder:   6,
				status:      domain.SoundStatusActive,
				track: domain.Track{
					ID:            "track-star",
					Title:         "星空冥想 · 放下今天的重量",
					Category:      "正念冥想",
					PlayCountText: "4,122 次播放",
					DurationText:  "18:00",
					Emoji:         "⭐",
					Author:        "listen 冥想室",
				},
			},
		},
	}
}

func (r *SoundContentRepository) ListCategories(_ context.Context) ([]domain.Category, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	items := make([]domain.Category, len(r.categories))
	copy(items, r.categories)
	return items, nil
}

func (r *SoundContentRepository) ListTracks(_ context.Context, query domain.TrackQuery) ([]domain.Track, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	items := make([]memorySoundTrack, 0, len(r.tracks))
	for _, item := range r.tracks {
		if item.trackStatus() != domain.SoundStatusActive {
			continue
		}
		if query.CategoryKey != "" && query.CategoryKey != "all" && item.categoryKey != query.CategoryKey {
			continue
		}
		items = append(items, item)
	}

	sort.Slice(items, func(i, j int) bool {
		if items[i].sortOrder == items[j].sortOrder {
			return items[i].track.ID < items[j].track.ID
		}
		return items[i].sortOrder < items[j].sortOrder
	})

	total := len(items)
	offset := (query.PageNo - 1) * query.PageSize
	if offset >= total {
		return []domain.Track{}, total, nil
	}
	end := offset + query.PageSize
	if end > total {
		end = total
	}

	result := make([]domain.Track, 0, end-offset)
	for _, item := range items[offset:end] {
		result = append(result, item.track)
	}
	return result, total, nil
}

func (r *SoundContentRepository) ListAdminTracks(_ context.Context, query domain.AdminTrackQuery) ([]domain.AdminTrack, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	keyword := strings.ToLower(strings.TrimSpace(query.Keyword))
	items := make([]memorySoundTrack, 0, len(r.tracks))
	for _, item := range r.tracks {
		if query.CategoryKey != "" && query.CategoryKey != "all" && item.categoryKey != query.CategoryKey {
			continue
		}
		if query.Status != "" && item.trackStatus() != query.Status {
			continue
		}
		if keyword != "" {
			title := strings.ToLower(item.track.Title)
			author := strings.ToLower(item.track.Author)
			if !strings.Contains(title, keyword) && !strings.Contains(author, keyword) && !strings.Contains(strings.ToLower(item.track.ID), keyword) {
				continue
			}
		}
		items = append(items, item)
	}

	sort.Slice(items, func(i, j int) bool {
		if items[i].sortOrder == items[j].sortOrder {
			return items[i].track.ID < items[j].track.ID
		}
		return items[i].sortOrder < items[j].sortOrder
	})

	total := len(items)
	offset := (query.PageNo - 1) * query.PageSize
	if offset >= total {
		return []domain.AdminTrack{}, total, nil
	}
	end := offset + query.PageSize
	if end > total {
		end = total
	}

	result := make([]domain.AdminTrack, 0, end-offset)
	for _, item := range items[offset:end] {
		result = append(result, item.toAdminTrack())
	}
	return result, total, nil
}

func (r *SoundContentRepository) GetAdminTrackByID(_ context.Context, trackID string) (domain.AdminTrack, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, item := range r.tracks {
		if item.track.ID == trackID {
			return item.toAdminTrack(), nil
		}
	}
	return domain.AdminTrack{}, domain.ErrSoundNotFound
}

func (r *SoundContentRepository) CreateAdminTrack(_ context.Context, input domain.AdminTrack) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, item := range r.tracks {
		if item.track.ID == input.ID {
			return domain.ErrInvalidInput
		}
	}
	r.tracks = append(r.tracks, memorySoundTrack{
		categoryKey: input.CategoryKey,
		sortOrder:   input.SortOrder,
		status:      input.Status,
		track: domain.Track{
			ID:            input.ID,
			Title:         input.Title,
			Category:      r.categoryLabel(input.CategoryKey),
			PlayCountText: input.PlayCountText,
			DurationText:  input.DurationText,
			Emoji:         input.Emoji,
			Author:        input.Author,
		},
	})
	return nil
}

func (r *SoundContentRepository) UpdateAdminTrack(_ context.Context, input domain.AdminTrack) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i, item := range r.tracks {
		if item.track.ID != input.ID {
			continue
		}
		item.categoryKey = input.CategoryKey
		item.sortOrder = input.SortOrder
		item.track.Title = input.Title
		item.track.Category = r.categoryLabel(input.CategoryKey)
		item.track.PlayCountText = input.PlayCountText
		item.track.DurationText = input.DurationText
		item.track.Emoji = input.Emoji
		item.track.Author = input.Author
		r.tracks[i] = item
		return nil
	}
	return domain.ErrSoundNotFound
}

func (r *SoundContentRepository) UpdateAdminTrackStatus(_ context.Context, trackID string, status string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i, item := range r.tracks {
		if item.track.ID != trackID {
			continue
		}
		item.status = status
		r.tracks[i] = item
		return nil
	}
	return domain.ErrSoundNotFound
}

func (r *SoundContentRepository) categoryLabel(categoryKey string) string {
	for _, item := range r.categories {
		if item.Key == categoryKey {
			return item.Label
		}
	}
	return ""
}

func (t memorySoundTrack) trackStatus() string {
	if strings.TrimSpace(t.status) == "" {
		return domain.SoundStatusActive
	}
	return t.status
}

func (t memorySoundTrack) toAdminTrack() domain.AdminTrack {
	return domain.AdminTrack{
		ID:            t.track.ID,
		CategoryKey:   t.categoryKey,
		Title:         t.track.Title,
		PlayCountText: t.track.PlayCountText,
		DurationText:  t.track.DurationText,
		Emoji:         t.track.Emoji,
		Author:        t.track.Author,
		SortOrder:     t.sortOrder,
		Status:        t.trackStatus(),
	}
}

var _ domain.Repository = (*SoundContentRepository)(nil)

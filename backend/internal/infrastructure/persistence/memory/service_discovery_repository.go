package memory

import (
	"context"
	"sort"
	"strings"
	"sync"
	"time"

	domain "listen/backend/internal/domain/service_discovery"
)

type ServiceDiscoveryRepository struct {
	mu         sync.RWMutex
	categories []domain.ServiceCategory
	providers  map[string]domain.ProviderPublicProfile
	items      []domain.ServiceItem
	presence   map[string]time.Time
}

func NewServiceDiscoveryRepository() *ServiceDiscoveryRepository {
	categories := []domain.ServiceCategory{
		{ID: "cat_all", Name: "全部", Icon: "sparkles", SortOrder: 1},
		{ID: "cat_food", Name: "饭搭子", Icon: "utensils", SortOrder: 2},
		{ID: "cat_movie", Name: "电影搭子", Icon: "film", SortOrder: 3},
		{ID: "cat_chat", Name: "散步聊天", Icon: "message-circle", SortOrder: 4},
		{ID: "cat_relax", Name: "心理疏导", Icon: "heart", SortOrder: 5},
	}

	providers := map[string]domain.ProviderPublicProfile{
		"p_pub_001": {
			ID:               "p_pub_001",
			DisplayName:      "暖心倾听师 · 小林",
			AvatarURL:        "https://mock.listen.local/avatar/p_pub_001.png",
			CityCode:         "310100",
			Bio:              "擅长情绪陪伴、轻聊天与晚间低压见面。",
			RatingAvg:        4.9,
			CompletedOrders:  128,
			Online:           true,
			VerificationText: "实名认证",
			Tags:             []string{"深夜聊天", "温柔倾听", "同城可约"},
			PriceFrom:        99,
			PriceUnit:        "30分钟",
		},
		"p_pub_002": {
			ID:               "p_pub_002",
			DisplayName:      "电影散步搭子 · 念念",
			AvatarURL:        "https://mock.listen.local/avatar/p_pub_002.png",
			CityCode:         "440100",
			Bio:              "适合想出门透口气的人，节奏轻松不尴尬。",
			RatingAvg:        4.8,
			CompletedOrders:  76,
			Online:           true,
			VerificationText: "高评分",
			Tags:             []string{"电影", "散步", "ENFP"},
			PriceFrom:        158,
			PriceUnit:        "小时",
		},
		"p_pub_003": {
			ID:               "p_pub_003",
			DisplayName:      "睡前放松向导 · 阿乔",
			AvatarURL:        "https://mock.listen.local/avatar/p_pub_003.png",
			CityCode:         "110100",
			Bio:              "偏向声音陪伴、睡前放松和安抚式聊天。",
			RatingAvg:        4.9,
			CompletedOrders:  205,
			Online:           true,
			VerificationText: "夜间优先",
			Tags:             []string{"助眠", "治愈声音", "INFP"},
			PriceFrom:        99,
			PriceUnit:        "30分钟",
		},
	}

	items := []domain.ServiceItem{
		{ID: "si_001", ProviderID: "p_pub_001", CategoryID: "cat_chat", Title: "线上倾听 30 分钟", Description: "适合情绪安抚、想找人说话、睡前放松。", PriceAmount: 99, PriceUnit: "30分钟", SupportOnline: true, SortOrder: 1},
		{ID: "si_002", ProviderID: "p_pub_001", CategoryID: "cat_chat", Title: "同城散步聊天 1 小时", Description: "适合下班后想透口气、想有人一起走一段路。", PriceAmount: 158, PriceUnit: "小时", SupportOnline: false, SortOrder: 2},
		{ID: "si_003", ProviderID: "p_pub_001", CategoryID: "cat_movie", Title: "电影陪伴 1 场", Description: "偏轻社交陪伴型，适合慢慢破冰。", PriceAmount: 188, PriceUnit: "次", SupportOnline: false, SortOrder: 3},
		{ID: "si_004", ProviderID: "p_pub_002", CategoryID: "cat_movie", Title: "电影散步陪伴 1 场", Description: "看电影后可散步聊天，低压轻社交。", PriceAmount: 188, PriceUnit: "次", SupportOnline: false, SortOrder: 1},
		{ID: "si_005", ProviderID: "p_pub_002", CategoryID: "cat_chat", Title: "下班散步聊天 1 小时", Description: "适合今晚想换换心情的人。", PriceAmount: 158, PriceUnit: "小时", SupportOnline: false, SortOrder: 2},
		{ID: "si_006", ProviderID: "p_pub_003", CategoryID: "cat_relax", Title: "睡前放松语音 30 分钟", Description: "适合焦虑、失眠、需要被稳稳接住的时刻。", PriceAmount: 99, PriceUnit: "30分钟", SupportOnline: true, SortOrder: 1},
	}

	now := time.Now()
	return &ServiceDiscoveryRepository{
		categories: categories,
		providers:  providers,
		items:      items,
		presence: map[string]time.Time{
			"p_pub_001": now,
			"p_pub_002": now.Add(-2 * time.Minute),
			"p_pub_003": now.Add(-8 * time.Minute),
		},
	}
}

func (r *ServiceDiscoveryRepository) ListCategories(_ context.Context) ([]domain.ServiceCategory, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	items := make([]domain.ServiceCategory, len(r.categories))
	copy(items, r.categories)
	sort.Slice(items, func(i, j int) bool {
		return items[i].SortOrder < items[j].SortOrder
	})
	return items, nil
}

func (r *ServiceDiscoveryRepository) ListPublicProviders(_ context.Context, query domain.PublicProviderQuery) ([]domain.ProviderPublicProfile, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	matched := make([]domain.ProviderPublicProfile, 0, len(r.providers))
	for _, provider := range r.providers {
		if query.CityCode != "" && provider.CityCode != query.CityCode {
			continue
		}
		if !r.matchCategory(provider.ID, query.CategoryID) {
			continue
		}
		if !r.matchKeyword(provider, query.Keyword) {
			continue
		}
		matched = append(matched, r.withOnline(provider))
	}

	sort.Slice(matched, func(i, j int) bool {
		if matched[i].RatingAvg == matched[j].RatingAvg {
			return matched[i].CompletedOrders > matched[j].CompletedOrders
		}
		return matched[i].RatingAvg > matched[j].RatingAvg
	})

	total := len(matched)
	start := (query.Page - 1) * query.PageSize
	if start >= total {
		return []domain.ProviderPublicProfile{}, total, nil
	}
	end := start + query.PageSize
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (r *ServiceDiscoveryRepository) GetPublicProviderByID(_ context.Context, providerID string) (domain.ProviderPublicProfile, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	provider, ok := r.providers[providerID]
	if !ok {
		return domain.ProviderPublicProfile{}, domain.ErrProviderNotFound
	}
	return r.withOnline(provider), nil
}

func (r *ServiceDiscoveryRepository) ListProviderServiceItems(_ context.Context, providerID string) ([]domain.ServiceItem, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if _, ok := r.providers[providerID]; !ok {
		return nil, domain.ErrProviderNotFound
	}

	items := make([]domain.ServiceItem, 0)
	for _, item := range r.items {
		if item.ProviderID != providerID {
			continue
		}
		items = append(items, item)
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].SortOrder < items[j].SortOrder
	})
	return items, nil
}

func (r *ServiceDiscoveryRepository) matchCategory(providerID string, categoryID string) bool {
	if categoryID == "" || categoryID == "cat_all" {
		return true
	}
	for _, item := range r.items {
		if item.ProviderID == providerID && item.CategoryID == categoryID {
			return true
		}
	}
	return false
}

func (r *ServiceDiscoveryRepository) matchKeyword(provider domain.ProviderPublicProfile, keyword string) bool {
	k := strings.TrimSpace(strings.ToLower(keyword))
	if k == "" {
		return true
	}
	if strings.Contains(strings.ToLower(provider.DisplayName), k) || strings.Contains(strings.ToLower(provider.Bio), k) {
		return true
	}
	for _, tag := range provider.Tags {
		if strings.Contains(strings.ToLower(tag), k) {
			return true
		}
	}
	for _, item := range r.items {
		if item.ProviderID == provider.ID && strings.Contains(strings.ToLower(item.Title), k) {
			return true
		}
	}
	return false
}

func (r *ServiceDiscoveryRepository) TouchProviderHeartbeat(providerID string, now time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.providers[providerID]; !ok {
		return domain.ErrProviderNotFound
	}
	r.presence[providerID] = now
	return nil
}

func (r *ServiceDiscoveryRepository) withOnline(provider domain.ProviderPublicProfile) domain.ProviderPublicProfile {
	p := provider
	lastSeen, ok := r.presence[p.ID]
	if !ok {
		p.Online = false
		return p
	}
	p.Online = time.Since(lastSeen) <= 5*time.Minute
	return p
}

var _ domain.Repository = (*ServiceDiscoveryRepository)(nil)

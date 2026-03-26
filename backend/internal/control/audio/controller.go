package audio

import (
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"ai-listen/backend/internal/support/authctx"
	"ai-listen/backend/pkg/ecode"
	"ai-listen/backend/pkg/httpx"
	"ai-listen/backend/pkg/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Controller struct {
	logger *zap.Logger
}

type AudioItem struct {
	ID            uint64    `json:"id"`
	Title         string    `json:"title"`
	Category      string    `json:"category"`
	DurationSec   int       `json:"durationSec"`
	Mood          string    `json:"mood"`
	Tags          []string  `json:"tags"`
	CoverURL      string    `json:"coverUrl"`
	AudioURL      string    `json:"audioUrl"`
	PlayCount     int       `json:"playCount"`
	FavoriteCount int       `json:"favoriteCount"`
	CreatedAt     time.Time `json:"createdAt"`
}

type audioState struct {
	mu         sync.RWMutex
	items      map[uint64]*AudioItem
	favorites  map[uint64]map[uint64]struct{}
	playRecord map[uint64]map[uint64]int
}

var audioStore = newAudioState()

func newAudioState() *audioState {
	now := time.Now()
	items := []*AudioItem{
		{
			ID:          1,
			Title:       "深夜情绪安抚",
			Category:    "sleep",
			DurationSec: 780,
			Mood:        "calm",
			Tags:        []string{"助眠", "放松", "呼吸"},
			CoverURL:    "https://cdn.listen.dev/audio/sleep-1.jpg",
			AudioURL:    "https://cdn.listen.dev/audio/sleep-1.mp3",
			CreatedAt:   now.Add(-72 * time.Hour),
		},
		{
			ID:          2,
			Title:       "通勤减压白噪音",
			Category:    "focus",
			DurationSec: 640,
			Mood:        "steady",
			Tags:        []string{"通勤", "减压", "专注"},
			CoverURL:    "https://cdn.listen.dev/audio/focus-1.jpg",
			AudioURL:    "https://cdn.listen.dev/audio/focus-1.mp3",
			CreatedAt:   now.Add(-48 * time.Hour),
		},
		{
			ID:          3,
			Title:       "睡前陪聊引导",
			Category:    "companion",
			DurationSec: 920,
			Mood:        "warm",
			Tags:        []string{"陪伴", "情绪", "睡前"},
			CoverURL:    "https://cdn.listen.dev/audio/companion-1.jpg",
			AudioURL:    "https://cdn.listen.dev/audio/companion-1.mp3",
			CreatedAt:   now.Add(-24 * time.Hour),
		},
	}
	state := &audioState{
		items:      make(map[uint64]*AudioItem, len(items)),
		favorites:  map[uint64]map[uint64]struct{}{},
		playRecord: map[uint64]map[uint64]int{},
	}
	for _, item := range items {
		cloned := *item
		cloned.Tags = append([]string(nil), item.Tags...)
		state.items[item.ID] = &cloned
	}
	return state
}

func (s *audioState) categories() []gin.H {
	s.mu.RLock()
	defer s.mu.RUnlock()
	counter := map[string]int{}
	for _, item := range s.items {
		counter[item.Category]++
	}
	result := make([]gin.H, 0, len(counter))
	for key, count := range counter {
		result = append(result, gin.H{
			"key":   key,
			"name":  categoryLabel(key),
			"count": count,
		})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i]["key"].(string) < result[j]["key"].(string)
	})
	return result
}

func (s *audioState) list(category, keyword string) []*AudioItem {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]*AudioItem, 0, len(s.items))
	keyword = strings.TrimSpace(strings.ToLower(keyword))
	for _, item := range s.items {
		if category != "" && item.Category != category {
			continue
		}
		if keyword != "" {
			matched := strings.Contains(strings.ToLower(item.Title), keyword) || strings.Contains(strings.ToLower(item.Mood), keyword)
			if !matched {
				for _, tag := range item.Tags {
					if strings.Contains(strings.ToLower(tag), keyword) {
						matched = true
						break
					}
				}
			}
			if !matched {
				continue
			}
		}
		result = append(result, cloneAudioItem(item))
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].CreatedAt.After(result[j].CreatedAt)
	})
	return result
}

func (s *audioState) detail(id uint64) (*AudioItem, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	item, ok := s.items[id]
	if !ok {
		return nil, false
	}
	return cloneAudioItem(item), true
}

func (s *audioState) recordPlay(userID, audioID uint64, progressSec int) (*AudioItem, int, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	item, ok := s.items[audioID]
	if !ok {
		return nil, 0, false
	}
	item.PlayCount++
	if _, ok := s.playRecord[userID]; !ok {
		s.playRecord[userID] = map[uint64]int{}
	}
	s.playRecord[userID][audioID] = progressSec
	return cloneAudioItem(item), progressSec, true
}

func (s *audioState) markFavorite(userID, audioID uint64) (*AudioItem, bool, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	item, ok := s.items[audioID]
	if !ok {
		return nil, false, false
	}
	if _, ok := s.favorites[userID]; !ok {
		s.favorites[userID] = map[uint64]struct{}{}
	}
	_, existed := s.favorites[userID][audioID]
	if !existed {
		s.favorites[userID][audioID] = struct{}{}
		item.FavoriteCount++
	}
	return cloneAudioItem(item), !existed, true
}

func cloneAudioItem(item *AudioItem) *AudioItem {
	if item == nil {
		return nil
	}
	copied := *item
	copied.Tags = append([]string(nil), item.Tags...)
	return &copied
}

func categoryLabel(key string) string {
	switch key {
	case "sleep":
		return "助眠"
	case "focus":
		return "专注"
	case "companion":
		return "陪伴"
	default:
		return key
	}
}

func (h *Controller) Categories(c *gin.Context) {
	response.Success(c, gin.H{
		"module": "audio",
		"action": "categories",
		"list":   audioStore.categories(),
	})
}

func (h *Controller) List(c *gin.Context) {
	category := c.Query("category")
	keyword := c.Query("keyword")
	list := audioStore.list(category, keyword)
	response.Success(c, gin.H{
		"module":   "audio",
		"action":   "list",
		"query":    gin.H{"category": category, "keyword": keyword, "pagination": httpx.PaginationQuery(c)},
		"total":    len(list),
		"list":     list,
		"hasQuery": category != "" || strings.TrimSpace(keyword) != "",
	})
}

func (h *Controller) Detail(c *gin.Context) {
	audioID, err := strconv.ParseUint(c.Param("audioId"), 10, 64)
	if err != nil || audioID == 0 {
		response.Fail(c, http.StatusBadRequest, ecode.BadRequest, gin.H{"reason": "invalid audioId"})
		return
	}
	item, ok := audioStore.detail(audioID)
	if !ok {
		response.Fail(c, http.StatusNotFound, ecode.NotFound, gin.H{"reason": "audio not found"})
		return
	}
	response.Success(c, gin.H{"module": "audio", "action": "detail", "audio": item})
}

func (h *Controller) PlayLog(c *gin.Context) {
	user, ok := authctx.CurrentUser(c)
	if !ok {
		return
	}
	audioID, err := strconv.ParseUint(c.Param("audioId"), 10, 64)
	if err != nil || audioID == 0 {
		response.Fail(c, http.StatusBadRequest, ecode.BadRequest, gin.H{"reason": "invalid audioId"})
		return
	}
	var req PlayLogRequest
	_ = c.ShouldBindJSON(&req)
	progressSec := req.ProgressSec
	if progressSec == 0 {
		progressSec = req.PositionSec
	}
	item, progress, found := audioStore.recordPlay(user.ID, audioID, progressSec)
	if !found {
		response.Fail(c, http.StatusNotFound, ecode.NotFound, gin.H{"reason": "audio not found"})
		return
	}
	response.Success(c, gin.H{
		"module":      "audio",
		"action":      "play_log",
		"audio":       item,
		"progressSec": progress,
	})
}

func (h *Controller) Favorite(c *gin.Context) {
	user, ok := authctx.CurrentUser(c)
	if !ok {
		return
	}
	audioID, err := strconv.ParseUint(c.Param("audioId"), 10, 64)
	if err != nil || audioID == 0 {
		response.Fail(c, http.StatusBadRequest, ecode.BadRequest, gin.H{"reason": "invalid audioId"})
		return
	}
	item, created, found := audioStore.markFavorite(user.ID, audioID)
	if !found {
		response.Fail(c, http.StatusNotFound, ecode.NotFound, gin.H{"reason": "audio not found"})
		return
	}
	response.Success(c, gin.H{"module": "audio", "action": "favorite", "audio": item, "created": created})
}

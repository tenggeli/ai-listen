package audio

import (
	"ai-listen/backend/pkg/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	logger *zap.Logger
}

func RegisterRoutes(group *gin.RouterGroup, logger *zap.Logger) {
	handler := &Handler{logger: logger}

	group.GET("/audio/categories", handler.Categories)
	group.GET("/audio", handler.List)
	group.GET("/audio/:audioId", handler.Detail)
	group.POST("/audio/:audioId/play-logs", handler.PlayLog)
	group.POST("/audio/:audioId/favorite", handler.Favorite)
}

func (h *Handler) Categories(c *gin.Context) {
	response.Success(c, gin.H{"module": "audio", "action": "categories"})
}

func (h *Handler) List(c *gin.Context) {
	response.Success(c, gin.H{"module": "audio", "action": "list", "category": c.Query("category")})
}

func (h *Handler) Detail(c *gin.Context) {
	response.Success(c, gin.H{"module": "audio", "action": "detail", "audioId": c.Param("audioId")})
}

func (h *Handler) PlayLog(c *gin.Context) {
	response.Success(c, gin.H{"module": "audio", "action": "play_log", "audioId": c.Param("audioId")})
}

func (h *Handler) Favorite(c *gin.Context) {
	response.Success(c, gin.H{"module": "audio", "action": "favorite", "audioId": c.Param("audioId")})
}

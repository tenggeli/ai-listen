package audio

import (
	"ai-listen/backend/pkg/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Controller struct {
	logger *zap.Logger
}

func (h *Controller) Categories(c *gin.Context) {
	response.Success(c, gin.H{"module": "audio", "action": "categories"})
}

func (h *Controller) List(c *gin.Context) {
	response.Success(c, gin.H{"module": "audio", "action": "list", "category": c.Query("category")})
}

func (h *Controller) Detail(c *gin.Context) {
	response.Success(c, gin.H{"module": "audio", "action": "detail", "audioId": c.Param("audioId")})
}

func (h *Controller) PlayLog(c *gin.Context) {
	response.Success(c, gin.H{"module": "audio", "action": "play_log", "audioId": c.Param("audioId")})
}

func (h *Controller) Favorite(c *gin.Context) {
	response.Success(c, gin.H{"module": "audio", "action": "favorite", "audioId": c.Param("audioId")})
}

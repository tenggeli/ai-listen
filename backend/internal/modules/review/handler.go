package review

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

	group.POST("/orders/:orderId/reviews", handler.Create)
	group.GET("/orders/:orderId/reviews", handler.List)
}

func (h *Handler) Create(c *gin.Context) {
	var req CreateReviewRequest
	_ = c.ShouldBindJSON(&req)
	response.Success(c, gin.H{"module": "review", "action": "create", "orderId": c.Param("orderId"), "request": req})
}

func (h *Handler) List(c *gin.Context) {
	response.Success(c, gin.H{"module": "review", "action": "list", "orderId": c.Param("orderId")})
}

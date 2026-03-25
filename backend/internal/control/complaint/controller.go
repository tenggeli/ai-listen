package complaint

import (
	"ai-listen/backend/pkg/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Controller struct {
	logger *zap.Logger
}

func RegisterRoutes(group *gin.RouterGroup, logger *zap.Logger) {
	handler := &Controller{logger: logger}

	group.POST("/orders/:orderId/complaints", handler.Create)
	group.GET("/complaints/:complaintId", handler.Detail)
}

func (h *Controller) Create(c *gin.Context) {
	var req CreateComplaintRequest
	_ = c.ShouldBindJSON(&req)
	response.Success(c, gin.H{"module": "complaint", "action": "create", "orderId": c.Param("orderId"), "request": req})
}

func (h *Controller) Detail(c *gin.Context) {
	response.Success(c, gin.H{"module": "complaint", "action": "detail", "complaintId": c.Param("complaintId")})
}

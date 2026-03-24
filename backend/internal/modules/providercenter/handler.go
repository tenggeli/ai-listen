package providercenter

import (
	"ai-listen/backend/pkg/httpx"
	"ai-listen/backend/pkg/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	logger *zap.Logger
}

func RegisterRoutes(group *gin.RouterGroup, logger *zap.Logger) {
	handler := &Handler{logger: logger}

	centerGroup := group.Group("/provider-center")
	centerGroup.POST("/apply", handler.Apply)
	centerGroup.GET("/apply/status", handler.ApplyStatus)
	centerGroup.PUT("/profile", handler.UpdateProfile)
	centerGroup.PUT("/service-items", handler.UpdateServiceItems)
	centerGroup.PUT("/work-status", handler.UpdateWorkStatus)
	centerGroup.GET("/orders", handler.GetOrders)
	centerGroup.POST("/invitations", handler.CreateInvitation)
	centerGroup.GET("/wallet", handler.GetWallet)
	centerGroup.POST("/accounts", handler.CreateAccount)
	centerGroup.GET("/accounts", handler.ListAccounts)
	centerGroup.POST("/withdraws", handler.CreateWithdraw)
	centerGroup.GET("/withdraws", handler.ListWithdraws)
	centerGroup.POST("/orders/:orderId/accept", handler.AcceptOrder)
	centerGroup.POST("/orders/:orderId/depart", handler.DepartOrder)
	centerGroup.POST("/orders/:orderId/arrive", handler.ArriveOrder)
	centerGroup.POST("/orders/:orderId/finish", handler.FinishOrder)
}

func (h *Handler) Apply(c *gin.Context) {
	response.Success(c, gin.H{"module": "provider_center", "action": "apply"})
}

func (h *Handler) ApplyStatus(c *gin.Context) {
	response.Success(c, gin.H{"module": "provider_center", "action": "apply_status"})
}

func (h *Handler) UpdateProfile(c *gin.Context) {
	httpx.NotImplemented(c, "provider_center.update_profile")
}

func (h *Handler) UpdateServiceItems(c *gin.Context) {
	httpx.NotImplemented(c, "provider_center.update_service_items")
}

func (h *Handler) UpdateWorkStatus(c *gin.Context) {
	var req WorkStatusRequest
	_ = c.ShouldBindJSON(&req)
	response.Success(c, gin.H{"module": "provider_center", "action": "update_work_status", "request": req})
}

func (h *Handler) GetOrders(c *gin.Context) {
	response.Success(c, gin.H{"module": "provider_center", "action": "get_orders", "query": httpx.PaginationQuery(c)})
}

func (h *Handler) CreateInvitation(c *gin.Context) {
	httpx.NotImplemented(c, "provider_center.create_invitation")
}

func (h *Handler) GetWallet(c *gin.Context) {
	response.Success(c, gin.H{"module": "provider_center", "action": "get_wallet"})
}

func (h *Handler) CreateAccount(c *gin.Context) {
	var req AccountRequest
	_ = c.ShouldBindJSON(&req)
	response.Success(c, gin.H{"module": "provider_center", "action": "create_account", "request": req})
}

func (h *Handler) ListAccounts(c *gin.Context) {
	response.Success(c, gin.H{"module": "provider_center", "action": "list_accounts"})
}

func (h *Handler) CreateWithdraw(c *gin.Context) {
	var req WithdrawRequest
	_ = c.ShouldBindJSON(&req)
	response.Success(c, gin.H{"module": "provider_center", "action": "create_withdraw", "request": req})
}

func (h *Handler) ListWithdraws(c *gin.Context) {
	response.Success(c, gin.H{"module": "provider_center", "action": "list_withdraws", "query": httpx.PaginationQuery(c)})
}

func (h *Handler) AcceptOrder(c *gin.Context) {
	response.Success(c, gin.H{"module": "provider_center", "action": "accept_order", "orderId": c.Param("orderId")})
}

func (h *Handler) DepartOrder(c *gin.Context) {
	response.Success(c, gin.H{"module": "provider_center", "action": "depart_order", "orderId": c.Param("orderId")})
}

func (h *Handler) ArriveOrder(c *gin.Context) {
	response.Success(c, gin.H{"module": "provider_center", "action": "arrive_order", "orderId": c.Param("orderId")})
}

func (h *Handler) FinishOrder(c *gin.Context) {
	response.Success(c, gin.H{"module": "provider_center", "action": "finish_order", "orderId": c.Param("orderId")})
}

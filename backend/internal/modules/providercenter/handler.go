package providercenter

import (
	"strconv"

	"ai-listen/backend/internal/store"
	"ai-listen/backend/internal/support/authctx"
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
	user, ok := authctx.CurrentUser(c)
	if !ok {
		return
	}
	var req ApplyRequest
	_ = c.ShouldBindJSON(&req)
	provider, err := store.Default().ApplyProvider(user.ID, req.RealName, req.IDCardNo)
	if err != nil {
		httpx.NotImplemented(c, "provider_center.apply")
		return
	}
	response.Success(c, gin.H{"module": "provider_center", "action": "apply", "provider": provider})
}

func (h *Handler) ApplyStatus(c *gin.Context) {
	user, ok := authctx.CurrentUser(c)
	if !ok {
		return
	}
	provider, err := store.Default().ProviderByUserID(user.ID)
	if err != nil {
		response.Success(c, gin.H{"module": "provider_center", "action": "apply_status", "provider": nil})
		return
	}
	response.Success(c, gin.H{"module": "provider_center", "action": "apply_status", "provider": provider})
}

func (h *Handler) UpdateProfile(c *gin.Context) {
	user, ok := authctx.CurrentUser(c)
	if !ok {
		return
	}
	var req UpdateProfileRequest
	_ = c.ShouldBindJSON(&req)
	provider, err := store.Default().UpdateProviderProfile(user.ID, req.DisplayName, req.Intro, req.Tags)
	if err != nil {
		httpx.NotImplemented(c, "provider_center.update_profile")
		return
	}
	response.Success(c, gin.H{"module": "provider_center", "action": "update_profile", "provider": provider})
}

func (h *Handler) UpdateServiceItems(c *gin.Context) {
	user, ok := authctx.CurrentUser(c)
	if !ok {
		return
	}
	var req UpdateServiceItemsRequest
	_ = c.ShouldBindJSON(&req)
	items := make([]store.ProviderServiceItem, 0, len(req.Items))
	for _, item := range req.Items {
		items = append(items, store.ProviderServiceItem{
			ServiceItemID: item.ServiceItemID,
			PriceAmount:   item.PriceAmount,
			PriceUnit:     item.PriceUnit,
		})
	}
	provider, err := store.Default().UpdateProviderServiceItems(user.ID, items)
	if err != nil {
		httpx.NotImplemented(c, "provider_center.update_service_items")
		return
	}
	response.Success(c, gin.H{"module": "provider_center", "action": "update_service_items", "provider": provider})
}

func (h *Handler) UpdateWorkStatus(c *gin.Context) {
	user, ok := authctx.CurrentUser(c)
	if !ok {
		return
	}
	var req WorkStatusRequest
	_ = c.ShouldBindJSON(&req)
	provider, err := store.Default().UpdateProviderWorkStatus(user.ID, req.WorkStatus)
	if err != nil {
		httpx.NotImplemented(c, "provider_center.update_work_status")
		return
	}
	response.Success(c, gin.H{"module": "provider_center", "action": "update_work_status", "provider": provider})
}

func (h *Handler) GetOrders(c *gin.Context) {
	user, ok := authctx.CurrentUser(c)
	if !ok {
		return
	}
	list, err := store.Default().OrdersByProvider(user.ID)
	if err != nil {
		response.Success(c, gin.H{"module": "provider_center", "action": "get_orders", "list": []any{}})
		return
	}
	response.Success(c, gin.H{"module": "provider_center", "action": "get_orders", "query": httpx.PaginationQuery(c), "list": list})
}

func (h *Handler) CreateInvitation(c *gin.Context) {
	httpx.NotImplemented(c, "provider_center.create_invitation")
}

func (h *Handler) GetWallet(c *gin.Context) {
	user, ok := authctx.CurrentUser(c)
	if !ok {
		return
	}
	provider, err := store.Default().ProviderByUserID(user.ID)
	if err != nil {
		httpx.NotImplemented(c, "provider_center.get_wallet")
		return
	}
	response.Success(c, gin.H{"module": "provider_center", "action": "get_wallet", "providerId": provider.ID, "balanceAmount": 0})
}

func (h *Handler) CreateAccount(c *gin.Context) {
	_, ok := authctx.CurrentUser(c)
	if !ok {
		return
	}
	var req AccountRequest
	_ = c.ShouldBindJSON(&req)
	response.Success(c, gin.H{"module": "provider_center", "action": "create_account", "request": req})
}

func (h *Handler) ListAccounts(c *gin.Context) {
	_, ok := authctx.CurrentUser(c)
	if !ok {
		return
	}
	response.Success(c, gin.H{"module": "provider_center", "action": "list_accounts"})
}

func (h *Handler) CreateWithdraw(c *gin.Context) {
	_, ok := authctx.CurrentUser(c)
	if !ok {
		return
	}
	var req WithdrawRequest
	_ = c.ShouldBindJSON(&req)
	response.Success(c, gin.H{"module": "provider_center", "action": "create_withdraw", "request": req})
}

func (h *Handler) ListWithdraws(c *gin.Context) {
	_, ok := authctx.CurrentUser(c)
	if !ok {
		return
	}
	response.Success(c, gin.H{"module": "provider_center", "action": "list_withdraws", "query": httpx.PaginationQuery(c)})
}

func (h *Handler) AcceptOrder(c *gin.Context) {
	user, ok := authctx.CurrentUser(c)
	if !ok {
		return
	}
	orderID, _ := strconv.ParseUint(c.Param("orderId"), 10, 64)
	order, err := store.Default().ProviderAcceptOrder(user.ID, orderID)
	if err != nil {
		httpx.NotImplemented(c, "provider_center.accept_order")
		return
	}
	response.Success(c, gin.H{"module": "provider_center", "action": "accept_order", "order": order})
}

func (h *Handler) DepartOrder(c *gin.Context) {
	user, ok := authctx.CurrentUser(c)
	if !ok {
		return
	}
	orderID, _ := strconv.ParseUint(c.Param("orderId"), 10, 64)
	order, err := store.Default().ProviderDepartOrder(user.ID, orderID)
	if err != nil {
		httpx.NotImplemented(c, "provider_center.depart_order")
		return
	}
	response.Success(c, gin.H{"module": "provider_center", "action": "depart_order", "order": order})
}

func (h *Handler) ArriveOrder(c *gin.Context) {
	user, ok := authctx.CurrentUser(c)
	if !ok {
		return
	}
	orderID, _ := strconv.ParseUint(c.Param("orderId"), 10, 64)
	order, err := store.Default().ProviderArriveOrder(user.ID, orderID)
	if err != nil {
		httpx.NotImplemented(c, "provider_center.arrive_order")
		return
	}
	response.Success(c, gin.H{"module": "provider_center", "action": "arrive_order", "order": order})
}

func (h *Handler) FinishOrder(c *gin.Context) {
	user, ok := authctx.CurrentUser(c)
	if !ok {
		return
	}
	orderID, _ := strconv.ParseUint(c.Param("orderId"), 10, 64)
	order, err := store.Default().ProviderFinishOrder(user.ID, orderID)
	if err != nil {
		httpx.NotImplemented(c, "provider_center.finish_order")
		return
	}
	response.Success(c, gin.H{"module": "provider_center", "action": "finish_order", "order": order})
}

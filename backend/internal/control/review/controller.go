package review

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"ai-listen/backend/internal/model"
	"ai-listen/backend/internal/support/authctx"
	"ai-listen/backend/pkg/ecode"
	"ai-listen/backend/pkg/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Controller struct {
	logger *zap.Logger
}

func (h *Controller) Create(c *gin.Context) {
	user, ok := authctx.CurrentUser(c)
	if !ok {
		return
	}
	orderID, err := strconv.ParseUint(c.Param("orderId"), 10, 64)
	if err != nil || orderID == 0 {
		response.Fail(c, http.StatusBadRequest, ecode.BadRequest, gin.H{"reason": "orderId is invalid"})
		return
	}

	var req CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.Content) == "" || req.Score < 1 || req.Score > 10 {
		response.Fail(c, http.StatusBadRequest, ecode.BadRequest, gin.H{"reason": "score(1-10) and content are required"})
		return
	}

	review, err := model.Default().CreateReview(user.ID, orderID, model.CreateReviewInput{
		Score:       req.Score,
		Content:     strings.TrimSpace(req.Content),
		Images:      req.Images,
		IsAnonymous: req.IsAnonymous,
	})
	if err != nil {
		h.handleReviewError(c, err)
		return
	}
	response.Success(c, gin.H{"module": "review", "action": "create", "review": review})
}

func (h *Controller) List(c *gin.Context) {
	orderID, err := strconv.ParseUint(c.Param("orderId"), 10, 64)
	if err != nil || orderID == 0 {
		response.Fail(c, http.StatusBadRequest, ecode.BadRequest, gin.H{"reason": "orderId is invalid"})
		return
	}

	list, err := model.Default().ReviewsByOrder(orderID)
	if err != nil {
		h.handleReviewError(c, err)
		return
	}
	response.Success(c, gin.H{"module": "review", "action": "list", "orderId": orderID, "list": list})
}

func (h *Controller) handleReviewError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, model.ErrOrderNotFound):
		response.Fail(c, http.StatusNotFound, ecode.OrderNotFound, gin.H{"reason": err.Error()})
	case errors.Is(err, model.ErrUnauthorized):
		response.Fail(c, http.StatusForbidden, ecode.Forbidden, gin.H{"reason": err.Error()})
	case errors.Is(err, model.ErrReviewAlreadyExists):
		response.Fail(c, http.StatusConflict, ecode.Conflict, gin.H{"reason": err.Error()})
	case errors.Is(err, model.ErrInvalidOrderStatus):
		response.Fail(c, http.StatusConflict, ecode.OrderStateInvalid, gin.H{"reason": err.Error()})
	default:
		response.Fail(c, http.StatusInternalServerError, ecode.InternalServerError, gin.H{"reason": err.Error()})
	}
}

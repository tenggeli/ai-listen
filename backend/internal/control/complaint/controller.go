package complaint

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

	var req CreateComplaintRequest
	if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.ComplaintType) == "" || strings.TrimSpace(req.Content) == "" {
		response.Fail(c, http.StatusBadRequest, ecode.BadRequest, gin.H{"reason": "complaintType and content are required"})
		return
	}

	complaint, err := model.Default().CreateComplaint(user.ID, orderID, model.CreateComplaintInput{
		ComplaintType:  strings.TrimSpace(req.ComplaintType),
		Content:        strings.TrimSpace(req.Content),
		EvidenceImages: req.EvidenceImages,
	})
	if err != nil {
		h.handleComplaintError(c, err)
		return
	}
	response.Success(c, gin.H{"module": "complaint", "action": "create", "complaint": complaint})
}

func (h *Controller) Detail(c *gin.Context) {
	complaintID, err := strconv.ParseUint(c.Param("complaintId"), 10, 64)
	if err != nil || complaintID == 0 {
		response.Fail(c, http.StatusBadRequest, ecode.BadRequest, gin.H{"reason": "complaintId is invalid"})
		return
	}

	complaint, err := model.Default().GetComplaint(complaintID)
	if err != nil {
		h.handleComplaintError(c, err)
		return
	}
	response.Success(c, gin.H{"module": "complaint", "action": "detail", "complaint": complaint})
}

func (h *Controller) handleComplaintError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, model.ErrComplaintNotFound), errors.Is(err, model.ErrOrderNotFound):
		response.Fail(c, http.StatusNotFound, ecode.NotFound, gin.H{"reason": err.Error()})
	case errors.Is(err, model.ErrUnauthorized):
		response.Fail(c, http.StatusForbidden, ecode.Forbidden, gin.H{"reason": err.Error()})
	default:
		response.Fail(c, http.StatusInternalServerError, ecode.InternalServerError, gin.H{"reason": err.Error()})
	}
}

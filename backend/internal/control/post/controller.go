package post

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"ai-listen/backend/internal/model"
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

func (h *Controller) Create(c *gin.Context) {
	user, ok := authctx.CurrentUser(c)
	if !ok {
		return
	}
	var req CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.Content) == "" {
		response.Fail(c, http.StatusBadRequest, ecode.BadRequest, gin.H{"reason": "content is required"})
		return
	}
	if req.VisibleScope == 0 {
		req.VisibleScope = 1
	}

	post, err := model.Default().CreatePost(user.ID, model.CreatePostInput{
		Content:      strings.TrimSpace(req.Content),
		Topic:        strings.TrimSpace(req.Topic),
		CityCode:     strings.TrimSpace(req.CityCode),
		VisibleScope: req.VisibleScope,
	})
	if err != nil {
		h.handlePostError(c, err)
		return
	}
	response.Success(c, gin.H{"module": "post", "action": "create", "post": post})
}

func (h *Controller) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	list, err := model.Default().Posts(page, pageSize)
	if err != nil {
		h.handlePostError(c, err)
		return
	}
	response.Success(c, gin.H{"module": "post", "action": "list", "query": httpx.PaginationQuery(c), "list": list})
}

func (h *Controller) Detail(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("postId"), 10, 64)
	if err != nil || postID == 0 {
		response.Fail(c, http.StatusBadRequest, ecode.BadRequest, gin.H{"reason": "postId is invalid"})
		return
	}

	post, err := model.Default().GetPost(postID)
	if err != nil {
		h.handlePostError(c, err)
		return
	}
	response.Success(c, gin.H{"module": "post", "action": "detail", "post": post})
}

func (h *Controller) Comment(c *gin.Context) {
	user, ok := authctx.CurrentUser(c)
	if !ok {
		return
	}
	postID, err := strconv.ParseUint(c.Param("postId"), 10, 64)
	if err != nil || postID == 0 {
		response.Fail(c, http.StatusBadRequest, ecode.BadRequest, gin.H{"reason": "postId is invalid"})
		return
	}

	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.Content) == "" {
		response.Fail(c, http.StatusBadRequest, ecode.BadRequest, gin.H{"reason": "content is required"})
		return
	}

	comment, post, err := model.Default().CreatePostComment(user.ID, postID, model.CreatePostCommentInput{
		Content:        strings.TrimSpace(req.Content),
		ReplyCommentID: req.ReplyCommentID,
	})
	if err != nil {
		h.handlePostError(c, err)
		return
	}
	response.Success(c, gin.H{"module": "post", "action": "comment", "comment": comment, "post": post})
}

func (h *Controller) Like(c *gin.Context) {
	user, ok := authctx.CurrentUser(c)
	if !ok {
		return
	}
	postID, err := strconv.ParseUint(c.Param("postId"), 10, 64)
	if err != nil || postID == 0 {
		response.Fail(c, http.StatusBadRequest, ecode.BadRequest, gin.H{"reason": "postId is invalid"})
		return
	}

	post, err := model.Default().LikePost(user.ID, postID)
	if err != nil {
		h.handlePostError(c, err)
		return
	}
	response.Success(c, gin.H{"module": "post", "action": "like", "post": post})
}

func (h *Controller) Unlike(c *gin.Context) {
	user, ok := authctx.CurrentUser(c)
	if !ok {
		return
	}
	postID, err := strconv.ParseUint(c.Param("postId"), 10, 64)
	if err != nil || postID == 0 {
		response.Fail(c, http.StatusBadRequest, ecode.BadRequest, gin.H{"reason": "postId is invalid"})
		return
	}

	post, err := model.Default().UnlikePost(user.ID, postID)
	if err != nil {
		h.handlePostError(c, err)
		return
	}
	response.Success(c, gin.H{"module": "post", "action": "unlike", "post": post})
}

func (h *Controller) handlePostError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, model.ErrPostNotFound):
		response.Fail(c, http.StatusNotFound, ecode.NotFound, gin.H{"reason": err.Error()})
	default:
		response.Fail(c, http.StatusInternalServerError, ecode.InternalServerError, gin.H{"reason": err.Error()})
	}
}

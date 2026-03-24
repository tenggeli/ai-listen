package post

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

	group.POST("/posts", handler.Create)
	group.GET("/posts", handler.List)
	group.GET("/posts/:postId", handler.Detail)
	group.POST("/posts/:postId/comments", handler.Comment)
	group.POST("/posts/:postId/likes", handler.Like)
	group.DELETE("/posts/:postId/likes", handler.Unlike)
}

func (h *Handler) Create(c *gin.Context) {
	httpx.NotImplemented(c, "post.create")
}

func (h *Handler) List(c *gin.Context) {
	response.Success(c, gin.H{"module": "post", "action": "list", "query": httpx.PaginationQuery(c)})
}

func (h *Handler) Detail(c *gin.Context) {
	response.Success(c, gin.H{"module": "post", "action": "detail", "postId": c.Param("postId")})
}

func (h *Handler) Comment(c *gin.Context) {
	httpx.NotImplemented(c, "post.comment")
}

func (h *Handler) Like(c *gin.Context) {
	response.Success(c, gin.H{"module": "post", "action": "like", "postId": c.Param("postId")})
}

func (h *Handler) Unlike(c *gin.Context) {
	response.Success(c, gin.H{"module": "post", "action": "unlike", "postId": c.Param("postId")})
}

package post

import (
	"ai-listen/backend/pkg/httpx"
	"ai-listen/backend/pkg/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Controller struct {
	logger *zap.Logger
}

func RegisterRoutes(group *gin.RouterGroup, logger *zap.Logger) {
	handler := &Controller{logger: logger}

	group.POST("/posts", handler.Create)
	group.GET("/posts", handler.List)
	group.GET("/posts/:postId", handler.Detail)
	group.POST("/posts/:postId/comments", handler.Comment)
	group.POST("/posts/:postId/likes", handler.Like)
	group.DELETE("/posts/:postId/likes", handler.Unlike)
}

func (h *Controller) Create(c *gin.Context) {
	httpx.NotImplemented(c, "post.create")
}

func (h *Controller) List(c *gin.Context) {
	response.Success(c, gin.H{"module": "post", "action": "list", "query": httpx.PaginationQuery(c)})
}

func (h *Controller) Detail(c *gin.Context) {
	response.Success(c, gin.H{"module": "post", "action": "detail", "postId": c.Param("postId")})
}

func (h *Controller) Comment(c *gin.Context) {
	httpx.NotImplemented(c, "post.comment")
}

func (h *Controller) Like(c *gin.Context) {
	response.Success(c, gin.H{"module": "post", "action": "like", "postId": c.Param("postId")})
}

func (h *Controller) Unlike(c *gin.Context) {
	response.Success(c, gin.H{"module": "post", "action": "unlike", "postId": c.Param("postId")})
}

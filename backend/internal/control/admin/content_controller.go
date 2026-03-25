package admin

import (
	"ai-listen/backend/pkg/httpx"
	"ai-listen/backend/pkg/response"
	"github.com/gin-gonic/gin"
)

func (h *Controller) ListPosts(c *gin.Context) {
	response.Success(c, gin.H{"module": "admin", "action": "list_posts", "query": httpx.PaginationQuery(c)})
}

func (h *Controller) HidePost(c *gin.Context) {
	response.Success(c, gin.H{"module": "admin", "action": "hide_post", "postId": c.Param("postId")})
}

func (h *Controller) ListAudio(c *gin.Context) {
	response.Success(c, gin.H{"module": "admin", "action": "list_audio", "query": httpx.PaginationQuery(c)})
}

func (h *Controller) OffShelfAudio(c *gin.Context) {
	response.Success(c, gin.H{"module": "admin", "action": "off_shelf_audio", "audioId": c.Param("audioId")})
}

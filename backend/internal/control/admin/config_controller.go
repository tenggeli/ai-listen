package admin

import (
	"ai-listen/backend/pkg/response"
	"github.com/gin-gonic/gin"
)

func (h *Controller) ListConfigs(c *gin.Context) {
	response.Success(c, gin.H{"module": "admin", "action": "list_configs"})
}

func (h *Controller) UpdateConfig(c *gin.Context) {
	response.Success(c, gin.H{"module": "admin", "action": "update_config", "configKey": c.Param("configKey")})
}

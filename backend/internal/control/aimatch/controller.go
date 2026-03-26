package aimatch

import (
	"ai-listen/backend/pkg/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Controller struct {
	logger *zap.Logger
}

func (h *Controller) Match(c *gin.Context) {
	var req MatchRequest
	_ = c.ShouldBindJSON(&req)
	response.Success(c, gin.H{
		"module":  "ai",
		"action":  "match",
		"request": req,
		"recommendations": []gin.H{
			{"providerId": 1001, "matchScore": 98},
			{"providerId": 1002, "matchScore": 94},
			{"providerId": 1003, "matchScore": 91},
		},
	})
}

func (h *Controller) GetMatch(c *gin.Context) {
	response.Success(c, gin.H{"module": "ai", "action": "get_match", "sessionId": c.Param("sessionId")})
}

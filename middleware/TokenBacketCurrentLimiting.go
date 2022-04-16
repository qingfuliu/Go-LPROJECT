package middleware

import (
	"MFile/generate/limiter"
	"github.com/gin-gonic/gin"
)

func CurrentLimiting(c *gin.Context) {
	type tokenAccess struct {
		CountToken int64 `json:"countToken" bind:"required"`
	}
	tokenAccess_ := &tokenAccess{}
	if err := c.ShouldBindJSON(tokenAccess_); err == nil {
		limiter.AcquireN(tokenAccess_.CountToken)
	} else {
		limiter.Acquire()
	}
	c.Next()
}

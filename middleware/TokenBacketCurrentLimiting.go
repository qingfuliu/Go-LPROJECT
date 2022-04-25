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

func CurrentLimiterRedis(c *gin.Context){
	type tokenAccess struct {
		CountToken int `json:"countToken" bind:"required"`
	}
	tokenAccess_ := &tokenAccess{}
	var ok bool
	if err := c.ShouldBindJSON(tokenAccess_); err == nil {
		ok=redis_limiter.AllowN(tokenAccess_.CountToken)
	} else {
		ok=redis_limiter.Allow()
	}
	if !ok{
		c.Abort()
	}
	c.Next()
}

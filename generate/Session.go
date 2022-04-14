package generate

import (
	"github.com/gin-gonic/gin"
	"time"
)

type Session interface {
	Get(key string) (interface{}, bool)
	Set(key string) bool
	SetExpire(minute int)
	GetExpire() time.Duration
	AddTo(c *gin.Context)
	String() string
	Json() string
}

type session struct {
	Vars   map[string]interface{} `json:"Vars"`
	Id     string                 `json:"-"`
	Expire time.Duration          `json:"-"`
}

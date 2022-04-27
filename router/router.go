package router

import (
	"MFile/controller"
	"MFile/middleware"
	"github.com/DeanThompson/ginpprof"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var Engine *gin.Engine = nil

func init() {
	Engine = gin.New()
	Engine.Use(gin.Recovery())
	Engine.Use(gin.Logger())
	Engine.Use(cors.Default())
	Engine.Use(middleware.CurrentLimiterRedis)
	Engine.GET("/UpLoadFile", controller.UpLoadFile)
	Engine.GET("/BackPointStart", controller.BackPointStart)
	Engine.GET("/BackPointProcess", controller.BackPointProcess)
	Engine.GET("/MergeFileChunk", controller.MergeFileChunk)
	Engine.GET("/DownLoadFile", controller.DownLoadFile)
	ginpprof.Wrap(Engine)

}

func logger(c *gin.Context) {

}

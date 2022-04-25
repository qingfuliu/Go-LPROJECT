package router

import (
	"MFile/controller"
	"MFile/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var Engine *gin.Engine = nil

func init() {
	Engine = gin.New()
	Engine.Use(gin.Recovery())
	Engine.Use(gin.Logger())
	Engine.Use(cors.Default())
	Engine.Use(middleware.CurrentLimiting)
	Engine.POST("/UpLoadFile", controller.UpLoadFile)
	Engine.POST("/BackPointStart", controller.BackPointStart)
	Engine.POST("/BackPointProcess", controller.BackPointProcess)
	Engine.POST("/MergeFileChunk", controller.MergeFileChunk)
	Engine.POST("/DownLoadFile", controller.DownLoadFile)
}

func logger(c *gin.Context) {

}

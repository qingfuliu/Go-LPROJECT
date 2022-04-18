package router

import (
	"MFile/controller"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
)

var Engine *gin.Engine = nil

func init() {
	Engine = gin.New()
	Engine.Use(gin.Recovery())
	Engine.Use(gin.Logger())
	Engine.Use(cors.Default())
	//Engine.Use(middleware.CurrentLimiting)
	//Engine.StaticFS("/file", http.Dir("E:\\书籍"))
	//Engine.GET("/like", controller.Likes)
	Engine.POST("/like", controller.Likes)
	Engine.POST("/UpLoadFile", controller.UpLoadFile)
	Engine.POST("/BackPointStart", controller.BackPointStart)
	Engine.POST("/BackPointProcess", controller.BackPointProcess)
	Engine.POST("/DownLoadFile", controller.DownLoadFile)
	Engine.StaticFS("/", http.Dir("E:\\书籍"))
	//versionApi := Engine.Group("/api/v1")
	//
	//grantApi := versionApi.Group("/grant")

}

func logger(c *gin.Context) {

}

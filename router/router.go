package router

import (
	"MFile/controller"
<<<<<<< HEAD
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
=======
	"MFile/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
>>>>>>> master
)

var Engine *gin.Engine = nil

func init() {
	Engine = gin.New()
	Engine.Use(gin.Recovery())
	Engine.Use(gin.Logger())
	Engine.Use(cors.Default())
<<<<<<< HEAD
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

=======
	Engine.Use(middleware.CurrentLimiting)
	Engine.POST("/UpLoadFile", controller.UpLoadFile)
	Engine.POST("/BackPointStart", controller.BackPointStart)
	Engine.POST("/BackPointProcess", controller.BackPointProcess)
	Engine.POST("/MergeFileChunk", controller.MergeFileChunk)
	Engine.POST("/DownLoadFile", controller.DownLoadFile)
>>>>>>> master
}

func logger(c *gin.Context) {

}

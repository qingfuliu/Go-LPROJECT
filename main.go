package main

import (
	"MFile/logger"
	"MFile/router"
	"go.uber.org/zap"
)

func main() {
	err := router.Engine.Run("127.0.0.1:5200")
	logger.MLogger.Info("gin err:", zap.Error(err))
}

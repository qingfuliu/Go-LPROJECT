package mysql

import (
	"MFile/logger"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
)

const (
//mysqlConnFormat = "%s:%s@tcp(%s:%d)/%s?parseTime=true"
//username        = "lqf"
//password        = "Wangfei222@"
//port            = 3306
//ip              = "192.168.1.103"
//dbName          = "blogs"
)

var MysqlDb *gorm.DB

func init() {
	var err error
	userName := "lqf"
	passWord := "Wangfei222@"
	ip := "192.168.1.103"
	port := 3306
	dbName := "blogs"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", userName, passWord, ip, port, dbName)
	MysqlDb, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.MLogger.Fatal("mysql connect failed!,err info is :", zap.Error(err))
		os.Exit(1)
	}

	logger.MLogger.Info("mysql connect successful")
}

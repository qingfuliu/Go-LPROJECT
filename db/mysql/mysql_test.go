package mysql

import (
	"MFile/models"
	"fmt"
	"log"
	"testing"
)

//select table_name as name from INFORMATION_SCHEMA.TABLES ;
type test struct {
	TableName_ string `json:"table_name" gorm:"column:TABLE_NAME"`
}

func (*test) TableName() string {
	return "INFORMATION_SCHEMA.TABLES"
}

func TestMysql(t *testing.T) {
	var tests []test
	if err := MysqlDb.Table("INFORMATION_SCHEMA.TABLES").Select("TABLE_NAME").Limit(10).Offset(0).Find(&tests).Error; err != nil {
		t.Fatal(err)
	}
	fmt.Println(tests)
}

func TestMysql1(t *testing.T) {

	test := models.FileChunkInFo{
		FileName:    "lqf111111",
		FileNameMd5: "asdasdfadsf",
		FileMd5:     "fileMd5",
		ChunkTotal:  10,
		FileExt:     ".png",
		ChunkNext:   1,
	}
	err := MysqlDb.Table("BackPointInfo").Where("fileName= ? ", test.FileName).Update("chunkNext", 1).Error
	//err := MysqlDb.Table("BackPointInfo").Select("*").Where("fileName=?", "lqf111").Set("gorm:query_option", "for update").FirstOrCreate(&test).Error
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(test)
	}
}

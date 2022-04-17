package mysql

import (
	"fmt"
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

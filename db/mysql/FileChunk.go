package mysql

import (
	"MFile/models"
)

func FindFirstOrCreate(fileChunkInfo models.FileChunkInFo) (*models.FileChunkInFo, error) {
	err := MysqlDb.Set("gorm:query_option", "for update").
		Where("filename=?", fileChunkInfo.FileName).
		FirstOrCreate(&fileChunkInfo).Error
	if err != nil {
		//	tx.Rollback()
		return nil, err
	}
	//	tx.Commit()
	return &fileChunkInfo, nil
}

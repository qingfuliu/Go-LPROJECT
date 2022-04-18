package controller

import (
	"MFile/config"
	"MFile/db/mysql"
	"MFile/generate/hash"
	logic "MFile/logic/file"
	"MFile/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"mime/multipart"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"
)

func UpLoadFile(c *gin.Context) {
	fileHeaders, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": "failed"})
		return
	}
	if errs := logic.UpLoadFile(fileHeaders); errs != nil {
		c.JSON(http.StatusOK, gin.H{"code": "failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": "ok"})
}

func DownLoadFile(c *gin.Context) {
	fileHeaders, err := c.MultipartForm()
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusOK, gin.H{"code": "failed"})
		return
	}
	now := time.Now().Format("20060102181209")
	for key, _ := range fileHeaders.Value {
		ext := path.Ext(key)
		fileName, _ := hash.Md5(strings.TrimSuffix(key, ext))
		fileName += "_" + now + ext
		c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("atachment;filename=%s", key))
		c.File(path.Join(config.DirName, fileName))
	}
	c.Status(200)
}

func BackPointStart(c *gin.Context) {
	fileName := c.PostForm("fileName")
	fileExt := path.Ext(fileName)
	fileNameMd5 := c.PostForm("fileNameMd5")
	chunkTotal, _ := strconv.Atoi(c.PostForm("chunkTotal"))
	fileMd5 := c.PostForm("fileMd5")

	fileChunkInfo, err := mysql.FindFirstOrCreate(models.FileChunkInFo{
		FileName:    strings.TrimSuffix(fileName, fileExt),
		FileNameMd5: fileNameMd5,
		FileMd5:     fileMd5,
		ChunkTotal:  chunkTotal,
		FileExt:     fileExt,
		ChunkNext:   1,
	})

	if err != nil || fileChunkInfo.FileMd5 != fileMd5 {
		c.JSON(http.StatusOK, gin.H{
			"ChunkNext": -1,
			"code":      "failed",
			"msg":       "文件校验失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ChunkNext": fileChunkInfo.ChunkNext,
		"code":      "ok",
	})
	return
}

func RemoveFileChunk(c *gin.Context) {

}

func BackPointProcess(c *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			c.Error(err)
		}
	}()
	fileName := c.PostForm("fileName")
	fileExt := path.Ext(fileName)
	fileName = strings.TrimSuffix(fileName, fileExt)
	chunkNum, _ := strconv.Atoi(c.PostForm("chunkNum"))
	chunkMd5 := c.PostForm("chunkMd5")

	//开启事务
	tx := mysql.MysqlDb.Begin()
	defer func() {
		tx.Commit()
	}()
	//校验chunk info
	var fileChunkInfo models.FileChunkInFo
	err = tx.Select("*").Where("fileName=?", fileName).First(&fileChunkInfo).Error
	if err != nil {
		c.JSON(http.StatusOK, "server busy")
		return
	}

	if fileChunkInfo.ChunkNext != chunkNum || fileChunkInfo.Finished == true {
		c.JSON(http.StatusOK, gin.H{
			"code":      "ok",
			"chunkNext": fileChunkInfo.ChunkNext,
		})
		return
	}

	fChunk, err := c.FormFile("chunk")
	if err != nil {
		c.JSON(http.StatusOK, "can not found chunk from form")
		return
	}
	var data []byte
	var file multipart.File

	if file, err = fChunk.Open(); err != nil {
		c.JSON(http.StatusOK, "can not found chunk from form")
		return
	}
	if data, err = io.ReadAll(file); err != nil {
		c.JSON(http.StatusOK, "read chunk failed")
		return
	}
	md5, _ := hash.Md5(string(data))
	if chunkMd5 != md5 {
		c.JSON(http.StatusOK, gin.H{
			"msg":       "md5 verification failed",
			"code":      "failed",
			"chunkNext": chunkNum + 1,
		})
		c.JSON(http.StatusOK, "md5 verification failed")
		return
	}
	file.Close()

	if err = logic.MakeFileChunk(data, fileName, chunkNum); err != nil {
		c.JSON(http.StatusOK, "md5 verification failed")
		return
	}

	if fileChunkInfo.ChunkTotal == chunkNum+1 {
		err = tx.Table("BackPointInfo").Where("fileName=?", fileName).Updates(map[string]interface{}{
			"chunkNext": fileChunkInfo.ChunkNext + 1,
			"finished":  true,
		}).Error
	} else {
		err = tx.Table("BackPointInfo").Where("fileName=?", fileName).Update("chunkNext", fileChunkInfo.ChunkNext+1).Error
	}
	if err = logic.MakeFileChunk(data, fileName, chunkNum); err != nil {
		c.JSON(http.StatusOK, "server busy")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":      "ok",
		"chunkNext": chunkNum + 1,
	})
	return
}

func MergeFileChunk(c *gin.Context) {

}

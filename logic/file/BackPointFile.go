package logic

import (
	"MFile/config"
	"MFile/generate/hash"
	"fmt"
	"io"
	"os"
	"path"
	"strconv"
	"strings"
)

//文件地址 BackPointDir/fileName_chunkNUm
func MakeFileChunk(data []byte, fileName string, chunkNum int) (err error) {
	var strBuilder strings.Builder
	strBuilder.WriteString(fileName)
	strBuilder.WriteString("_")
	strBuilder.WriteString(strconv.Itoa(chunkNum))
	var newFile *os.File
	if newFile, err = os.Create(path.Join(config.BackPointDir, strBuilder.String())); err == nil {
		defer newFile.Close()
		_, err = newFile.Write(data)
		if err != nil {
			newFile.Sync()
		}
	}
	return
}

func MergeFileChunk(fileName, fileNameExt, fileMd5 string, chunkTotal int) (err error) {
	file, err := os.Create(path.Join(config.DirName, fileName+fileNameExt))
	if err != nil {
		return
	}
	defer file.Close()
	var chunk *os.File
	for i := 1; i <= chunkTotal; i++ {
		chunkPath := path.Join(config.BackPointDir, fileName+"_"+strconv.Itoa(i))
		chunk, err = os.Open(chunkPath)
		if err == nil {
			_, err = io.Copy(file, chunk)
			if err != nil {
				break
			}
		}
		chunk.Close()
		if err != nil {
			break
		}
		err = os.Remove(chunkPath)
	}
	file.Close()
	file, err = os.Open(path.Join(config.DirName, fileName+fileNameExt))
	data, err := io.ReadAll(file)
	if err == nil && hash.Md5Byte(data) != fileMd5 {
		return fmt.Errorf("文件MD5校验失败")
	}
	return
}

func RemoveFileChunk(fileName string, chunkTotal int) (err error) {
	for i := 1; i < chunkTotal; i++ {
		err = os.Remove(path.Join(config.BackPointDir, fileName+"_"+strconv.Itoa(i)))
		if err != nil {
			return
		}
	}
	return
}

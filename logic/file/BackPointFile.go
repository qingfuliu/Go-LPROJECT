package logic

import (
	"MFile/config"
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

func MergeFileChunk(fileName, fileMd5 string, chunkTotal int) error {

	return nil
}

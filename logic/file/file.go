package logic

import (
	"MFile/config"
	"MFile/generate/hash"
	"io"
	"mime/multipart"
	"os"
	"path"
	"strings"
	"time"
)

func UpLoadFile(f *multipart.Form) (errs []error) {
	now := time.Now().Format("20060102181209")
	for _, value := range f.File {
		for _, val := range value {

			iFile, err := val.Open()
			if err != nil {
				errs = append(errs, err)
				continue
			}

			ext := path.Ext(val.Filename)
			fileName, _ := hash.Md5(strings.TrimSuffix(val.Filename, ext))
			fileName += "_" + now + ext
			oFile, err := os.Create(path.Join(config.DirName, fileName))
			if err != nil {
				errs = append(errs, err)
				continue
			}

			_, err = io.Copy(oFile, iFile)
			if err != nil {
				errs = append(errs, err)
				continue
			}

			iFile.Close()
			oFile.Close()
		}
	}
	return
}

func DeleteFile(fileName string) error {
	filePath := path.Join(config.DirName, fileName)
	if err := os.Remove(filePath); err != nil {
		return err
	}
	return nil
}

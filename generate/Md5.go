package generate

import (
	md52 "crypto/md5"
	"fmt"
	"io"
)

func Md5(str string) (string, error) {
	md5 := md52.New()
	_, err := io.WriteString(md5, str)
	if err != nil {
		//log
		return "", err
	}
	return fmt.Sprintf("%X", md5.Sum(nil)), nil
}

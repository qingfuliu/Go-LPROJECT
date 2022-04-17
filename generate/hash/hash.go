package hash

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"github.com/spaolacci/murmur3"
	"hash/crc32"
	"hash/crc64"
	"io"
)

func Md5(str string) (string, error) {
	md5 := md5.New()
	_, err := io.WriteString(md5, str)
	if err != nil {
		//log
		return "", err
	}
	return fmt.Sprintf("%X", md5.Sum(nil)), nil
}

func Sha1(str string) (string, error) {
	sha1 := sha1.New()
	_,err:=io.WriteString(sha1,str)
	if err!=nil{
		return "",fmt.Errorf("sha1 error")
	}
	return fmt.Sprintf("%X",sha1.Sum(nil)),nil
}


func Sha256(str string)(string,error){
	sha256:=sha256.New()
	_,err:=io.WriteString(sha256,str)
	if err!=nil{
		return "",err
	}
	return fmt.Sprintf("%X",sha256.Sum(nil)),nil
}

func Sha512(str string)(string,error){
	sha512:=sha512.New()
	_,err:=io.WriteString(sha512,str)
	if err!=nil{
		return "",nil
	}
	return fmt.Sprintf("%X",sha512.Sum(nil)),nil
}

func CRC32(str string)(uint32,error){
		return crc32.ChecksumIEEE([]byte(str)),nil
}

func CRC64(str string)(uint64,error){
	crc64:=crc64.New(crc64.MakeTable(crc64.ISO))
	_,err:=io.WriteString(crc64,str)
	if err!=nil{
		return 0,err
	}
	return crc64.Sum64(),nil
}

func MurMur3(str string)uint64{
	return murmur3.Sum64([]byte(str))
}
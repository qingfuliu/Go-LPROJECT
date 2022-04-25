package authentication

import (
	"MFile/models"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

const (
	key              = "abcd"
	assessTokenLife  = time.Minute * 10
	refreshTokenLife = time.Hour
)

func NewJWTToken(userName, userId string) (string, string, error) {
	now := time.Now()
	jwtModels := &models.JWTModel{
		UserName: userName,
		UserId:   userId,
		TokenType: "assess",
		StandardClaims: jwt.StandardClaims{
			Issuer:    "lqfServer",
			IssuedAt:  now.Unix(),
			ExpiresAt: now.Add(assessTokenLife).Unix(),
		},
	}
	var err error
	var assessToken,refreshToken string
	assessToken,err=jwt.NewWithClaims(jwt.SigningMethodES256,jwtModels).SignedString(key)
	if err!=nil{
		return "","",err
	}
	jwtModels.StandardClaims.ExpiresAt=now.Add(refreshTokenLife).Unix()
	jwtModels.TokenType="refresh"
	refreshToken,err=jwt.NewWithClaims(jwt.SigningMethodES256,jwtModels).SignedString(key)
	if err!=nil{
		return "","",nil
	}
	return assessToken,refreshToken,nil
}

func ParseToken(tokenString string)(model *models.JWTModel,err error){
	token,err:=jwt.ParseWithClaims(tokenString,&models.JWTModel{},func(t *jwt.Token)(interface{},error){
			if _,ok:=t.Method.(*jwt.SigningMethodHMAC);!ok{
				return nil,fmt.Errorf("valid error")
			}
			return key,nil
	})
	if err!=nil{
		return
	}
	model,ok:=token.Claims.(*models.JWTModel)
	if !ok || !token.Valid{
		return nil,fmt.Errorf("parse error")
	}
	return
}


func RefreshToken(refreshToken string)(string, error){
	refresh,err:=ParseToken(refreshToken)
	if err!=nil||refresh.TokenType!="refresh"{
		return "",err
	}

	refresh.TokenType="assess"
	now:=time.Now()
	refresh.ExpiresAt=now.Add(assessTokenLife).Unix()
	refresh.IssuedAt=now.Unix()

	var assessToken string
	assessToken,err=jwt.NewWithClaims(jwt.SigningMethodES256,refresh).SignedString(key)
	if err!=nil{
		return "",err
	}
	return assessToken,nil
}
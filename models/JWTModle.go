package models

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type JWTModel struct {
	UserName string `json:"userName"`
	UserId   string `json:"userId"`
	TokenType string `json:"tokenType"`
	jwt.StandardClaims
}

func (j *JWTModel) Valid() error {
	if len(j.UserId) == 0 || len(j.UserName) == 0 {
		return errors.New("incomplete information")
	}
	if time.Now().Unix() > j.ExpiresAt {
		return errors.New("token expired")
	}
	return nil
}

package common

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

// 定义一个jwt加密的密钥
var jwtKey = []byte("a_secret_crect")

type Claims struct {
	UserId uint
	jwt.StandardClaims
}

func ReleaseToken(key uint) (string, error) {
	expirationTime := time.Now().Add(7 * 24 * time.Hour)
	claims := &Claims{
		UserId: key,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(), //token的有效期
			IssuedAt:  time.Now().Unix(),     //token发放的时间
			Issuer:    "2021080903002",       //发放token的人
			Subject:   "user token",          //主题
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseToken(tokenString string) (*jwt.Token, *Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (i interface{}, err error) {
		return jwtKey, nil
	})

	return token, claims, err
}

package jwt

import (
	"errors"
	"github.com/aiyouyo/bluebell/settingts"
	"github.com/dgrijalva/jwt-go"
	"time"
)

// TokenExpireDuration token过期时间
//const TokenExpireDuration = time.Hour * 24 * 365

var secret = []byte("aiyoyo")

var (
	ErrorInvalidToken = errors.New("无效的Token")
	//ErrorWithoutToken = errors.New("未携带Token")
)

type BlueBellClaims struct {
	UserName string `json:"userName"`
	UserID   int64  `json:"user_id"`
	jwt.StandardClaims
}

func GenToken(username string, userid int64) (token string, err error) {

	tokenExpireDuration := 24 * time.Hour * time.Duration(settingts.Conf.TokenConfig.Duration)
	bbc := BlueBellClaims{
		UserName: username,
		UserID:   userid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenExpireDuration).Unix(), // 过期时间
			Issuer:    "bluebell.com",                             // 签发人
		},
	}
	// 使用指定的签名方法创建签名对象
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, bbc)
	token, err = tokenClaims.SignedString(secret)
	return
}

// ParseToken 解析JWT
func ParseToken(token string) (*BlueBellClaims, error) {

	bbc := new(BlueBellClaims) // 存放解析出来的数据

	tokenClaims, err := jwt.ParseWithClaims(token, bbc, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		return nil, err
	}

	if tokenClaims.Valid {
		return bbc, nil
	}
	//tokenClaims.Claims

	return nil, ErrorInvalidToken
}

func ValidToken(token string) {

}

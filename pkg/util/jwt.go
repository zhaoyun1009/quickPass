package util

import (
	"QuickPass/pkg/setting"
	"QuickPass/pkg/token_cache"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtSecret []byte

type Claims struct {
	Agency   string `json:"agency"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Role     int8   `json:"role"`
	jwt.StandardClaims
}

// GenerateToken generate tokens used for auth
func GenerateToken(agency, username, nickname string, Role int8) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(time.Duration(setting.AppSetting.TokenExpireTime) * time.Second) // 1一个小时的过期时间

	claims := Claims{
		agency,
		username,
		nickname,
		Role,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "quick-pass",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)

	err = token_cache.SetCacheToken(agency, username, token)
	if err != nil {
		log.Println("token_cache.SetCacheToken:", err)
	}
	return token, err
}

// ParseToken parsing token
func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}

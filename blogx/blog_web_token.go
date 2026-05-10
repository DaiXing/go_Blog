package blogx

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// 秘钥。
var jwtKeys = []byte("Jlw987hTQsfwe")

// jwt 生成token
func jwtBuildToken(userId uint, username string) string {
	// 最重要的是，过期时间，userId
	tokenInfo := PxJwtToken{
		UserId:   userId,   // 业务参数
		Username: username, // 业务参数
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 3)), // 过期时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		}}

	tmp := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenInfo)
	token, err := tmp.SignedString(jwtKeys) // 加密。
	CheckErr("SignedString", err)
	Logger.Info("生成jwt ", "token", token)
	return token
}

// jwt 验证token
func jwtVerifyToken(token string) *PxJwtToken {
	if len(token) == 0 {
		panic("token empty")
	}

	// 返回 *Token
	token2, err2 := jwt.ParseWithClaims(
		token,
		&PxJwtToken{},
		func(t *jwt.Token) (any, error) {
			return jwtKeys, nil
		})
	CheckErr("ParseWithClaims", err2)
	// 转为 TokenInfo
	tokenInfo, ok := token2.Claims.(*PxJwtToken)
	if !ok {
		panic("jwt error ")
	}
	return tokenInfo
}

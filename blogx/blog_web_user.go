package blogx

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// 注册。
func handleRegister(ctx *gin.Context) {
	// 转换参数格式。
	var req PxRegisterReq
	err1 := ctx.ShouldBindJSON(&req)
	CheckErr("ShouldBindJSON", err1)

	// 不能重复。
	var count int64
	err2 := Db.Table("users").Where("username = ?", req.Username).Count(&count).Error
	CheckErr("Count", err2)
	if count > 0 {
		ctx.JSON(http.StatusForbidden, &PxBaseResp{
			Error: "username exists ",
		})
		return
	}
	err3 := Db.Table("users").Where("email = ?", req.Email).Count(&count).Error
	CheckErr("Count", err3)
	if count > 0 {
		ctx.JSON(http.StatusForbidden, &PxBaseResp{
			Error: "email exists ",
		})
		return
	}

	// 插入。
	user := User{
		Username: req.Username,
		Password: PasswordEncode(req.Password),
		Email:    req.Email,
	}
	err4 := Db.Create(&user).Error
	CheckErr("Create", err4)
	ctx.JSON(http.StatusOK, &PxBaseResp{
		Desc: "OK",
	})
}

// 登录。
func handleLogin(ctx *gin.Context) {
	// 转换参数格式。
	var req PxLoginReq
	err1 := ctx.ShouldBindJSON(&req)
	CheckErr("ShouldBindJSON", err1)

	if len(req.Username) == 0 {
		panic("username is empty")
	}
	if len(req.Password) == 0 {
		panic("password is empty")
	}

	// 查用户。
	userx := User{
		Username: req.Username,
		Password: PasswordEncode(req.Password),
	}
	var user2 User
	err3 := Db.Where(&userx).Find(&user2).Error
	CheckErr("Find", err3)
	// 不存在。
	if user2.ID == 0 {
		ctx.JSON(http.StatusNotFound, &PxBaseResp{
			Error: "user not found ",
		})
		return
	}

	// jwt token
	token := jwtBuildToken(user2.ID, user2.Username)
	resp := PxLoginResp{
		Token: token,
		Time:  time.Now().Format(time.DateTime),
	}
	ctx.JSON(http.StatusOK, resp)
}

// AOP 验证token 。
func aopVerifyToken(ctx *gin.Context) {
	// 解析token
	tokenx := ctx.GetHeader(Key_usertoken)
	pxToken := jwtVerifyToken(tokenx)

	// 查DB
	var userx User
	err := Db.Where("id = ?", pxToken.UserId).First(&userx).Error
	CheckErr("查user", err)
	if userx.ID == 0 {
		panic(fmt.Sprintf("user not found : %d ", pxToken.UserId))
	}

	// 保存。
	ctx.Set(Key_userid, pxToken.UserId)
	ctx.Set(Key_username, pxToken.Username)
	ctx.Set(Key_user, &userx)
}

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
		panic("jwt error : ")
	}
	return tokenInfo
}

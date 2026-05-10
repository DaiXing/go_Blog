package blogx

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// 注册。
func handleRegister(ctx *gin.Context) {
	// 转换参数格式。
	var req PxRegisterReq
	err1 := ctx.ShouldBindJSON(&req)
	// CheckErr("ShouldBindJSON", err1)
	if err1 != nil {
		WebBadRequest(ctx, err1.Error())
		return
	}

	// 不能重复。
	var count int64
	err2 := Db.Table("users").Where("username = ?", req.Username).Count(&count).Error
	CheckErr("Count", err2)
	if count > 0 {
		WebForbidden(ctx, "username exists ")
		return
	}
	err3 := Db.Table("users").Where("email = ?", req.Email).Count(&count).Error
	CheckErr("Count", err3)
	if count > 0 {
		WebForbidden(ctx, "email exists ")
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
	// CheckErr("ShouldBindJSON", err1)
	if err1 != nil {
		WebBadRequest(ctx, err1.Error())
		return
	}

	if len(req.Username) == 0 {
		WebBadRequest(ctx, "username is empty")
		return
	}
	if len(req.Password) == 0 {
		WebBadRequest(ctx, "password is empty")
		return
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
		WebNotFound(ctx, "user not found ")
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

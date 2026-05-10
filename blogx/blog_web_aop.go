package blogx

import (
	"time"

	"github.com/gin-gonic/gin"
)

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
		// panic(fmt.Sprintf("user not found : %d ", pxToken.UserId))
		WebNotFound(ctx, "user not found ")
		return
	}

	// 保存。
	ctx.Set(Key_userid, pxToken.UserId)
	ctx.Set(Key_username, pxToken.Username)
	ctx.Set(Key_user, &userx)
}

// 记录日志 。
func aopLogReq(ctx *gin.Context) {
	startTime := time.Now()
	uri := ctx.Request.URL.Path
	method := ctx.Request.Method

	ctx.Next()

	costMillis := time.Since(startTime).Milliseconds()
	Logger.Info("http请求", "method", method, "uri", uri, "costMillis", costMillis)
}

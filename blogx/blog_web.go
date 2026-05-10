package blogx

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	// "github.com/golang-jwt/jwt"
)

// server。
var webServer *gin.Engine

const Key_user = "user"
const Key_userid = "user_id"
const Key_username = "user_name"
const Key_usertoken = "user_token"
const Key_postid = "post_id"

// 初始化。
func WebInit() {
	webServer = gin.Default()

	webSetUrl()

	// 运行。默认 8080
	port := ConfigParams.Init.ServerPort
	addr := ":" + strconv.FormatInt(int64(port), 10)
	Logger.Info("web 初始化完成 ", "监听", addr)

	err := webServer.Run(addr)
	CheckErr("webServer 启动", err)
}

// 配置URL
func webSetUrl() {
	// 记录日志。
	webServer.Use(aopLogReq)
	// 包装异常。
	webServer.Use(gin.CustomRecovery(func(ctx *gin.Context, err any) {
		errMsg := "Exception "
		if err != nil {
			errMsg = fmt.Sprintf("%s %v ", errMsg, err) // 加上异常
		}
		// 必须调用 abort 阻断流程。
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, &PxBaseResp{
			Error: errMsg,
		})
	}))

	// 健康检查。
	webServer.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, &PxBaseResp{
			Desc: "blog web server . " + time.Now().Format(time.DateTime),
		})
	})

	// 注册。
	webServer.POST("/register", handleRegister)

	// 登录。 返回token
	webServer.POST("/login", handleLogin)

	// 博客。 需要token
	group1 := webServer.Group("/post")
	group1.Use(aopVerifyToken)               // 验证token。
	group1.POST("/add", handlePostAdd)       // 新建。
	group1.POST("/update", handlePostUpdate) // 更新。
	group1.POST("/delete", handlePostDelete) // 删除。

	// 评论。 需要token
	group2 := webServer.Group("/comment")
	group2.Use(aopVerifyToken)            // 验证token。
	group2.POST("/add", handleCommentAdd) // 新建。

	// 查询。不需要token。
	webServer.POST("/post/query-list", handlePostQueryList)
	webServer.GET("/post/query-one", handlePostQueryOne)
}

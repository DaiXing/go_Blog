package blogx

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 评论。新建。
func handleCommentAdd(ctx *gin.Context) {
	userid := ctx.GetUint(Key_userid)

	// 转换参数格式。
	var req PxCommentAddReq
	err1 := ctx.ShouldBindJSON(&req)
	CheckErr("ShouldBindJSON", err1)

	// 查文章。
	var post1 Post
	err2 := Db.First(&post1, req.PostId).Error
	CheckErr("First", err2)

	// 插入。
	comment1 := Comment{
		UserId:  userid,
		PostId:  post1.ID,
		Content: req.Content,
	}
	err3 := Db.Create(&comment1).Error
	CheckErr("Create", err3)

	ctx.JSON(http.StatusOK, &PxCommentAddResp{
		PxBaseResp: PxBaseResp{
			Desc: "OK",
		},
		CommentAdded: comment1,
	})
}

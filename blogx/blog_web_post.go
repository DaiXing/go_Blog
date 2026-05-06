package blogx

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 文章。新建。
func handlePostAdd(ctx *gin.Context) {
	userid := ctx.GetUint(Key_userid)

	// 转换参数格式。
	var req PxPostAddReq
	err1 := ctx.ShouldBindJSON(&req)
	CheckErr("ShouldBindJSON", err1)

	post1 := Post{
		UserId:  userid,
		Title:   req.Title,
		Content: req.Content,
	}
	err2 := Db.Create(&post1).Error
	CheckErr("Create", err2)

	ctx.JSON(http.StatusOK, &PxPostAddResp{
		PxBaseResp: PxBaseResp{
			Desc: "OK handlePostAdd",
		},
		PostAdded: &post1,
	})
}

// 文章。修改。
func handlePostUpdate(ctx *gin.Context) {
	userid := ctx.GetUint(Key_userid)

	// 转换参数格式。
	var req PxPostUpdateReq
	err1 := ctx.ShouldBindJSON(&req)
	CheckErr("ShouldBindJSON", err1)

	// 先查询。 使用ID
	var post1 Post
	err2 := Db.First(&post1, req.PostId).Error
	CheckErr("First", err2)

	if post1.UserId != userid {
		ctx.JSON(http.StatusNotFound, &PxBaseResp{
			Error: "post user not match",
		})
		return
	}

	// 用ID更新。
	// post2 := Post{
	// 	Title:   req.Title,
	// 	Content: req.Content,
	// }
	// post2.ID = post1.ID
	// err3 := Db.Updates(post2).Error// 错误。危险。

	// 安全更新。
	updateMap := map[string]any{
		"title":   req.Title,
		"content": req.Content,
	}
	err3 := Db.Model(&Post{}).Where("id = ? ", post1.ID).Updates(updateMap).Error
	CheckErr("Updates", err3)

	// 再查一次。
	var post3 Post
	err4 := Db.First(&post3, post1.ID).Error
	CheckErr("First", err4)

	ctx.JSON(http.StatusOK, &PxPostUpdateResp{
		PxBaseResp: PxBaseResp{
			Desc: "OK handlePostUpdate",
		},
		PostUpdated: &post3,
	})
}

// 文章。删除。
func handlePostDelete(ctx *gin.Context) {
	userid := ctx.GetUint(Key_userid)

	// 转换参数格式。
	var req PxPostDeleteReq
	err1 := ctx.ShouldBindJSON(&req)
	CheckErr("ShouldBindJSON", err1)

	// 先查询。 使用ID
	var post1 Post
	err2 := Db.First(&post1, req.PostId).Error
	CheckErr("First", err2)

	if post1.UserId != userid {
		ctx.JSON(http.StatusNotFound, &PxBaseResp{
			Error: "post user not match",
		})
		return
	}

	// 用ID删除。
	err3 := Db.Delete(&Post{}, post1.ID).Error
	CheckErr("Delete", err3)

	ctx.JSON(http.StatusOK, &PxPostUpdateResp{
		PxBaseResp: PxBaseResp{
			Desc: "OK handlePostDelete",
		},
	})
}

// 文章。查单个。
func handlePostQueryOne(ctx *gin.Context) {
	postId, ok := ctx.GetQuery(Key_postid)
	if !ok {
		panic("param postId empty")
	}

	// 查文章。
	var post1 Post
	// 关联查询。
	err1 := Db.Preload("Comments", func(db *gorm.DB) *gorm.DB {
		return db.Order("id desc").Limit(2) // 只查N条评论。
	}).Preload("User").First(&post1, postId).Error
	CheckErr("First", err1)

	if post1.ID == 0 {
		ctx.JSON(http.StatusNotFound, &PxBaseResp{
			Error: "post not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, &PxPostOneResp{
		Post: post1,
	})
}

// 文章。查列表。
func handlePostQueryList(ctx *gin.Context) {
	var req PxPostQueryListReq
	err1 := ctx.ShouldBindJSON(&req)
	CheckErr("ShouldBindJSON", err1)

	// 事务。
	tx := Db.Model(&Post{}) // 构建查询。
	if req.PostId > 0 {     // 可选字段
		tx = tx.Where("id = ?", req.PostId)
	}
	titlex := strings.TrimSpace(req.Title)
	if len(titlex) > 0 { // 可选字段
		// tx = tx.Where("title like '%?%'", titlex)// 错误。
		tx = tx.Where("title like ?", "%"+titlex+"%") // 正确。
	}
	// 分页
	indexBegin := req.PageSize * (req.PageNo - 1)       // 分页
	tx.Offset(int(indexBegin)).Limit(int(req.PageSize)) // 分页。

	// 查列表。
	var posts []Post
	err2 := tx.Find(&posts).Error
	CheckErr("Find", err2)

	ctx.JSON(http.StatusOK, &PxPostListResp{
		Posts: posts,
	})
}

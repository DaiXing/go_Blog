package testx

import (
	"fmt"
	"strconv"
	"testing"

	"gitee.com/dx/go_blog/blogx"
)

var userToken string

func Test1(tt *testing.T) {
	// 正常注册。
	var req1 blogx.PxRegisterReq
	req1.Username = "Apple"
	req1.Email = "apple@cc.com"
	req1.Password = "98j76sw"

	fmt.Println("测试： 正常注册")
	tt.Run("正常注册", func(t *testing.T) {
		_, resp1 := blogx.PostJson[blogx.PxBaseResp](url_register, req1, nil)
		if len(resp1.Error) > 0 {
			tt.Fatal("注册错误")
		}
	})

	fmt.Println("测试： 不能重复注册")
	tt.Run("不能重复注册", func(t *testing.T) {
		_, resp2 := blogx.PostJson[blogx.PxBaseResp](url_register, req1, nil)
		if len(resp2.Error) == 0 {
			tt.Fatal("重复注册。异常。")
		}
	})

	fmt.Println("测试： 正常登录")
	tt.Run("正常登录", func(t *testing.T) {
		var req2 blogx.PxLoginReq
		req2.Username = req1.Username
		req2.Password = req1.Password
		_, resp2 := blogx.PostJson[blogx.PxLoginResp](url_login, req2, nil)
		if len(resp2.Error) > 0 {
			tt.Fatal("登录错误")
		}
		userToken = resp2.Token
	})

	fmt.Println("测试： 错误登录")
	tt.Run("错误登录", func(t *testing.T) {
		var req2 blogx.PxLoginReq
		req2.Username = "不存在的name"
		req2.Password = "不存在的密码"
		_, resp2 := blogx.PostJson[blogx.PxLoginResp](url_login, req2, nil)
		if len(resp2.Error) == 0 {
			tt.Fatal("登录错误")
		}
	})
	//==================
	// 发表文章
	var req3 blogx.PxPostAddReq
	req3.Title = "文章1"
	req3.Content = "文章1的内容啊"

	// 不带token。
	fmt.Println("测试： 发表文章，不带token。")
	_, resp3 := blogx.PostJson[blogx.PxPostAddResp](url_post_add, req3, nil)
	if len(resp3.Error) == 0 {
		tt.Fatal("发表文章，不带token。")
	}

	// 带上错误的token
	fmt.Println("测试： 发表文章，带上错误的token。")
	_, resp4 := blogx.PostJsonWithToken[blogx.PxPostAddResp](url_post_add, req3, "tokenxxxx")
	if len(resp4.Error) == 0 {
		tt.Fatal("发表文章，带上错误的token。")
	}
	//==================

	// 带上正确的token
	fmt.Println("测试： 发表文章，带上正确的token。")
	_, resp5 := blogx.PostJsonWithToken[blogx.PxPostAddResp](url_post_add, req3, userToken)
	if len(resp5.Error) > 0 {
		tt.Fatal("发表文章，带上正确的token。")
	}
	if resp5.PostAdded.ID == 0 {
		tt.Fatal("创建文章，没有ID")
	}

	// 修改文章。 不能修改别人的文章。
	fmt.Println("测试： 修改文章。 不能修改别人的文章。")
	req6 := blogx.PxPostUpdateReq{
		PostId:  1,
		Title:   "文章：修改了标题 ",
		Content: "文章：修改了内容",
	}
	_, resp6 := blogx.PostJsonWithToken[blogx.PxPostUpdateResp](url_post_update, &req6, userToken)
	if len(resp6.Error) == 0 {
		tt.Fatal("修改文章。 不能修改别人的文章。")
	}

	// 修改文章。 可以修改自己的文章。
	fmt.Println("测试： 修改文章。 可以修改自己的文章。")
	req7 := req6
	req7.PostId = resp5.PostAdded.ID // 自己的文章。
	_, resp7 := blogx.PostJsonWithToken[blogx.PxPostUpdateResp](url_post_update, &req7, userToken)
	if len(resp7.Error) > 0 {
		tt.Fatal("修改文章。 可以修改自己的文章。")
	}
	if req7.Title != resp7.PostUpdated.Title {
		tt.Fatal("修改文章。 标题没有被修改")
	}
	if req7.Title != resp7.PostUpdated.Title {
		tt.Fatal("修改文章。 内容没有被修改")
	}
	//==================
	// 删除文章。正常删除。
	fmt.Println("测试： 删除文章。正常删除。")
	var req8 blogx.PxPostDeleteReq
	req8.PostId = req7.PostId
	req8.Reason = "测试删除"
	_, resp8 := blogx.PostJsonWithToken[blogx.PxBaseResp](url_post_delete, req8, userToken)
	if len(resp8.Error) > 0 {
		tt.Fatal("删除文章。正常删除。")
	}

	// 删除后，再查一次。
	resp8x := blogx.GetJson[blogx.PxPostOneResp](
		url_post_query_one + "?post_id=" + strconv.FormatUint(uint64(req7.PostId), 10))
	if len(resp8x.Error) == 0 {
		tt.Fatal("删除后，还能查到")
	}

	// 删除文章。删除不存在的文章。
	fmt.Println("测试： 删除文章。删除不存在的文章。")
	var req9 blogx.PxPostDeleteReq
	req9.PostId = 800009
	_, resp9 := blogx.PostJsonWithToken[blogx.PxBaseResp](url_post_delete, req9, userToken)
	if len(resp9.Error) == 0 {
		tt.Fatal("删除文章。删除不存在的文章。")
	}

	//==================

	// 先查出一个文章
	fmt.Println("测试： 先查出一个文章")
	req10 := blogx.PxPostQueryListReq{
		PageNo:   1,
		PageSize: 1,
	}
	_, resp10 := blogx.PostJson[blogx.PxPostListResp](url_post_query_list, &req10, nil)
	post10 := resp10.Posts[0]

	// 评论。 新增评论。
	fmt.Println("测试： 评论。 新增评论。")
	req11 := blogx.PxCommentAddReq{
		PostId:  post10.ID,
		Content: "评论：新来的一个评论",
	}
	_, resp11 := blogx.PostJsonWithToken[blogx.PxCommentAddResp](url_comment_add, req11, userToken)
	if len(resp11.Error) > 0 {
		tt.Fatal("评论。 新增评论。")
	}
	comment11 := resp11.CommentAdded

	// 查看评论。在第一个位置。
	fmt.Println("测试： 评论。 查看评论。在第一个位置。")
	resp12 := blogx.GetJson[blogx.PxPostOneResp](
		url_post_query_one + "?post_id=" + strconv.FormatUint(uint64(comment11.PostId), 10))
	if len(resp12.Error) > 0 {
		tt.Fatal("查看评论。在第一个位置。")
	}
	comment12 := resp12.Post.Comments[0]
	fmt.Println("comment11 =", blogx.ToJsonString(comment11))
	fmt.Println("comment12 =", blogx.ToJsonString(comment12))
	if comment11.Content != comment12.Content {
		tt.Fatal("查看评论。文本不一致。")
	}

}

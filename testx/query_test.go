package testx

import (
	"fmt"
	"strconv"
	"testing"

	"gitee.com/dx/go_blog/blogx"
)

func TestQuery(t *testing.T) {

	// 查单个文章。 id不存在。
	fmt.Println("测试： 查单个文章。 id不存在。")
	resp1 := blogx.GetJson[blogx.PxPostOneResp](url_post_query_one + "?post_id=333")
	if len(resp1.Error) == 0 {
		t.Fatal("查单个文章，错误")
	}

	// 查全部文章。
	fmt.Println("测试： 查全部文章。")
	var req2 blogx.PxPostQueryListReq
	req2.PostId = 0
	req2.Title = ""
	req2.PageNo = 1
	req2.PageSize = 1000
	_, resp2 := blogx.PostJson[blogx.PxPostListResp](url_post_query_list, req2, nil)
	if len(resp2.Posts) == 0 {
		t.Fatal("查询文章列表错误")
	}
	post0 := resp2.Posts[0]

	// 再查单个文章。
	fmt.Println("测试： 再查单个文章。")
	resp3 := blogx.GetJson[blogx.PxPostOneResp](
		url_post_query_one + "?post_id=" + strconv.FormatInt(int64(post0.ID), 10),
	)
	if resp3.Post.ID == 0 {
		t.Fatal("查单个文章，错误")
	}
	if resp3.Post.User.ID == 0 {
		t.Fatal("查单个文章，但是没有user")
	}
	if len(resp3.Post.Comments) == 0 {
		t.Fatal("查单个文章，但是没有评论")
	}
}

package blogx

import "github.com/golang-jwt/jwt/v5"

// 实体类。
type PxRegisterReq struct {
	Username string
	Password string
	Email    string
}
type PxLoginReq struct {
	Username string
	Password string
}
type PxJwtToken struct {
	jwt.RegisteredClaims // 继承。
	UserId               uint
	Username             string
}
type PxBaseResp struct {
	Error string
	Desc  string
}
type PxLoginResp struct {
	PxBaseResp // 继承。
	Token      string
	Time       string
}
type PxPostAddReq struct {
	Title   string
	Content string
}
type PxPostAddResp struct {
	PxBaseResp // 继承。
	PostAdded  *Post
}
type PxPostUpdateReq struct {
	PostId  uint
	Title   string
	Content string
}
type PxPostUpdateResp struct {
	PxBaseResp  // 继承。
	PostUpdated *Post
}
type PxPostDeleteReq struct {
	PostId uint
	Reason string
}
type PxCommentAddReq struct {
	PostId  uint
	Content string
}
type PxCommentAddResp struct {
	PxBaseResp   // 继承。
	CommentAdded Comment
}
type PxPostOneResp struct {
	PxBaseResp      // 继承。
	Post       Post // 文章。
}
type PxPostQueryListReq struct { // 查列表。
	PostId   uint
	Title    string
	PageNo   uint
	PageSize uint
}
type PxPostListResp struct { // 查列表。
	PxBaseResp        // 继承。
	Posts      []Post // 文章。
}

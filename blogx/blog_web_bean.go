package blogx

import "github.com/golang-jwt/jwt/v5"

// 实体类。
type PxRegisterReq struct {
	Username string `json:"username" binding:"required,min=5,max=20"`
	Password string `json:"password" binding:"required,min=8,max=20"`
	Email    string `json:"email" binding:"required,email"`
}
type PxLoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type PxJwtToken struct {
	jwt.RegisteredClaims        // 继承。
	UserId               uint   `json:"userId"`
	Username             string `json:"username"`
}
type PxBaseResp struct {
	Error string `json:"error"`
	Desc  string `json:"desc"`
}
type PxLoginResp struct {
	PxBaseResp        // 继承。
	Token      string `json:"token"`
	Time       string `json:"time"`
}
type PxPostAddReq struct {
	Title   string `json:"title" binding:"required,min=3"`
	Content string `json:"content" binding:"required,min=5"`
}
type PxPostAddResp struct {
	PxBaseResp       // 继承。
	PostAdded  *Post `json:"postAdded"`
}
type PxPostUpdateReq struct {
	PostId  uint   `json:"postId"`
	Title   string `json:"title" binding:"required,min=3"`
	Content string `json:"content" binding:"required,min=5"`
}
type PxPostUpdateResp struct {
	PxBaseResp        // 继承。
	PostUpdated *Post `json:"postUpdated"`
}
type PxPostDeleteReq struct {
	PostId uint   `json:"postId"`
	Reason string `json:"reason" binding:"required"`
}
type PxCommentAddReq struct {
	PostId  uint   `json:"postId"`
	Content string `json:"content" binding:"required,min=5"`
}
type PxCommentAddResp struct {
	PxBaseResp           // 继承。
	CommentAdded Comment `json:"commentAdded"`
}
type PxPostOneResp struct {
	PxBaseResp      // 继承。
	Post       Post `json:"post"` // 文章。
}
type PxPostQueryListReq struct { // 查列表。
	PostId   uint   `json:"postId"`
	Title    string `json:"title"`
	PageNo   uint   `json:"pageNo"`
	PageSize uint   `json:"pageSize"`
}
type PxPostListResp struct { // 查列表。
	PxBaseResp        // 继承。
	Posts      []Post `json:"posts"` // 文章。
}

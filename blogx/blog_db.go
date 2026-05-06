package blogx

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// 数据库连接。
var Db *gorm.DB

// 数据库，初始化。
func DbInit() {
	dbx, err := gorm.Open(sqlite.Open("blog.db"), &gorm.Config{})
	CheckErr("gorm.Open", err)
	Db = dbx
	Logger.Info("DB 初始化完成")

	DbReloadTable()
	DbInsertRows()
	DbQueryRows()
}

// 重建表。
func DbReloadTable() {
	err1 := Db.Migrator().DropTable(&User{}, &Post{}, &Comment{})
	CheckErr("DropTable", err1)

	err2 := Db.Migrator().CreateTable(&User{}, &Post{}, &Comment{})
	CheckErr("CreateTable", err2)

	tableNames, err := Db.Migrator().GetTables()
	CheckErr("db.Migrator().GetTables()", err)
	Logger.Info("DB 重建表完成 ", "tableNames", tableNames)
}

type User struct {
	gorm.Model        // 包含几个常用字段
	Username   string `gorm:"index;size:100;not null;"`
	Password   string `gorm:"size:100;null null;"`
	Email      string `gorm:"index;size:100;not null;"`
	Address    string `gorm:"size:300;"`
}

// 博客文章。
type Post struct {
	gorm.Model        // 包含几个常用字段
	Title      string `gorm:"index;size:100;not null;"`
	Content    string `gorm:"size:5000;"`
	UserId     uint   `gorm:"index;"`
	User       User
	Comments   []Comment // 评论。
}

// 文章评论。
type Comment struct {
	gorm.Model        // 包含几个常用字段
	Content    string `gorm:"size:300;"`
	UserId     uint   `gorm:"index;"`
	User       User
	PostId     uint `gorm:"index;"`
	Post       Post
}

// 插入一些行。
func DbInsertRows() {
	users := []User{
		{
			Username: "Jack",
			Password: PasswordEncode("123123"),
			Email:    "jack@cc.com",
			Address:  "NewYork",
		},
		{
			Username: "Tom",
			Password: PasswordEncode("787878"),
			Email:    "tom@cc.com",
			Address:  "Chicago",
		},
	}
	err1 := Db.Create(users).Error
	CheckErr("db.Create(&users)", err1)

	posts := []Post{
		{
			Title:   "文章：伊朗封锁波斯湾",
			Content: "内容：战争还在继续啊。。。。",
			UserId:  users[1].ID,
		}, {
			Title:   "文章：纽约股市崩盘",
			Content: "内容：股市崩了3000点。。。。",
			UserId:  users[1].ID,
		},
	}
	err2 := Db.Create(posts).Error
	CheckErr("db.Create(posts)", err2)

	comments := []Comment{
		{
			Content: "评论：伊朗太强了！",
			UserId:  posts[0].UserId,
			PostId:  posts[0].ID,
		},
		{
			Content: "评论：伊朗估计撑不住了！",
			UserId:  posts[0].UserId,
			PostId:  posts[0].ID,
		},
		{
			Content: "评论：美国又要股灾了！",
			UserId:  posts[1].UserId,
			PostId:  posts[1].ID,
		},
	}
	err3 := Db.Create(comments).Error
	CheckErr("db.Create(comments)", err3)
	Logger.Info("DB insert rows done ")
}

// 查询全部行
func DbQueryRows() {
	var users []User
	Db.Find(&users)
	var posts []Post
	Db.Find(&posts)
	var comments []Comment
	Db.Find(&comments)
	Logger.Info("DB 查询 ", "userList", ToJsonString(users))
	Logger.Info("DB 查询 ", "postList", ToJsonString(posts))
	Logger.Info("DB 查询 ", "commentList", ToJsonString(comments))
}

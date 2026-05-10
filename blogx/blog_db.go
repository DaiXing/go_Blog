package blogx

import (
	"time"

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

	// 先插入一些数据。
	if ConfigParams.Init.DbInsertRowsEnable {
		DbInsertRows()
	}

	DbQueryRows()
}

// 重建表。
func DbReloadTable() {
	// 先删除表。
	if ConfigParams.Init.DbDropTableEnable {
		err1 := Db.Migrator().DropTable(&User{}, &Post{}, &Comment{})
		CheckErr("DropTable", err1)
	}

	// 建表
	err2 := Db.Migrator().CreateTable(&User{}, &Post{}, &Comment{})
	CheckErr("CreateTable", err2)

	// 查表
	tableNames, err := Db.Migrator().GetTables()
	CheckErr("db.Migrator().GetTables()", err)

	Logger.Info("DB 重建表完成 ", "tableNames", tableNames)
}

// 复制 gorm.Model 。修改json格式等。
type BaseModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
type User struct {
	BaseModel        // 包含几个常用字段
	Username  string `gorm:"index;size:100;not null;" json:"username"`
	Password  string `gorm:"size:100;null null;" json:"-"`
	Email     string `gorm:"uniqueIndex;size:100;not null;" json:"email"`
	Address   string `gorm:"size:300;" json:"address"`
}

// 博客文章。
type Post struct {
	BaseModel           // 包含几个常用字段
	Title     string    `gorm:"index;size:100;not null;" json:"title"`
	Content   string    `gorm:"size:5000;" json:"content"`
	UserId    uint      `gorm:"index;" json:"userId"`
	User      User      `json:"user"`
	Comments  []Comment `json:"comments"` // 评论。
}

// 文章评论。
type Comment struct {
	BaseModel        // 包含几个常用字段
	Content   string `gorm:"size:300;" json:"content"`
	UserId    uint   `gorm:"index;" json:"userId"`
	User      User   `json:"user"`
	PostId    uint   `gorm:"index;" json:"postId"`
	Post      Post   `json:"post"`
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

	Logger.Info("DB 插入行数据，完成 ")
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

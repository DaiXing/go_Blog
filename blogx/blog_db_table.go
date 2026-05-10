package blogx

import (
	"time"

	"gorm.io/gorm"
)

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

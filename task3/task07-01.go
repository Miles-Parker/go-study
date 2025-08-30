package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

//题目：钩子函数
//继续使用博客系统的模型。
//要求 ：
//为 Post 模型添加一个钩子函数，在文章创建时自动更新用户的文章数量统计字段。
//为 Comment 模型添加一个钩子函数，在评论删除时检查文章的评论数量，如果评论数量为 0，则更新文章的评论状态为 "无评论"。

type User struct {
	gorm.Model
	Name      string
	PostCount int `gorm:"default:0"`
	Posts     []Post
}

type Post struct {
	gorm.Model
	Title         string
	Content       string
	CommentStatus string `gorm:"default:'有评论'"`
	UserID        uint
	Comments      []Comment
}

type Comment struct {
	gorm.Model
	Content string
	PostID  uint
}

func (p *Post) AfterCreate(tx *gorm.DB) error {
	return tx.Model(&User{}).Where("id = ?", p.UserID).
		Update("post_count", gorm.Expr("post_count + 1")).Error
}

func (c *Comment) AfterDelete(tx *gorm.DB) error {
	var count int64
	if err := tx.Model(&Comment{}).Where("post_id = ?", c.PostID).Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		return tx.Model(&Post{}).Where("id = ?", c.PostID).
			Update("comment_status", "无评论").Error
	}
	return nil
}

func main() {
	dsn := "root:root@tcp(10.4.8.22:3306)/gorm?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&User{}, &Post{}, &Comment{})
}

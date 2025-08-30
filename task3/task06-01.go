package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

//题目：基于上述博客系统的模型定义。
//要求 ：
//编写Go代码，使用Gorm查询某个用户发布的所有文章及其对应的评论信息。
//编写Go代码，使用Gorm查询评论数量最多的文章信息。

type User struct {
	gorm.Model
	Name  string
	Posts []Post
}

type Post struct {
	gorm.Model
	Title    string
	UserID   uint
	Comments []Comment
}

type Comment struct {
	gorm.Model
	Content string
	PostID  uint
}

func main() {
	dsn := "root:root@tcp(10.4.8.22:3306)/gorm?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// 查询用户1的所有文章及评论
	userID := 1
	var user User
	err = db.Preload("Posts.Comments").First(&user, userID).Error
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("用户%d的文章及评论:\n", userID)
	for _, post := range user.Posts {
		fmt.Printf("- 文章: %s (评论数: %d)\n", post.Title, len(post.Comments))
	}

	// 查询评论最多的文章
	var popularPost Post
	err = db.Model(&Post{}).
		Select("posts.*, COUNT(comments.id) as comment_count").
		Joins("LEFT JOIN comments ON comments.post_id = posts.id").
		Group("posts.id").
		Order("comment_count DESC").
		First(&popularPost).Error
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\n评论最多的文章: %s (评论数: %d)\n",
		popularPost.Title,
		len(popularPost.Comments))
}

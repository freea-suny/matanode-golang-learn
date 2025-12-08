/*
 * @Description:
 * @version: 1.0.0
 * @Author: sun.yong
 * @Date: 2025-12-07 18:38:35
 * @LastEditors: sun.yong
 * @LastEditTime: 2025-12-08 09:17:30
 */
package upgrade

import (
	"fmt"

	"gorm.io/gorm"
)

/**
进阶gorm
题目1：模型定义
假设你要开发一个博客系统，有以下几个实体： User （用户）、 Post （文章）、 Comment （评论）。
要求 ：
使用Gorm定义 User 、 Post 和 Comment 模型，其中 User 与 Post 是一对多关系（一个用户可以发布多篇文章）， Post 与 Comment 也是一对多关系（一篇文章可以有多个评论）。
编写Go代码，使用Gorm创建这些模型对应的数据库表。
*/

type User struct {
	ID      uint `gorm:"primaryKey"`
	Name    string
	Email   string `gorm:"uniqueIndex"`
	PostNum uint   `gorm:"default:0"`
}

type Post struct {
	ID            uint `gorm:"primaryKey"`
	UserID        uint
	Title         string
	Content       string
	CommentNum    int
	CommentStatus string
	Comments      []Comment
}
type Comment struct {
	ID      uint `gorm:"primaryKey"`
	PostID  uint
	UserID  uint
	Content string
}

func CreatTable(db *gorm.DB) {
	fmt.Println("创建数据库表")
	// db.Migrator().DropTable("users")
	db.AutoMigrate(&User{}, &Post{}, &Comment{})
	fmt.Println("创建数据库表成功")

	// db.Create(&User{Name: "张三", Email: "zhangsan@example.com"})
	db.Create(&Post{UserID: 1, Title: "第一篇文章", Content: "这是第一篇文章的内容"})
	db.Create(&Comment{PostID: 1, UserID: 1, Content: "这是第一篇文章的评论"})
	db.Create(&Post{UserID: 1, Title: "第2篇文章", Content: "这是第2篇文章的内容"})
	db.Create(&Comment{PostID: 2, UserID: 1, Content: "这是第2篇文章的评论1"})
	db.Create(&Comment{PostID: 2, UserID: 1, Content: "这是第2篇文章的评论2"})
}

/*
*
题目2：关联查询
基于上述博客系统的模型定义。
要求 ：
编写Go代码，使用Gorm查询某个用户发布的所有文章及其对应的评论信息。
编写Go代码，使用Gorm查询评论数量最多的文章信息。
*/
func QueryPsotAndPreComment(db *gorm.DB, userId uint) ([]Post, error) {
	var posts []Post
	err := db.Model(&Post{}).Where("user_id = ?", userId).Preload("Comments").Find(&posts).Error
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func QueryPostWithMaxComment(db *gorm.DB) (Post, error) {
	var post Post
	db.Model(Post{}).Order("Comment_num desc").First(&post)
	// 检查是否查询到了文章
	if post.ID == 0 {
		return post, fmt.Errorf("没有查询到评论数量最多的文章")
	}
	return post, nil
}

/**
题目3：钩子函数
继续使用博客系统的模型。
要求 ：
为 Post 模型添加一个钩子函数，在文章创建时自动更新用户的文章数量统计字段。
为 Comment 模型添加一个钩子函数，在评论删除时检查文章的评论数量，如果评论数量为 0，则更新文章的评论状态为 "无评论"。
*/

func (p *Post) BeforeCreate(tx *gorm.DB) (err error) {
	// 更新用户的文章数量统计字段
	err = tx.Model(&User{}).Where("id = ?", p.UserID).Update("post_num", gorm.Expr("post_num + ?", 1)).Error
	if err != nil {
		return err
	}
	return nil
}

func (c *Comment) BeforeCreate(tx *gorm.DB) (err error) {
	// 更新文章的评论数量统计字段
	err = tx.Model(&Post{}).Where("id = ?", c.PostID).Update("comment_num", gorm.Expr("comment_num + ?", 1)).Error
	if err != nil {
		return err
	}
	return nil
}

// 为 Comment 模型添加一个钩子函数，在评论删除时检查文章的评论数量，如果评论数量为 0，则更新文章的评论状态为 "无评论"。
func (c *Comment) BeforeDelete(tx *gorm.DB) (err error) {
	// 检查文章是否存在
	var post Post
	err = tx.Model(&Post{}).Where("id = ?", c.PostID).First(&post).Error
	if err != nil {
		return err
	}
	// 更新文章的评论数量统计字段
	tx.Model(&Post{}).Where("id = ?", c.PostID).Update("comment_num", gorm.Expr("comment_num - ?", 1))
	if err != nil {
		return err
	}
	// 如果评论数量为 0，则更新文章的评论状态为 "无评论"
	if post.CommentNum == 0 {
		tx.Model(&Post{}).Where("id = ?", c.PostID).Update("comment_status", "无评论")
	}
	return nil
}

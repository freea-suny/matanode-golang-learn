package models

import "gorm.io/gorm"

/*
 * @Description: 文章模型
 * @version: 1.0.0
 * @Author: sun.yong
 * @Date: 2025-12-10 21:48:22
 * @LastEditors: sun.yong
 * @LastEditTime: 2025-12-10 21:49:32
 */

type Post struct {
	gorm.Model
	Title    string
	Content  string
	UserID   uint
	User     User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Comments []Comment `gorm:"foreignKey:PostID" json:"comments,omitempty"`
}

/*
 * @Description:
 * @version: 1.0.0
 * @Author: sun.yong
 * @Date: 2025-12-10 22:00:13
 * @LastEditors: sun.yong
 * @LastEditTime: 2025-12-10 22:10:29
 */
package models

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	PostID  uint
	UserID  uint
	User    User `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Post    Post `gorm:"foreignKey:PostID" json:"post,omitempty"`
	Content string
}

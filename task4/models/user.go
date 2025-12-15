/*
 * @Description: 用户模型
 * @version: 1.0.0
 * @Author: sun.yong
 * @Date: 2025-12-10 20:42:28
 * @LastEditors: sun.yong
 * @LastEditTime: 2025-12-10 21:47:17
 */
package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserName string
	Email    string
	Password string
	//表示文章切片，post中的userId指向user的id,json序列化时，如果为空，则忽略
	Posts    []Post    `gorm:"foreignKey:UserID" json:"posts,omitempty"`
	Comments []Comment `gorm:"foreignKey:UserID" json:"comments,omitempty"`
}

// 保存用户信息之前，先加密用户密码（钩子函数）
func (u *User) BeforeSave(tx *gorm.DB) (err error) {
	//加密用户密码，使用bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// 验证密码
func (u *User) VerifyPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

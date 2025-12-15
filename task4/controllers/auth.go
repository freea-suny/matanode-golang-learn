/*
 * @Description:
 * @version: 1.0.0
 * @Author: sun.yong
 * @Date: 2025-12-11 23:45:36
 * @LastEditors: sun.yong
 * @LastEditTime: 2025-12-15 15:42:30
 */
package controllers

import (
	"net/http"
	"task4/configs"
	"task4/models"
	"task4/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 定义注册请求结构体
type RegisterRequest struct {
	UserName string `json:"user_name" binding:"required,min=3,max=10"`
	Password string `json:"password" binding:"required,min=6,max=12"`
	Email    string `json:"email" binding:"required,email"`
}

// 定义登录请求结构体
type LoginRequest struct {
	UserName string `json:"user_name" binding:"required,min=3,max=10"`
	Password string `json:"password" binding:"required,min=6,max=12"`
}

// 定义注册处理函数(不需要返回值，因为返回的信息会放到gin的响应中)
func Register(c *gin.Context) {
	//绑定请求参数到结构体
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorRes(c, "请求数据异常，注册失败", http.StatusBadRequest)
		return
	}

	// TODO: 实现注册逻辑
	// 1. 检查用户名是否已存在
	var user models.User
	if err := configs.GetDB().First(&user, "user_name = ? ", req.UserName).Error; err == nil {
		utils.ErrorRes(c, "用户名已存在", http.StatusBadRequest)
		return
	} else if err != nil && err != gorm.ErrRecordNotFound {
		utils.ErrorRes(c, "数据库错误", http.StatusInternalServerError)
		return
	}
	//2. 检查邮箱是否存在
	// 检查邮箱是否已存在
	if err := configs.GetDB().First(&user, "email = ? ", req.Email).Error; err == nil {
		utils.ErrorRes(c, "邮箱已存在", http.StatusBadRequest)
		return
	} else if err != nil && err != gorm.ErrRecordNotFound {
		utils.ErrorRes(c, "数据库错误", http.StatusInternalServerError)
		return
	}

	// 3. 存储用户信息到数据库
	user = models.User{
		UserName: req.UserName,
		Password: req.Password,
		Email:    req.Email,
	}
	if err := configs.GetDB().Create(&user).Error; err != nil {
		utils.ErrorRes(c, "数据库错误", http.StatusInternalServerError)
		return
	}
	// 4. 返回成功响应
	utils.SuccessRes(c, "注册成功", nil)
}

// 定义登录函数
func Login(c *gin.Context) {
	//绑定请求信息到结构体中
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorRes(c, "请求数据异常，登录失败", http.StatusBadRequest)
	}

	//获取请求信息中的用户id
	userName := req.UserName
	//根据用户名去查询用户信息是否存在
	var user models.User
	if err := configs.GetDB().First(&user, "user_name = ?", userName).Error; err != nil {
		//如果不存在返回错误信息
		if err != gorm.ErrRecordNotFound {
			utils.ErrorRes(c, "登录异常", http.StatusBadRequest)
			return
		} else {
			utils.ErrorRes(c, "用户不存在", http.StatusBadRequest)
		}
	}

	//如果存在则获取用户密码和登录密码做校验
	if !user.VerifyPassword(req.Password) {
		utils.ErrorRes(c, "密码错误", http.StatusBadRequest)
		return
	}
	//如果校验成功则返回登录成功信息
	//生成令牌
	token, err := utils.GenerateToken(user.ID, userName)
	//如果令牌生成失败
	if err != nil {
		utils.ErrorRes(c, "登录失败", http.StatusInternalServerError)
		return
	}
	//返回令牌
	utils.SuccessRes(c, "登录成功", gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.UserName,
			"email":    user.Email,
		},
	})

	// utils.SuccessRes(c,  "登录成功", map[string]interface{}{
	// 	"token": token,
	// 	"user": map[string]interface{}{
	// 		"id":       user.ID,
	// 		"username": user.UserName,
	// 		"email":    user.Email,
	// 	},
	// })
}

// GetProfile 获取当前登录用户的个人资料
// @Description: 从Gin上下文获取已认证用户信息并返回
// @param c: Gin上下文，包含已认证用户信息
func GetProfile(c *gin.Context) {
	// 从上下文获取"user"键对应的用户数据
	// 此数据由AuthMiddleware在认证成功后设置到上下文中
	user, exists := c.Get("user")

	// 检查用户数据是否存在
	// 如果不存在，说明认证中间件未正确设置用户信息
	if !exists {
		// 返回401未授权错误，提示用户未找到
		utils.ErrorRes(c, "User not found", http.StatusUnauthorized)
		return
	}

	// 用户数据存在，返回200成功响应
	// 包含"Profile retrieved"消息和用户信息
	utils.SuccessRes(c, "Profile retrieved", user)
}

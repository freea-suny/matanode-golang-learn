/*
 * @Description:
 * @version: 1.0.0
 * @Author: sun.yong
 * @Date: 2025-12-14 10:06:22
 * @LastEditors: sun.yong
 * @LastEditTime: 2025-12-15 23:06:52
 */
package controllers

import (
	"net/http"
	"strconv"
	"task4/configs"
	"task4/models"
	"task4/utils"

	"github.com/gin-gonic/gin"
)

/**
实现博客文章的完整 CRUD 操作，
包括文章创建、查询、更新和删除。
支持文章列表分页、单篇文章详情查看等功能。
使用面向对象的控制器实现方式(空结构体)
*/

type PostController struct{}

type CreatePostRequest struct {
	Title   string `json:"title" binding:"required,min=1,max=100"`
	Content string `json:"content" binding:"required,min=1`
}

type UpdatePostRequest struct {
	ID      uint   `json:"id" binding:"required"`
	Title   string `json:"title" binding:"required,min=1,max=100"`
	Content string `json:"content" binding:"required,min=1`
}

// 创建文章
func (pc *PostController) CreatePost(c *gin.Context) {
	var req CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorRes(c, "请求入参不合法", http.StatusBadRequest)
		return
	}

	//从上下文获取用户ID
	userID := c.GetUint("user_id")
	if userID == 0 {
		utils.ErrorRes(c, "用户未登录", http.StatusUnauthorized)
		return
	}
	//创建实体
	post := models.Post{
		Title:   req.Title,
		Content: req.Content,
		UserID:  uint(userID),
	}
	//保存
	if err := configs.GetDB().Create(&post).Error; err != nil {
		utils.ErrorRes(c, "创建文章失败", http.StatusInternalServerError)
		return
	}
	//返回成功
	configs.GetDB().Preload("User").First(&post, post.ID)
	utils.SuccessRes(c, "创建文章成功", post)

}

// 分页查询
func (pc *PostController) GetPostsPage(c *gin.Context) {
	//从url上获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	//计算offset
	offset := (page - 1) * pageSize

	//查询，预加载用户GetPostsPage
	var posts []models.Post
	err := configs.GetDB().Preload("User").
		Limit(pageSize).
		Offset(offset).
		Find(&posts).
		Order("created_at DESC").Error
	if err != nil {
		utils.ErrorRes(c, "分页查询文章失败", http.StatusInternalServerError)
		return
	}

	//获取总数GetPostsPage
	var total int64
	err = configs.GetDB().Model(&models.Post{}).Count(&total).Error
	if err != nil {
		utils.ErrorRes(c, "获取文章总数失败", http.StatusInternalServerError)
		return
	}

	//返回成功
	utils.SuccessRes(c, "分页查询文章成功", gin.H{
		"posts":    posts,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

// 查询单篇文章详情
func (pc *PostController) GetPostDetail(c *gin.Context) {
	//从url上获取文章id
	id, _ := strconv.Atoi(c.Param("id"))

	//查询，预加载用户GetPostDetail
	var post models.Post
	err := configs.GetDB().Preload("User").
		First(&post, id).Error
	if err != nil {
		utils.ErrorRes(c, "查询文章详情失败", http.StatusInternalServerError)
		return
	}

	//返回成功
	utils.SuccessRes(c, "查询文章详情成功", post)
}

// 修改文章
func (pc *PostController) UpdatePost(c *gin.Context) {
	//绑定参数UpdatePost
	var req UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorRes(c, "请求入参不合法", http.StatusBadRequest)
		return
	}

	//根据文章id查询文章是否存在
	var post models.Post
	if err := configs.GetDB().Preload("User").First(&post, req.ID).Error; err != nil {
		utils.ErrorRes(c, "文章不存在", http.StatusNotFound)
		return
	}

	//检查是否是文章作者
	if post.UserID != c.GetUint("user_id") {
		utils.ErrorRes(c, "您没有权限修改这篇文章", http.StatusForbidden)
		return
	}

	//更新文章
	post.Title = req.Title
	post.Content = req.Content
	err := configs.GetDB().Save(&post).Error
	if err != nil {
		utils.ErrorRes(c, "更新文章失败", http.StatusInternalServerError)
		return
	}
	//返回成功
	utils.SuccessRes(c, "更新文章成功", post)

}

// 删除文章
func (pc *PostController) DeletePost(c *gin.Context) {

	//从url上获取文章id
	id, _ := strconv.Atoi(c.Param("id"))

	//根据文章id查询文章是否存在
	var post models.Post
	if err := configs.GetDB().Preload("User").First(&post, id).Error; err != nil {
		utils.ErrorRes(c, "文章不存在", http.StatusNotFound)
		return
	}

	//检查是否是文章作者
	if post.UserID != c.GetUint("user_id") {
		utils.ErrorRes(c, "您没有权限删除这篇文章", http.StatusForbidden)
		return
	}
	//删除文章
	err := configs.GetDB().Delete(&post).Error
	if err != nil {
		utils.ErrorRes(c, "删除文章失败", http.StatusInternalServerError)
		return
	}
	//返回成功
	utils.SuccessRes(c, "删除文章成功", nil)
}

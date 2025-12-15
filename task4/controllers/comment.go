/*
 * @Description:
 * @version: 1.0.0
 * @Author: sun.yong
 * @Date: 2025-12-14 14:05:51
 * @LastEditors: sun.yong
 * @LastEditTime: 2025-12-15 23:30:03
 */
package controllers

import (
	"net/http"
	"task4/configs"
	"task4/models"
	"task4/utils"

	"github.com/gin-gonic/gin"
)

type CommentController struct{}

/**
处理文章评论的创建和查询功能，
实现用户对文章的互动交流。
支持按文章 ID 查询评论列表，确保评论与文章的正确关联。
*/

// 定义请求结构体
type CommentReq struct {
	PostID   uint   `json:"post_id" binding:"required"`
	Content  string `json:"content"`
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
}

// 创建评论
func (cc *CommentController) CreateComment(c *gin.Context) {
	//绑定请求参数到结构体
	var req CommentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorRes(c, "请求数据异常，创建评论失败", http.StatusBadRequest)
		return
	}

	//校验文章是否存在
	var post models.Post
	if err := configs.GetDB().First(&post, "id = ?", req.PostID).Error; err != nil {
		utils.ErrorRes(c, "文章不存在", http.StatusBadRequest)
		return
	}

	//检查用户是否存在
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorRes(c, "用户不存在", http.StatusBadRequest)
		return
	}
	//创建评论
	comment := models.Comment{
		PostID:  req.PostID,
		UserID:  userID.(uint),
		Content: req.Content,
	}
	if err := configs.GetDB().Create(&comment).Error; err != nil {
		utils.ErrorRes(c, "创建评论失败", http.StatusInternalServerError)
		return
	}
	utils.SuccessRes(c, "创建评论成功", comment)
}

// 获取文章的评论列表(按文章Id分页)
func (cc *CommentController) GetComments(c *gin.Context) {
	//绑定请求参数到结构体
	var req CommentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorRes(c, "请求数据异常，获取评论失败", http.StatusBadRequest)
		return
	}
	//校验文章是否存在
	var post models.Post
	if err := configs.GetDB().First(&post, "id = ?", req.PostID).Error; err != nil {
		utils.ErrorRes(c, "文章不存在", http.StatusBadRequest)
		return
	}
	//查询评论列表
	var comments []models.Comment
	offset := (req.Page - 1) * req.PageSize
	if err := configs.GetDB().Where("post_id = ?", req.PostID).
		Offset(offset).
		Limit(req.PageSize).
		Order("created_at desc").
		Find(&comments).
		Error; err != nil {
		utils.ErrorRes(c, "查询评论失败", http.StatusInternalServerError)
		return
	}

	//获取评论总数
	var total int64
	configs.GetDB().Model(&models.Comment{}).
		Where("post_id = ?", req.PostID).
		Count(&total)
	utils.SuccessRes(c, "查询评论成功", gin.H{
		"comments":  comments,
		"total":     total,
		"page":      req.Page,
		"page_size": req.PageSize,
	})

}

package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
 * @Description: 响应体工具包(用于统一的响应格式)
 * @version: 1.0.0
 * @Author: sun.yong
 * @Date: 2025-12-10 23:11:45
 * @LastEditors: sun.yong
 * @LastEditTime: 2025-12-11 23:00:04
 */

// 定义一个响应的结构体
type Response struct {
	//响应消息
	Msg string `json:"msg"`
	//响应数据
	Data interface{} `json:"data"`
	//是否成功
	Success bool `json:"success"`
}

// 定义成功响应方法
func SuccessRes(c *gin.Context, msg string, data interface{}) {
	//创建新的自定义实例
	res := Response{
		Msg:     msg,
		Data:    data,
		Success: true,
	}

	c.JSON(http.StatusOK, res)
}

// 定义失败响应方法
func ErrorRes(c *gin.Context, msg string, code int) {
	//创建新的自定义实例
	res := Response{
		Msg:     msg,
		Success: false,
	}

	c.JSON(code, res)
}

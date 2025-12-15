/*
 * @Description:
 * @version: 1.0.0
 * @Author: sun.yong
 * @Date: 2025-12-11 22:19:12
 * @LastEditors: sun.yong
 * @LastEditTime: 2025-12-15 14:59:19
 */
package middleware

import (
	"net/http"
	"strings"
	"task4/configs"
	"task4/models"
	"task4/utils"

	"github.com/gin-gonic/gin"
)

/*
 * @Description: 认证中间件,用户请求入口,验证用户是否携带有效令牌,如果有效,则将用户信息存储到上下文,否则返回401错误
 * @version: 1.0.0
 * @Author: sun.yong
 * @Date: 2025-12-11 22:22:00
 * @LastEditors: sun.yong
 * @LastEditTime: 2025-12-11 22:58:48
 */

// 定义一个认证中间件函数
func AuthMiddelWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		//从请求头中获取Authorization(包含token,格式为：Bearer token)
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			//如果没有Authorization头,则返回401错误
			utils.ErrorRes(c, "未授权", http.StatusUnauthorized)
			//直接终止请求
			c.Abort()
			return
		}

		authHeaderParts := strings.Split(authHeader, " ")
		if len(authHeaderParts) != 2 || authHeaderParts[0] != "Bearer" {
			//如果Authorization头格式错误,则返回401错误
			utils.ErrorRes(c, "未授权", http.StatusUnauthorized)
			//直接终止请求
			c.Abort()
			return
		}

		//获取令牌
		token := authHeaderParts[1]

		//解析认证令牌
		clamis, err := utils.ParseToken(token)
		if err != nil {
			//令牌解析失败
			utils.ErrorRes(c, "未授权", http.StatusUnauthorized)
			//直接终止请求
			c.Abort()
			return
		}

		//获取用户信息
		var user models.User
		err = configs.GetDB().First(&user, clamis.UserID).Error
		if err != nil {
			//用户不存在
			utils.ErrorRes(c, "用户不存在", http.StatusUnauthorized)
			//直接终止请求
			c.Abort()
			return
		}

		//将用户信息存储到gin框架的上下文中
		c.Set("user", user)
		c.Set("user_id", user.ID)
		//继续往下执行
		c.Next()

	}

}

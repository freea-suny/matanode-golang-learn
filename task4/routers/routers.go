/*
 * @Description:
 * @version: 1.0.0
 * @Author: sun.yong
 * @Date: 2025-12-14 14:49:17
 * @LastEditors: sun.yong
 * @LastEditTime: 2025-12-15 23:27:33
 */
package routers

import (
	"task4/controllers"
	"task4/middleware"
	"task4/utils"

	"github.com/gin-gonic/gin"
)

/**
定义 API 接口的路由映射，
组织和管理所有 HTTP 端点。
通过路由分组实现权限控制，
将公开接口和需要认证的接口分开管理。
*/

func SetUpRouters() *gin.Engine {
	// 创建新的 Gin 引擎实例
	r := gin.New()

	//使用中间件（类似于拦截器）
	r.Use(middleware.LoggerMiddleware())
	r.Use(middleware.GlobalErrorHandler())
	r.Use(gin.Recovery()) // 恢复中间件，用于处理panic

	//创建控制器实例，类似于controller对象
	postController := &controllers.PostController{}
	commentController := &controllers.CommentController{}

	//根API路由组
	api := r.Group("/api/v1")
	{
		//认证相关的路由无需认证
		auth := api.Group("/auth")
		{
			auth.POST("/register", controllers.Register)
			auth.POST("/login", controllers.Login)
		}
		//需要认证的路由组
		protected := api.Group("/protected")
		protected.Use(middleware.AuthMiddelWare())
		{
			//用户信息
			protected.GET("/profile", controllers.GetProfile)
			//文章相关路由
			posts := protected.Group("/posts")
			{
				posts.POST("/", postController.CreatePost)
				posts.POST("/update", postController.UpdatePost)
				posts.DELETE("/delete/:id", postController.DeletePost)
			}
			//评论相关路由
			comments := protected.Group("/comments")
			{
				comments.POST("/", commentController.CreateComment)
			}
		}

		//公开路由
		public := api.Group("/public")
		{
			public.GET("/posts", postController.GetPostsPage)
			public.GET("/posts/:id", postController.GetPostDetail)
			public.POST("/comments", commentController.GetComments)
		}
	}

	//健康检查SetUpRouters
	r.GET("/health", func(c *gin.Context) {
		utils.SuccessRes(c, "健康检查通过", nil)
	})

	return r
}

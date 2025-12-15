/*
 * @Description:
 * @version: 1.0.0
 * @Author: sun.yong
 * @Date: 2025-12-11 23:18:38
 * @LastEditors: sun.yong
 * @LastEditTime: 2025-12-14 15:32:30
 */
package middleware

import (
	"net/http"
	"runtime/debug"
	"task4/utils"
	"time" // 导入time包用于时间格式化

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

/*
 * @Description: 日志中间件 - 使用Logrus记录结构化HTTP请求日志
 * @version: 1.0.0
 * @Author: sun.yong
 * @Date: 2025-12-11 23:12:48
 * @LastEditors: sun.yong
 * @LastEditTime: 2025-12-14 15:04:08
 */

// LoggerMiddleware 创建一个基于Logrus的Gin日志中间件
// 返回类型：gin.HandlerFunc - Gin框架的中间件函数类型
func LoggerMiddleware() gin.HandlerFunc {
	// 创建独立的Logrus日志实例
	// 使用独立实例而不是全局实例，便于后续扩展和配置隔离
	logger := logrus.New()

	// 设置日志输出格式为JSON
	// JSON格式便于日志收集系统（如ELK、Splunk）解析和索引
	logger.SetFormatter(&logrus.JSONFormatter{})

	// 使用Gin的LoggerWithFormatter方法，自定义日志格式化逻辑
	// param gin.LogFormatterParams：Gin提供的日志格式化参数结构
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// 使用Logrus记录结构化HTTP请求日志
		logger.WithFields(logrus.Fields{
			"status_code": param.StatusCode,                     // HTTP响应状态码（如200、404、500）
			"latency":     param.Latency,                        // 请求处理耗时（如12.345678ms）
			"client_ip":   param.ClientIP,                       // 客户端IP地址
			"method":      param.Method,                         // HTTP请求方法（如GET、POST、PUT、DELETE）
			"path":        param.Path,                           // 请求路径（如/api/v1/posts）
			"user_agent":  param.Request.UserAgent(),            // 客户端User-Agent信息
			"error":       param.ErrorMessage,                   // 错误信息（如果有）
			"timestamp":   param.TimeStamp.Format(time.RFC3339), // 日志时间戳（RFC3339格式）
		}).Info("HTTP Request") // 使用Info级别记录日志，日志消息为"HTTP Request"

		// 返回空字符串，禁用Gin默认的日志输出
		// 因为我们已经使用Logrus记录了自定义格式的日志
		return ""
	})
}

// 全局异常处理中间件
func GlobalErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 使用defer确保异常处理函数在任何情况下都会执行 （包括c.Next()之后的代码）
		defer func() {
			//捕获panic异常
			if r := recover(); r != nil {
				//记录结构化日志
				logrus.WithFields(logrus.Fields{
					"error":  r,
					"path":   c.Request.URL.Path,
					"method": c.Request.Method,
					"ip":     c.ClientIP(),
					"stack":  string(debug.Stack()),
				}).Error("Recovered from panic")

				//统一返回错误响应
				utils.ErrorRes(c, "服务器内部错误", http.StatusInternalServerError)

				// 终止后续中间件和路由处理
				c.Abort()
			}
		}()
		// 继续处理后续中间件和路由
		c.Next()
	}
}

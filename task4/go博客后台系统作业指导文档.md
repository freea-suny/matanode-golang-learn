# Go 语言博客项目详细开发指南（MySQL 版本）

## 第一阶段：项目基础搭建（1-2 天）

### 步骤 1：项目目录结构创建

**作用描述**：建立清晰的项目架构，为后续开发提供组织框架。合理的目录结构有助于代码维护、团队协作和功能模块划分，确保项目具有良好的可扩展性。

```
golang_blog/
├── cmd/
│   └── main.go                    # 程序入口
├── config/
│   └── database.go               # 数据库配置
├── controllers/                   # 控制器层
├── middleware/                   # 中间件层
├── models/                       # 数据模型层
├── routes/                       # 路由配置
├── utils/                        # 工具函数
├── go.mod                        # 模块定义
├── go.sum                        # 依赖锁定
├── .env                          # 环境变量
└── .env.example                 # 环境变量模板
```

### 步骤 2：初始化项目依赖

**作用描述**：定义项目所需的第三方库依赖，确保开发环境一致性。通过 go.mod 文件管理依赖版本，避免因依赖版本不一致导致的兼容性问题。

创建 `go.mod` 文件：

```go:/d:/workspace/web3/golang_blog-master/golang_blog-master/go.mod
module golang_blog

go 1.24.2

require (
    github.com/gin-gonic/gin v1.11.0
    github.com/golang-jwt/jwt/v5 v5.3.0
    github.com/sirupsen/logrus v1.9.3
    github.com/joho/godotenv v1.5.1
    gorm.io/driver/mysql v1.6.0
    gorm.io/gorm v1.31.0
)
```

### 步骤 3：环境变量配置

**作用描述**：实现配置与代码分离，提高项目的可移植性和安全性。通过环境变量管理敏感信息（如数据库密码、JWT 密钥），便于不同环境（开发/测试/生产）的部署。

创建 `.env.example` 文件：

```bash:/d:/workspace/web3/golang_blog-master/golang_blog-master/.env.example
# MySQL数据库配置
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_mysql_password
DB_NAME=golang_blog

# JWT配置
JWT_SECRET=your_jwt_secret_key_here

# 服务器配置
PORT=8080
GIN_MODE=debug
```

## 第二阶段：数据库配置（1 天）

### 步骤 4：MySQL 数据库配置

**作用描述**：建立与 MySQL 数据库的连接，为数据持久化提供基础。通过 GORM ORM 框架简化数据库操作，实现数据库连接的统一管理和复用。

创建 `config/database.go` 文件：

```go:/d:/workspace/web3/golang_blog-master/golang_blog-master/config/database.go
package config

import (
    "fmt"
    "log"
    "os"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDatabase() {
    var err error

    // 从环境变量获取配置
    dbHost := getEnv("DB_HOST", "localhost")
    dbPort := getEnv("DB_PORT", "3306")
    dbUser := getEnv("DB_USER", "root")
    dbPassword := getEnv("DB_PASSWORD", "")
    dbName := getEnv("DB_NAME", "golang_blog")

    // MySQL连接字符串
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        dbUser, dbPassword, dbHost, dbPort, dbName)

    // 连接数据库
    DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
    })

    if err != nil {
        log.Fatal("Failed to connect to MySQL:", err)
    }

    log.Println("MySQL connected successfully")
}

func GetDB() *gorm.DB {
    return DB
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}
```

## 第三阶段：数据模型设计（2 天）

### 步骤 5：用户模型

**作用描述**：定义用户数据结构和业务逻辑，包括用户注册、登录、密码加密等功能。通过 GORM 标签实现数据库表结构映射，确保数据完整性和安全性。

创建 `models/user.go` 文件：

```go:/d:/workspace/web3/golang_blog-master/golang_blog-master/models/user.go
package models

import (
    "time"
    "golang.org/x/crypto/bcrypt"
    "gorm.io/gorm"
)

type User struct {
    ID        uint           `gorm:"primaryKey" json:"id"`
    Username  string         `gorm:"uniqueIndex;size:50;not null" json:"username"`
    Email     string         `gorm:"uniqueIndex;size:100;not null" json:"email"`
    Password  string         `gorm:"size:255;not null" json:"-"`
    Posts     []Post         `gorm:"foreignKey:UserID" json:"posts,omitempty"`
    Comments  []Comment      `gorm:"foreignKey:UserID" json:"comments,omitempty"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// 密码加密
func (u *User) BeforeSave(tx *gorm.DB) error {
    if u.Password != "" {
        hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
        if err != nil {
            return err
        }
        u.Password = string(hashedPassword)
    }
    return nil
}

// 验证密码
func (u *User) CheckPassword(password string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
    return err == nil
}
```

### 步骤 6：文章模型

**作用描述**：定义博客文章的数据结构，建立与用户模型的关联关系。实现文章的基本 CRUD 操作，支持文章与评论的一对多关系。

创建 `models/post.go` 文件：

```go:/d:/workspace/web3/golang_blog-master/golang_blog-master/models/post.go
package models

import (
    "time"
    "gorm.io/gorm"
)

type Post struct {
    ID        uint           `gorm:"primaryKey" json:"id"`
    Title     string         `gorm:"size:200;not null" json:"title"`
    Content   string         `gorm:"type:text;not null" json:"content"`
    UserID    uint           `gorm:"not null;index" json:"user_id"`
    User      User           `gorm:"foreignKey:UserID" json:"user"`
    Comments  []Comment      `gorm:"foreignKey:PostID" json:"comments,omitempty"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
```

### 步骤 7：评论模型

**作用描述**：定义评论数据模型，建立与用户和文章的双向关联。实现评论的创建和查询功能，支持评论的级联删除。

创建 `models/comment.go` 文件：

```go:/d:/workspace/web3/golang_blog-master/golang_blog-master/models/comment.go
package models

import (
    "time"
    "gorm.io/gorm"
)

type Comment struct {
    ID        uint           `gorm:"primaryKey" json:"id"`
    Content   string         `gorm:"type:text;not null" json:"content"`
    UserID    uint           `gorm:"not null;index" json:"user_id"`
    User      User           `gorm:"foreignKey:UserID" json:"user"`
    PostID    uint           `gorm:"not null;index" json:"post_id"`
    Post      Post           `gorm:"foreignKey:PostID" json:"post"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
```

## 第四阶段：工具函数开发（1 天）

### 步骤 8：JWT 工具

**作用描述**：实现用户身份认证的核心功能，包括 JWT 令牌的生成和验证。通过 JWT 实现无状态的身份验证机制，支持用户会话管理。

创建 `utils/jwt.go` 文件：

```go:/d:/workspace/web3/golang_blog-master/golang_blog-master/utils/jwt.go
package utils

import (
    "time"
    "github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(getEnv("JWT_SECRET", "default_secret"))

type Claims struct {
    UserID uint `json:"user_id"`
    jwt.RegisteredClaims
}

func GenerateToken(userID uint) (string, error) {
    claims := Claims{
        UserID: userID,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Issuer:    "golang_blog",
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtSecret)
}

func ParseToken(tokenString string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        return jwtSecret, nil
    })

    if err != nil {
        return nil, err
    }

    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }

    return nil, err
}
```

### 步骤 9：响应工具

**作用描述**：统一 API 响应格式，提高代码复用性和可维护性。通过标准化的响应结构，确保前端能够正确处理成功和错误情况。

创建 `utils/response.go` 文件：

```go:/d:/workspace/web3/golang_blog-master/golang_blog-master/utils/response.go
package utils

import "github.com/gin-gonic/gin"

type Response struct {
    Success bool        `json:"success"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
    Error   string      `json:"error,omitempty"`
}

func Success(c *gin.Context, code int, message string, data interface{}) {
    c.JSON(code, Response{
        Success: true,
        Message: message,
        Data:    data,
    })
}

func Error(c *gin.Context, code int, message string, err error) {
    response := Response{
        Success: false,
        Message: message,
    }

    if err != nil {
        response.Error = err.Error()
    }

    c.JSON(code, response)
}
```

## 第五阶段：中间件开发（2 天）

### 步骤 10：认证中间件

**作用描述**：实现请求级别的身份验证，保护需要登录才能访问的 API 接口。通过中间件拦截请求，验证 JWT 令牌的有效性，确保只有合法用户才能访问受保护资源。

创建 `middleware/auth.go` 文件：

```go:/d:/workspace/web3/golang_blog-master/golang_blog-master/middleware/auth.go
package middleware

import (
    "net/http"
    "strings"
    "golang_blog/config"
    "golang_blog/models"
    "golang_blog/utils"
    "github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            utils.Error(c, http.StatusUnauthorized, "Authorization header required", nil)
            c.Abort()
            return
        }

        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            utils.Error(c, http.StatusUnauthorized, "Invalid authorization format", nil)
            c.Abort()
            return
        }

        claims, err := utils.ParseToken(parts[1])
        if err != nil {
            utils.Error(c, http.StatusUnauthorized, "Invalid token", err)
            c.Abort()
            return
        }

        var user models.User
        if err := config.GetDB().First(&user, claims.UserID).Error; err != nil {
            utils.Error(c, http.StatusUnauthorized, "User not found", err)
            c.Abort()
            return
        }

        c.Set("user", user)
        c.Set("userID", user.ID)
        c.Next()
    }
}
```

### 步骤 11：日志中间件（在 Gin 框架中，中间件是一种可以拦截 HTTP 请求和响应的特殊函数，它可以在请求到达控制器之前或响应返回客户端之前执行一些操作）

**作用描述**：记录 HTTP 请求的详细信息，便于问题排查和性能监控。通过日志中间件收集请求方法、路径、状态码、响应时间等关键信息。

创建 `middleware/logger.go` 文件：

```go:/d:/workspace/web3/golang_blog-master/golang_blog-master/middleware/logger.go
package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// LoggerMiddleware 日志记录中间件
func LoggerMiddleware() gin.HandlerFunc {
	// 创建logrus实例
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// 记录请求日志
		logger.WithFields(logrus.Fields{
			"status_code":  param.StatusCode,
			"latency":      param.Latency,
			"client_ip":    param.ClientIP,
			"method":       param.Method,
			"path":         param.Path,
			"user_agent":   param.Request.UserAgent(),
			"error":        param.ErrorMessage,
			"timestamp":    param.TimeStamp.Format(time.RFC3339),
		}).Info("HTTP Request")

		return ""
	})
}

// ErrorHandlerMiddleware 全局错误处理中间件
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logrus.WithFields(logrus.Fields{
					"error": err,
					"path":  c.Request.URL.Path,
					"method": c.Request.Method,
				}).Error("Panic recovered")

				c.JSON(500, gin.H{
					"code":    500,
					"message": "Internal server error",
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}
```

## 第六阶段：控制器开发（3-4 天）

### 步骤 12：认证控制器

**作用描述**：处理用户注册、登录和个人资料查询等认证相关功能。实现用户身份验证流程，包括密码验证、JWT 令牌生成和用户信息管理。

创建 `controllers/auth.go` 文件：

```go:/d:/workspace/web3/golang_blog-master/golang_blog-master/controllers/auth.go
package controllers

import (
    "net/http"
    "golang_blog/config"
    "golang_blog/models"
    "golang_blog/utils"
    "github.com/gin-gonic/gin"
)

type RegisterRequest struct {
    Username string `json:"username" binding:"required,min=3,max=50"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}

func Register(c *gin.Context) {
    var req RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        utils.Error(c, http.StatusBadRequest, "Invalid request data", err)
        return
    }

    // 检查用户名是否已存在
    var existingUser models.User
    if err := config.GetDB().Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
        utils.Error(c, http.StatusConflict, "Username already exists", nil)
        return
    }

    // 检查邮箱是否已存在
    if err := config.GetDB().Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
        utils.Error(c, http.StatusConflict, "Email already exists", nil)
        return
    }

    user := models.User{
        Username: req.Username,
        Email:    req.Email,
        Password: req.Password,
    }

    if err := config.GetDB().Create(&user).Error; err != nil {
        utils.Error(c, http.StatusInternalServerError, "Failed to create user", err)
        return
    }

    utils.Success(c, http.StatusCreated, "User registered successfully", gin.H{
        "id":       user.ID,
        "username": user.Username,
        "email":    user.Email,
    })
}

func Login(c *gin.Context) {
    var req LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        utils.Error(c, http.StatusBadRequest, "Invalid request data", err)
        return
    }

    var user models.User
    if err := config.GetDB().Where("username = ?", req.Username).First(&user).Error; err != nil {
        utils.Error(c, http.StatusUnauthorized, "Invalid credentials", nil)
        return
    }

    if !user.CheckPassword(req.Password) {
        utils.Error(c, http.StatusUnauthorized, "Invalid credentials", nil)
        return
    }

    token, err := utils.GenerateToken(user.ID)
    if err != nil {
        utils.Error(c, http.StatusInternalServerError, "Failed to generate token", err)
        return
    }

    utils.Success(c, http.StatusOK, "Login successful", gin.H{
        "token": token,
        "user": gin.H{
            "id":       user.ID,
            "username": user.Username,
            "email":    user.Email,
        },
    })
}

func GetProfile(c *gin.Context) {
    user, exists := c.Get("user")
    if !exists {
        utils.Error(c, http.StatusUnauthorized, "User not found", nil)
        return
    }

    utils.Success(c, http.StatusOK, "Profile retrieved", user)
}
```

### 步骤 13：文章控制器

**作用描述**：实现博客文章的完整 CRUD 操作，包括文章创建、查询、更新和删除。支持文章列表分页、单篇文章详情查看等功能。

创建 `controllers/post.go` 文件：

```go:/d:/workspace/web3/golang_blog-master/golang_blog-master/controllers/post.go
package controllers

import (
	"strconv"

	"golang_blog/config"
	"golang_blog/models"
	"golang_blog/utils"
	"github.com/gin-gonic/gin"
)

type PostController struct{}

type CreatePostRequest struct {
	Title   string `json:"title" binding:"required,min=1,max=200"`
	Content string `json:"content" binding:"required,min=1"`
}

type UpdatePostRequest struct {
	Title   string `json:"title" binding:"required,min=1,max=200"`
	Content string `json:"content" binding:"required,min=1"`
}

// CreatePost 创建文章
func (pc *PostController) CreatePost(c *gin.Context) {
	var req CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "User not authenticated")
		return
	}

	post := models.Post{
		Title:   req.Title,
		Content: req.Content,
		UserID:  userID.(uint),
	}

	if err := config.DB.Create(&post).Error; err != nil {
		utils.InternalServerError(c, "Failed to create post")
		return
	}

	// 预加载用户信息
	config.DB.Preload("User").First(&post, post.ID)

	utils.Success(c, post)
}

// GetPosts 获取文章列表
func (pc *PostController) GetPosts(c *gin.Context) {
	var posts []models.Post

	// 分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	// 查询文章列表，预加载用户信息
	if err := config.DB.Preload("User").
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&posts).Error; err != nil {
		utils.InternalServerError(c, "Failed to get posts")
		return
	}

	// 获取总数
	var total int64
	config.DB.Model(&models.Post{}).Count(&total)

	utils.Success(c, gin.H{
		"posts":     posts,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetPost 获取单个文章详情
func (pc *PostController) GetPost(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid post ID")
		return
	}

	var post models.Post
	if err := config.DB.Preload("User").Preload("Comments.User").First(&post, postID).Error; err != nil {
		utils.NotFound(c, "Post not found")
		return
	}

	utils.Success(c, post)
}

// UpdatePost 更新文章
func (pc *PostController) UpdatePost(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid post ID")
		return
	}

	var req UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "User not authenticated")
		return
	}

	var post models.Post
	if err := config.DB.First(&post, postID).Error; err != nil {
		utils.NotFound(c, "Post not found")
		return
	}

	// 检查是否是文章作者
	if post.UserID != userID.(uint) {
		utils.Forbidden(c, "You can only update your own posts")
		return
	}

	// 更新文章
	post.Title = req.Title
	post.Content = req.Content

	if err := config.DB.Save(&post).Error; err != nil {
		utils.InternalServerError(c, "Failed to update post")
		return
	}

	// 预加载用户信息
	config.DB.Preload("User").First(&post, post.ID)

	utils.Success(c, post)
}

// DeletePost 删除文章
func (pc *PostController) DeletePost(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid post ID")
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "User not authenticated")
		return
	}

	var post models.Post
	if err := config.DB.First(&post, postID).Error; err != nil {
		utils.NotFound(c, "Post not found")
		return
	}

	// 检查是否是文章作者
	if post.UserID != userID.(uint) {
		utils.Forbidden(c, "You can only delete your own posts")
		return
	}

	// 删除文章（软删除）
	if err := config.DB.Delete(&post).Error; err != nil {
		utils.InternalServerError(c, "Failed to delete post")
		return
	}

	utils.Success(c, gin.H{"message": "Post deleted successfully"})
}
```

### 步骤 14：评论控制器

**作用描述**：处理文章评论的创建和查询功能，实现用户对文章的互动交流。支持按文章 ID 查询评论列表，确保评论与文章的正确关联。

创建 `controllers/comment.go` 文件：

```go:/d:/workspace/web3/golang_blog-master/golang_blog-master/controllers/comment.go
package controllers

import (
	"strconv"

	"golang_blog/config"
	"golang_blog/models"
	"golang_blog/utils"
	"github.com/gin-gonic/gin"
)

type CommentController struct{}

type CreateCommentRequest struct {
	Content string `json:"content" binding:"required,min=1,max=1000"`
}

// CreateComment 创建评论
func (cc *CommentController) CreateComment(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("post_id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid post ID")
		return
	}

	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "User not authenticated")
		return
	}

	// 检查文章是否存在
	var post models.Post
	if err := config.DB.First(&post, postID).Error; err != nil {
		utils.NotFound(c, "Post not found")
		return
	}

	comment := models.Comment{
		Content: req.Content,
		UserID:  userID.(uint),
		PostID:  uint(postID),
	}

	if err := config.DB.Create(&comment).Error; err != nil {
		utils.InternalServerError(c, "Failed to create comment")
		return
	}

	// 预加载用户信息
	config.DB.Preload("User").First(&comment, comment.ID)

	utils.Success(c, comment)
}

// GetComments 获取文章的评论列表
func (cc *CommentController) GetComments(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("post_id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid post ID")
		return
	}

	// 检查文章是否存在
	var post models.Post
	if err := config.DB.First(&post, postID).Error; err != nil {
		utils.NotFound(c, "Post not found")
		return
	}

	var comments []models.Comment

	// 分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	// 查询评论列表，预加载用户信息
	if err := config.DB.Preload("User").
		Where("post_id = ?", postID).
		Order("created_at ASC").
		Limit(pageSize).
		Offset(offset).
		Find(&comments).Error; err != nil {
		utils.InternalServerError(c, "Failed to get comments")
		return
	}

	// 获取总数
	var total int64
	config.DB.Model(&models.Comment{}).Where("post_id = ?", postID).Count(&total)

	utils.Success(c, gin.H{
		"comments":  comments,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}
```

## 第七阶段：路由配置（1 天）

### 步骤 15：路由配置

**作用描述**：定义 API 接口的路由映射，组织和管理所有 HTTP 端点。通过路由分组实现权限控制，将公开接口和需要认证的接口分开管理。

创建 `routes/routes.go` 文件：

```go:/d:/workspace/web3/golang_blog-master/golang_blog-master/routes/routes.go
package routes

import (
	"golang_blog/controllers"
	"golang_blog/middleware"
	"github.com/gin-gonic/gin"
)

// SetupRoutes 设置路由
func SetupRoutes() *gin.Engine {
	r := gin.New()

	// 使用中间件
	r.Use(middleware.LoggerMiddleware())
	r.Use(middleware.ErrorHandlerMiddleware())
	r.Use(gin.Recovery())

	// 创建控制器实例
	authController := &controllers.AuthController{}
	postController := &controllers.PostController{}
	commentController := &controllers.CommentController{}

	// API路由组
	api := r.Group("/api/v1")
	{
		// 认证相关路由（无需认证）
		auth := api.Group("/auth")
		{
			auth.POST("/register", authController.Register)
			auth.POST("/login", authController.Login)
		}

		// 需要认证的路由
		authenticated := api.Group("")
		authenticated.Use(middleware.AuthMiddleware())
		{
			// 用户信息
			authenticated.GET("/profile", authController.GetProfile)

			// 文章相关路由
			posts := authenticated.Group("/posts")
			{
				posts.POST("", postController.CreatePost)
				posts.PUT("/:id", postController.UpdatePost)
				posts.DELETE("/:id", postController.DeletePost)
			}

			// 评论相关路由
			comments := authenticated.Group("/posts/:post_id/comments")
			{
				comments.POST("", commentController.CreateComment)
			}
		}

		// 公开路由（无需认证）
		public := api.Group("")
		{
			// 文章公开路由
			public.GET("/posts", postController.GetPosts)
			public.GET("/posts/:id", postController.GetPost)
		}

		// 评论公开路由（单独分组避免路由冲突）
		comments := api.Group("/comments")
		{
			comments.GET("/post/:post_id", commentController.GetComments)
		}
	}

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"message": "Blog API is running",
		})
	})

	return r
}
```

## 第八阶段：项目整合和测试（2 天）

### 步骤 16：程序入口

**作用描述**：整合所有模块，启动 HTTP 服务器。负责初始化数据库连接、加载环境变量、设置路由并启动服务监听。

创建 `cmd/main.go` 文件：

```go:/d:/workspace/web3/golang_blog-master/golang_blog-master/cmd/main.go
package main

import (
	"log"
	"os"

	"golang_blog/config"
	"golang_blog/routes"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		logrus.Warn("No .env file found, using system environment variables")
	}

	// 初始化日志
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.Info("Starting blog application...")

	// 初始化数据库
	config.InitDatabase()

	// 设置路由
	r := routes.SetupRoutes()

	// 启动服务器
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	port = ":" + port

	logrus.WithField("port", port).Info("Server starting")

	if err := r.Run(port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
```

### 步骤 17：安装依赖并运行

**作用描述**：完成项目的最终部署准备，确保所有依赖正确安装，项目能够正常运行。

```bash
# 安装依赖
go mod tidy

# 运行项目
go run cmd/main.go
```

## 开发时间安排建议

**第一周：**

- 周一：项目基础搭建 + 环境配置
- 周二：数据库配置 + MySQL 学习
- 周三：数据模型设计
- 周四：工具函数开发
- 周五：中间件开发

**第二周：**

- 周一：认证控制器
- 周二：文章控制器
- 周三：评论控制器
- 周四：路由配置
- 周五：项目整合测试

**第三周：**

- 功能完善和优化
- API 测试
- 部署准备

这个增强版的开发指南为每个步骤都添加了详细的作用描述，帮助你理解每个开发环节在整个项目中的具体作用和意义。通过这种方式，你不仅知道如何实现代码，还能理解为什么要这样实现，从而更好地掌握整个项目的架构和设计思路。

/*
 * @Description: jwt工具包
 * @version: 1.0.0
 * @Author: sun.yong
 * @Date: 2025-12-10 22:13:21
 * @LastEditors: sun.yong
 * @LastEditTime: 2025-12-11 22:21:34
 */
package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// 定义一个令牌结构体
// Claims 自定义声明类型 并内嵌jwt.RegisteredClaims
// jwt包自带的jwt.RegisteredClaims只包含了官方字段
// 假设我们这里需要额外记录一个username字段，所以要自定义结构体
// 如果想要保存更多信息，都可以添加到这个结构体中
type Claims struct {
	//自定义存储用户信息
	UserID   uint   `json:"user_id"`
	UserName string `json:"user_name"`
	jwt.RegisteredClaims
}

func getEnv() string {
	if value := os.Getenv("JWT_SECRET"); value == "" {
		return "default_secret_key"
	} else {
		return value
	}
}

// 获取密钥
var secretKey = []byte(getEnv())

// 生成令牌
func GenerateToken(userID uint, userName string) (string, error) {
	//创建新的自定义实例
	claims := Claims{
		UserID:   userID,
		UserName: userName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 自定义过期时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),                     // 自定义签发时间
			Issuer:    "suny",                                             //自定义签发人
		},
	}

	//使用HS256算法生成令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//使用密钥对令牌签名
	return token.SignedString(secretKey)
}

// 解析令牌是否有效
func ParseToken(tokenString string) (*Claims, error) {
	//解析令牌
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}
	// 类型断言，将解析后的Claims转换为自定义的Claims类型
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("不合法的token")
}

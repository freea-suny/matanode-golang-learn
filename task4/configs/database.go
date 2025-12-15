/*
 * @Description:
 * @version: 1.0.0
 * @Author: sun.yong
 * @Date: 2025-12-09 11:34:45
 * @LastEditors: sun.yong
 * @LastEditTime: 2025-12-15 15:33:15
 */
package configs

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 定义一个全局的数据库连接对象
var DB *gorm.DB

// 初始化数据库连接
func InitDB() {
	//从环境变量中获取数据库配置信息
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "3306")
	user := getEnv("DB_USER", "root")
	password := getEnv("DB_PASSWORD", "1234")
	dbname := getEnv("DB_NAME", "matanode")

	//mysql连接dsn
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, dbname)

	//连接数据库
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		//gorm自带的日志（用于sql的执行记录）
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		//go自带的标准库日志
		log.Fatal("连接数据库失败:", err)
	}

	log.Println("连接数据库成功")

}

func GetDB() *gorm.DB {
	return DB
}

func getEnv(key string, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

/*
 * @Description:
 * @version: 1.0.0
 * @Author: sun.yong
 * @Date: 2025-12-14 15:59:17
 * @LastEditors: sun.yong
 * @LastEditTime: 2025-12-15 15:33:01
 */
package main

import (
	"os"
	"path/filepath"
	"task4/configs"
	"task4/routers"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

/**
整合所有模块，
启动 HTTP 服务器。
负责初始化数据库连接、
加载环境变量、
设置路由并启动服务监听。
*/

func main() {
	//加载环境变量
	// 尝试多种方式加载.env文件
	envFiles := []string{
		".env",
		"../.env",
		"./.env",
		filepath.Join(os.Getenv("PWD"), ".env"),
	}

	loaded := false
	for _, envFile := range envFiles {
		if err := godotenv.Load(envFile); err == nil {
			logrus.Infof("成功加载环境变量文件: %s", envFile)
			loaded = true
			break
		}
	}

	if !loaded {
		logrus.Warn("未找到.env文件，将使用系统环境变量")
	}

	//初始化日志
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.Info("开始启动博客服务后台...")

	//初始化数据库
	configs.InitDB()

	//设置路由
	router := routers.SetUpRouters()

	//启动服务器
	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "8080"
	}
	port := ":" + serverPort

	logrus.WithField("port", port).Info("博客服务后台启动成功")
	if err := router.Run(port); err != nil {
		logrus.WithError(err).Fatal("博客服务后台启动失败")
	}

}

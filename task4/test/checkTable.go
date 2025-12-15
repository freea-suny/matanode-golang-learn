/*
 * @Description:
 * @version: 1.0.0
 * @Author: sun.yong
 * @Date: 2025-12-14 16:25:34
 * @LastEditors: sun.yong
 * @LastEditTime: 2025-12-14 16:30:10
 */
package main

import (
	"task4/configs"
	"task4/models"

	"github.com/sirupsen/logrus"
)

func main() {

	configs.InitDB()

	//检查数据库表是否存在
	db := configs.GetDB()

	//初始化日志
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.Info("开始初始化表")

	//检查数据库表是否存在
	if db.Migrator().HasTable(&models.User{}) {
		logrus.Info("用户表已存在")
	} else {
		logrus.Info("用户表不存在，创建用户表")
		db.AutoMigrate(&models.User{})
	}
	//检查数据库表是否存在
	if db.Migrator().HasTable(&models.Post{}) {
		logrus.Info("文章表已存在")
	} else {
		logrus.Info("文章表不存在，创建文章表")
		db.AutoMigrate(&models.Post{})
	}
	//检查数据库表是否存在
	if db.Migrator().HasTable(&models.Comment{}) {
		logrus.Info("评论表已存在")
	} else {
		logrus.Info("评论表不存在，创建评论表")
		db.AutoMigrate(&models.Comment{})
	}
	logrus.Info("初始化表结束")

}

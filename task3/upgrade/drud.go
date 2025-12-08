/*
 * @Description:
 * @version: 1.0.0
 * @Author: sun.yong
 * @Date: 2025-12-05 11:04:29
 * @LastEditors: sun.yong
 * @LastEditTime: 2025-12-05 13:42:18
 */
package upgrade

import (
	"fmt"

	"gorm.io/gorm"
)

/**
题目1：基本CRUD操作
假设有一个名为 students 的表，包含字段 id （主键，自增）、 name （学生姓名，字符串类型）、
age （学生年龄，整数类型）、 grade （学生年级，字符串类型）。
要求 ：
编写SQL语句向 students 表中插入一条新记录，学生姓名为 "张三"，年龄为 20，年级为 "三年级"。
编写SQL语句查询 students 表中所有年龄大于 18 岁的学生信息。
编写SQL语句将 students 表中姓名为 "张三" 的学生年级更新为 "四年级"。
编写SQL语句删除 students 表中年龄小于 15 岁的学生记录
*/

// 实体
type student struct {
	Id    uint `gorm:"primaryKey","autoIncrement"`
	Name  string
	Age   int
	Grade string
}

func Drud(db *gorm.DB) {
	//新增该表
	err := db.AutoMigrate(&student{})
	if err != nil {
		panic("自动迁移失败: " + err.Error())
	}

	//插入一条新记录，学生姓名为 "张三"，年龄为 20，年级为 "三年级"。
	result := db.Create(&student{Name: "张三", Age: 20, Grade: "三年级"})
	if result.Error != nil {
		panic("创建记录出错: " + result.Error.Error())
	} else {
		fmt.Println("创建完成")
	}

	//查询 students 表中所有年龄大于 18 岁的学生信息。
	var students []student
	db.Where("age > ?", 18).Find(&students)
	for _, stu := range students {
		println("年龄大于18岁的学生姓名为:", stu.Name)
	}

	//更新，将 students 表中姓名为 "张三" 的学生年级更新为 "四年级"。
	db.Model(&student{}).Where("name = ? ", "张三").Update("grade", "四年级")
	fmt.Println("更新完成")

	//删除 students 表中年龄小于 15 岁的学生记录
	db.Where("age<?", 15).Delete(&student{})
	fmt.Println("删除完成")

}

/*
 * @Description:
 * @version: 1.0.0
 * @Author: sun.yong
 * @Date: 2025-12-04 17:59:58
 * @LastEditors: sun.yong
 * @LastEditTime: 2025-12-08 09:18:08
 */
package main

import (
	"fmt"
	"task3/constant"
	"task3/upgrade"

	"github.com/glebarez/sqlite"
	// _ "github.com/logoove/sqlite"
	"gorm.io/gorm"
)

// 定义一个与数据库表对应的模型结构体
type User struct {
	ID    uint `gorm:"primaryKey"`
	Name  string
	Email string `gorm:"uniqueIndex"`
}

func main() {
	db, err := gorm.Open(sqlite.Open(constant.DBPath), &gorm.Config{})
	if err != nil {
		panic("连接数据库失败: " + err.Error())
	}
	//执行CRUD操作
	// upgrade.Drud(db)

	//充值
	// upgrade.Charge(db, 1, 100)
	// upgrade.Charge(db, 2, 100)

	//转帐
	// fmt.Println(upgrade.Transfer(db, 1, 2, 100))

	//sqlx使用
	// employees, err := upgrade.QueryEmployeesByDepartment("技术部")
	// if err != nil {
	// 	panic("查询员工失败: " + err.Error())
	// }
	// fmt.Printf("%+v\n", employees)

	// highestSalary, err := upgrade.QureyHighestSalary()
	// if err != nil {
	// 	panic("查询最高工资失败: " + err.Error())
	// }
	// fmt.Printf("最高工资: %v\n", highestSalary)

	// books, err := upgrade.QueryBooksByPrice(50)
	// if err != nil {
	// 	panic("查询图书失败: " + err.Error())
	// }
	// fmt.Printf("%+v\n", books)

	// upgrade.CreatTable(db)
	post, _ := upgrade.QueryPostWithMaxComment(db)
	fmt.Printf("%+v\n", post)

}

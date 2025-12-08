/*
 * @Description:
 * @version: 1.0.0
 * @Author: sun.yong
 * @Date: 2025-12-05 15:49:47
 * @LastEditors: sun.yong
 * @LastEditTime: 2025-12-07 18:33:57
 */
package upgrade

import (
	"task3/constant"

	"github.com/jmoiron/sqlx"
)

/*
*
题目1：使用SQL扩展库进行查询
假设你已经使用Sqlx连接到一个数据库，并且有一个 employees 表，包含字段 id 、 name 、 department 、 salary 。
要求 ：
编写Go代码，使用Sqlx查询 employees 表中所有部门为 "技术部" 的员工信息，并将结果映射到一个自定义的 Employee 结构体切片中。
编写Go代码，使用Sqlx查询 employees 表中工资最高的员工信息，并将结果映射到一个 Employee 结构体中。
*/

type Employee struct {
	ID         int     `db:"id"`
	Name       string  `db:"name"`
	Department string  `db:"department"`
	Salary     float64 `db:"salary"`
}

func QueryEmployeesByDepartment(department string) ([]Employee, error) {
	dbx, err := sqlx.Open("sqlite", constant.DBPath) // 注意驱动名是 sqlite
	if err != nil {
		panic("连接数据库失败: " + err.Error())
	}
	defer dbx.Close()

	//注：这个go的数据类型对应的是sqlite的数据类型
	dbx.MustExec(`create table if not exists employees (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT, department TEXT, salary REAL);`)
	//mustexec是不需要返回值的，因为它是用来执行数据库操作的，而不是用来查询数据的，遇到问题会直接是程序崩溃
	dbx.MustExec("INSERT INTO employees(name, department, salary) VALUES (?, ?, ?)", "张三", "技术部", 5000)
	dbx.MustExec("INSERT INTO employees(name, department, salary) VALUES (?, ?, ?)", "李四", "技术部", 6000)
	dbx.MustExec("INSERT INTO employees(name, department, salary) VALUES (?, ?, ?)", "王五", "技术部", 7000)

	//查询返回切片
	var employees []Employee
	err1 := dbx.Select(&employees, "SELECT * FROM employees WHERE department = ?", department)
	if err1 != nil {
		return nil, err1
	}
	return employees, nil
}

func QureyHighestSalary() (Employee, error) {
	dbx, err := sqlx.Open("sqlite", constant.DBPath) // 注意驱动名是 sqlite
	if err != nil {
		panic("连接数据库失败: " + err.Error())
	}
	defer dbx.Close()

	//查询返回单条记录
	var employee Employee
	err2 := dbx.Get(&employee, "SELECT * FROM employees ORDER BY salary DESC LIMIT 1")
	if err2 != nil {
		return employee, err2
	}
	return employee, nil
}

/*
*
题目2：实现类型安全映射
假设有一个 books 表，包含字段 id 、 title 、 author 、 price 。
要求 ：
定义一个 Book 结构体，包含与 books 表对应的字段。
编写Go代码，使用Sqlx执行一个复杂的查询，例如查询价格大于 50 元的书籍，并将结果映射到 Book 结构体切片中，确保类型安全。
*/
type Book struct {
	ID     int     `db:"id"`
	Title  string  `db:"title"`
	Author string  `db:"author"`
	Price  float64 `db:"price"`
}

func QueryBooksByPrice(price float64) ([]Book, error) {
	db, err := sqlx.Open("sqlite", constant.DBPath) // 注意驱动名是 sqlite
	if err != nil {
		panic("连接数据库失败: " + err.Error())
	}
	defer db.Close()

	//创建books表
	db.MustExec(`create table if not exists books(id INTEGER PRIMARY KEY AUTOINCREMENT, title TEXT, author TEXT, price REAL);`)

	db.MustExec("INSERT INTO books(title, author, price) VALUES (?, ?, ?)", "书名1", "作者1", 60)
	db.MustExec("INSERT INTO books(title, author, price) VALUES (?, ?, ?)", "书名2", "作者2", 70)
	db.MustExec("INSERT INTO books(title, author, price) VALUES (?, ?, ?)", "书名3", "作者3", 80)
	var books []Book
	err2 := db.Select(&books, "select * from books where price > ? ", price)
	if err2 != nil {
		return nil, err2
	}
	return books, nil
}

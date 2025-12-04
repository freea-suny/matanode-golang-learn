/*
 * @Description:
 * @version: 1.0.0
 * @Author: sun.yong
 * @Date: 2025-12-04 09:21:38
 * @LastEditors: sun.yong
 * @LastEditTime: 2025-12-04 09:53:47
 */
package main

import (
	"fmt"
)

/**
题目 ：定义一个 Shape 接口，包含 Area() 和 Perimeter() 两个方法。然后创建 Rectangle 和 Circle 结构体，实现 Shape 接口。
在主函数中，创建这两个结构体的实例，并调用它们的 Area() 和 Perimeter() 方法。
考察点 ：接口的定义与实现、面向对象编程风格。
*/

type Shape interface {
	Area() float64      //面积
	Perimeter() float64 //周长
}

type Rectangle struct {
	length float64
	width  float64
}

type circle struct {
	radius float64 //半径
}

func (r *Rectangle) Area() float64 {
	return r.length * r.width
}

func (c *circle) Area() float64 {
	return 3.14 * c.radius * c.radius
}

func (r *Rectangle) Perimeter() float64 {
	return 2 * (r.length + r.width)
}

func (c *circle) Perimeter() float64 {
	return 2 * 3.14 * c.radius
}

// func main() {
// 	var shape Shape
// 	//创建矩形实例
// 	shape = &Rectangle{length: 5, width: 3}
// 	//创建圆形实例
// 	shape = &circle{radius: 4}

// 	//调用矩形和圆形的方法
// 	fmt.Printf("矩形面积：%.2f，周长：%.2f\n", shape.Area(), shape.Perimeter())
// 	fmt.Printf("圆形面积：%.2f，周长：%.2f\n", shape.Area(), shape.Perimeter())

// }

/**
题目 ：使用组合的方式创建一个 Person 结构体，包含 Name 和 Age 字段，再创建一个 Employee 结构体，组合 Person 结构体并添加 EmployeeID 字段。
为 Employee 结构体实现一个 PrintInfo() 方法，输出员工的信息。
考察点 ：组合的使用、方法接收者。
*/

type Person struct {
	Name string
	Age  int
}

type Employee struct {
	Person
	EmployeeID int
}

func (e *Employee) PrintInfo(a int) {
	fmt.Printf("姓名：%s，年龄：%d，员工ID：%d\n", e.Name, e.Age, e.EmployeeID)
	e.Person.Age = a
}

// func main() {
// 	emp := &Employee{
// 		Person: Person{
// 			Name: "张三",
// 			Age:  30,
// 		},
// 		EmployeeID: 1001,
// 	}

// 	emp.PrintInfo(18)
// 	fmt.Printf("改后姓名：%s，年龄：%d，员工ID：%d\n", emp.Name, emp.Age, emp.EmployeeID)
// }

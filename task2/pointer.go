/*
 * @Description:
 * @version: 1.0.0
 * @Author: sun.yong
 * @Date: 2025-12-03 16:41:27
 * @LastEditors: sun.yong
 * @LastEditTime: 2025-12-03 21:41:52
 */
package main

import (
	"fmt"
)

/**
题目 ：编写一个Go程序，定义一个函数，该函数接收一个整数指针作为参数，在函数内部将该指针指向的值增加10，
然后在主函数中调用该函数并输出修改后的值。
考察点 ：指针的使用、值传递与引用传递的区别。
*/

// 定义一个结构体
type Student struct {
	Name string
	age  int
}

// 定义一个值传递接收者
func (s Student) setName(name string) {
	s.Name = name
}

// 定义一个指针传递接收者
func (s *Student) setAge(age int) {
	s.age = age
}

func updateStudent() {
	//创建一个Student实例
	student := Student{Name: "suny", age: 18}
	//定义一个指针
	p := &student.age
	student.setName("zhangsan")
	student.setAge(*p + 10)

	//输出修改后的结构体
	fmt.Println(student)
}

/**
题目 ：实现一个函数，接收一个整数切片的指针，将切片中的每个元素乘以2。
考察点 ：指针运算、切片操作。
*/

func testSlice(sl *[]int) {
	for i := 0; i < len(*sl); i++ {
		(*sl)[i] = (*sl)[i] * 2
	}
	fmt.Println("*2后：", *sl)
}

// func main() {
// 	updateStudent()
// 	testSlice(&([]int{1, 2, 3, 4, 5}))
// }

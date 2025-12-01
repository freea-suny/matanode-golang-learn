/*
 * @Description: 判断一个整数是否是回文数
 * @version: 1.0.0
 * @Author: sy
 * @Date: 2025-11-30 23:55:03
<<<<<<< HEAD
 * @LastEditors: sun.yong
 * @LastEditTime: 2025-12-01 00:22:58
=======
 * @LastEditors: sy
 * @LastEditTime: 2025-12-01 00:20:42
>>>>>>> 7e0cdddca306708b335ee28e22a03e2c0ca93c55
 */
package main

import (
	"fmt"
	"strconv"
)

// 是否回文数（从左往右读和从右往左读是一样的整数，如：121，1001等等）
func isPalindrome1(x int) bool {
	//如果x是负数，则不是回文数
	if x < 0 {
		return false
	}

	//将整数转换为字符串
	strX := strconv.Itoa(x)

	//定义两个指针，分别指向字符串的开头和结尾
	head := 0
	tail := len(strX) - 1

	//循环比较字符串的头尾字符，直到指针相遇
	for head < tail {
		h := strX[head]
		t := strX[tail]
		fmt.Println(string(h))
		fmt.Println(t)
		if strX[head] != strX[tail] {
			return false
		}
		head++
		tail--
	}
	return true
}

// func main() {
// 	fmt.Println(isPalindrome1(121))
// }

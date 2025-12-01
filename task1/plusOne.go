/*
 * @Description:
 * @version: 1.0.0
 * @Author: sun.yong
 * @Date: 2025-12-01 22:15:45
 * @LastEditors: sun.yong
 * @LastEditTime: 2025-12-01 22:25:58
 */
package main

/**
给定一个表示 大整数 的整数数组 digits，其中 digits[i] 是整数的第 i 位数字。这些数字按从左到右，从最高位到最低位排列。这个大整数不包含任何前导 0。

将大整数加 1，并返回结果的数字数组。
*/

import (
	"strconv"
)

func plusOne(digits []int) []int {

	//将切片转换成字符串
	str := ""
	for _, v := range digits {
		str += strconv.Itoa(v)
	}
	//将字符串转换成整数，然后+1
	num, _ := strconv.Atoi(str)
	//将整数传换成字符串
	num = num + 1
	//将字符串转换成切片返回plusOne
	result := make([]int, 0)
	for _, v := range strconv.Itoa(num) {
		digit, _ := strconv.Atoi(string(v))
		result = append(result, digit)
	}

	return result

}

// func main() {
// 	fmt.Println(plusOne([]int{9, 9, 9}))
// }

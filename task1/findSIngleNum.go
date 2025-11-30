package main

import (
	"fmt"
)

func main() {
	nums := [...]int{4, 4, 3, 2, 2, 1, 1}
	singleNum := singleNumber(nums)
	fmt.Println("数组中，出现次数为1个的元素为: ", singleNum)
}

func singleNumber(nums [7]int) int {
	if len(nums) == 0 {
		fmt.Println("输入的数组为空")
		return 0
	}

	//定义一个map集合，用来存储数组中的元素，key为元素值，value为元素的次数
	numMap := make(map[int]int)

	//遍历数组
	for _, num := range nums {
		numMap[num] += 1
	}

	//遍历map集合，找到vlaue值为1的key,即为只出现一次的元素
	var singleNum int

	for key, value := range numMap {
		if value == 1 {
			singleNum = key
			break
		}
	}
	return singleNum
}

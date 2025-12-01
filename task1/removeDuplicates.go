/*
 * @Description:
 * @version: 1.0.0
 * @Author: sun.yong
 * @Date: 2025-12-01 22:34:23
 * @LastEditors: sun.yong
 * @LastEditTime: 2025-12-01 22:55:30
 */
package main

import (
	"fmt"
)

/**给你一个 非严格递增排列 的数组 nums ，请你 原地 删除重复出现的元素，使每个元素 只出现一次 ，返回删除后数组的新长度。元素的 相对顺序 应该保持 一致 。然后返回 nums 中唯一元素的个数。

考虑 nums 的唯一元素的数量为 k。去重后，返回唯一元素的数量 k。

nums 的前 k 个元素应包含 排序后 的唯一数字。下标 k - 1 之后的剩余元素可以忽略。*/

func removeDuplicates(nums []int) int {
	if len(nums) < 1 {
		return 0
	}

	//定义一个map保存已经出现的元素
	maps := make(map[int]bool)
	newNUms := make([]int, 0)
	for i, v := range nums {
		fmt.Println("i:", i)
		if maps[v] {
			//如果已经出现过，则删除当前元素
			// newNUms = append(newNUms, nums[i+1:len(nums)-1]...)
			continue
		} else {
			//如果没有出现过，则添加到map中
			maps[v] = true
			newNUms = append(newNUms, v)
		}
	}
	fmt.Println("nums:", newNUms)
	return len(newNUms)
}

// func main() {
// 	//测试用例
// 	nums := []int{1, 1, 2, 3, 1, 4, 5, 5}
// 	fmt.Println(removeDuplicates(nums))
// }

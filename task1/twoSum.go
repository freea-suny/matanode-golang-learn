/*
 * @Description:
 * @version: 1.0.0
 * @Author: sun.yong
 * @Date: 2025-12-01 23:16:51
 * @LastEditors: sun.yong
 * @LastEditTime: 2025-12-01 23:25:25
 */
package main

/**
给定一个整数数组 nums 和一个整数目标值 target，请你在该数组中找出 和为目标值 target  的那 两个 整数，并返回它们的数组下标。

你可以假设每种输入只会对应一个答案，并且你不能使用两次相同的元素。

你可以按任意顺序返回答案。
*/

func twoSum(nums []int, target int) []int {
	//for循环切片，每个元素都与之相加尝试，知道最后一个
	for i := 0; i < len(nums); i++ {
		for j := i + 1; j < len(nums); j++ {
			if nums[i]+nums[j] == target {
				return []int{i, j}
			}
		}
	}
	return []int{}
}

// func main() {
// 	fmt.Println(twoSum([]int{2, 7, 11, 15}, 22))
// }

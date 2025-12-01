/*
 * @Description:
 * @version: 1.0.0
 * @Author: sun.yong
 * @Date: 2025-12-01 22:58:46
 * @LastEditors: sun.yong
 * @LastEditTime: 2025-12-01 23:18:23
 */
package main

/**
以数组 intervals 表示若干个区间的集合，其中单个区间为 intervals[i] = [starti, endi] 。
请你合并所有重叠的区间，并返回 一个不重叠的区间数组，该数组需恰好覆盖输入中的所有区间 。
*/

func merge(intervals [][]int) [][]int {
	//定义一个二维数组
	newNums := make([][]int, 0)
	//循环二维数组
	j := 0
	for i, v := range intervals {
		if i < j {

			continue
		}
		if i == len(intervals)-1 {
			newNums = append(newNums, v)
			break
		}
		if v[0] <= intervals[i+1][0] && v[1] <= intervals[i+1][1] && v[1] >= intervals[i+1][0] {
			newNums = append(newNums, []int{v[0], intervals[i+1][1]})
			j = i + 2
		} else {
			newNums = append(newNums, v)
			j = i + 1
		}

	}

	return newNums
}

// func main() {
// 	fmt.Println(merge([][]int{{1, 3}, {2, 6}, {8, 15}, {15, 18}}))
// }

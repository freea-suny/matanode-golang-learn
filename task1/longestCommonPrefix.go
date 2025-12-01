/*
 * @Description:
 * @version: 1.0.0
 * @Author: sun.yong
 * @Date: 2025-12-01 21:41:38
 * @LastEditors: sun.yong
 * @LastEditTime: 2025-12-01 22:12:20
 */
package main

/**
编写一个函数来查找字符串数组中的最长公共前缀。

如果不存在公共前缀，返回空字符串 ""。

*/

func longestCommonPrefix(strs []string) string {

	if len(strs) < 1 {
		return ""
	}

	//取切片的第一个元素，并且拿第一个元素的第一个字节去和其它元素的第一个字节进行比较，如果都一样，继续，否则结束
	pre := strs[0]
	prefix := ""
outter:
	for i := 0; i < len(pre); i++ {
		p := pre[i]
		for j := 1; j < len(strs); j++ {
			if len(strs[j]) < i+1 {
				break outter
			}
			if p != strs[j][i] {
				break outter
			}
			if j == len(strs)-1 {
				//最后一个字符串
				prefix += string(p)
			}
		}
	}
	return prefix
}

// func main() {
// 	fmt.Println(longestCommonPrefix([]string{"flower", "flow", "flowht"}))
// }

/*
 * @Description:
 * @version: 1.0.0
 * @Author: sun.yong
 * @Date: 2025-12-01 20:44:49
 * @LastEditors: sun.yong
 * @LastEditTime: 2025-12-01 21:34:45
 */
package main

/**
给定一个只包括 '('，')'，'{'，'}'，'['，']' 的字符串 s ，判断字符串是否有效。

有效字符串需满足：

左括号必须用相同类型的右括号闭合。
左括号必须以正确的顺序闭合。
每个右括号都有一个对应的相同类型的左括号。

*/

func isValid(s string) bool {
	if len([]rune(s)) < 2 {
		return false
	}

	//字符串的个数必须是偶数isValid才有效
	if len([]rune(s))%2 != 0 {
		return false
	}

	//定义一个切片
	stack := []string{}
	//定义一个map
	matchingBracket := map[string]string{
		")": "(",
		"}": "{",
		"]": "[",
	}
	//循环字符串中
	for _, v := range []rune(s) {
		switch v {
		case '(', '{', '[':
			stack = append(stack, string(v))
		case ')', '}', ']':
			if len(stack) > 0 && stack[len(stack)-1] == matchingBracket[string(v)] {
				stack = stack[:len(stack)-1]
			} else {
				return false
				//Output:

			}
		default:
			return false
		}
	}

	//如果栈为空，则说明所有的括号都被匹配了
	if len(stack) == 0 {
		return true
	}
	//如果栈不为空，则说明有括号没有被匹配
	return false
}

// func main() {
// 	s := "((({}))){}{[{]"
// 	fmt.Println(isValid(s))
// }

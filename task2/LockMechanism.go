/*
 * @Description:
 * @version: 1.0.0
 * @Author: sun.yong
 * @Date: 2025-12-04 11:00:33
 * @LastEditors: sun.yong
 * @LastEditTime: 2025-12-04 11:19:17
 */
package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

/**题目 ：编写一个程序，使用 sync.Mutex 来保护一个共享的计数器。启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。
考察点 ： sync.Mutex 的使用、并发数据安全。
*/

func inrement(i *int) {

	for j := 0; j < 1000; j++ {
		*i++
	}
}

func sum() {
	counter := 0
	wg := sync.WaitGroup{}
	lock := sync.Mutex{}
	//启动10个协程
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			lock.Lock()
			defer lock.Unlock()
			defer wg.Done()
			inrement(&counter)
		}()
	}
	wg.Wait()
	fmt.Println(counter)
}

// func main() {
// 	for i := 1; i < 100; i++ {
// 		sum()
// 	}

// }

/*
*
题目 ：使用原子操作（ sync/atomic 包）实现一个无锁的计数器。启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。
考察点 ：原子操作、并发数据安全。
*/
func inrement1(i *int32) {

	for j := 0; j < 1000; j++ {
		atomic.AddInt32(i, 1)
	}
}

func sum1() {
	var counter int32 = 0
	wg := sync.WaitGroup{}
	// lock := sync.Mutex{}
	//启动10个协程
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			// lock.Lock()
			// defer lock.Unlock()
			defer wg.Done()
			inrement1(&counter)
		}()
	}
	wg.Wait()
	fmt.Println(counter)
}

// func main() {
// 	for i := 1; i < 100; i++ {
// 		sum1()
// 	}

// }

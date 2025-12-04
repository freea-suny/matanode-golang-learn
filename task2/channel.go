/*
 * @Description:
 * @version: 1.0.0
 * @Author: sun.yong
 * @Date: 2025-12-04 09:54:16
 * @LastEditors: sun.yong
 * @LastEditTime: 2025-12-04 11:00:20
 */
package main

import (
	"fmt"
	"sync"
	"time"
)

/*
*
题目 ：编写一个程序，使用通道实现两个协程之间的通信。一个协程生成从1到10的整数，并将这些整数发送到通道中，另一个协程从通道中接收这些整数并打印出来。
考察点 ：通道的基本使用、协程间通信。
*/
func producer(ch chan<- int, done chan<- bool) {
	for i := 1; i <= 10; i++ {
		ch <- i
	}
	close(ch)
	//消息发送者可以使用通道来传递是否发送成功，主携程可以接受这个信号，判断是否继续下面的代码
	done <- true
}

func consumer(ch <-chan int) {
	for v := range ch {
		fmt.Printf("接收到的整数：%d\n", v)
	}
}

func producer2(ch chan<- int, wg *sync.WaitGroup) {

	//执行结束，销毁
	defer wg.Done()
	for i := 1; i <= 10; i++ {
		ch <- i
	}
	close(ch)

}

func consumer2(ch <-chan int, wg *sync.WaitGroup) {

	//执行结束，销毁
	defer wg.Done()
	for v := range ch {
		fmt.Printf("接收到的整数：%d\n", v)
	}
}

// func main() {
// 	// 	ch := make(chan int, 3)
// 	// 	done := make(chan bool)
// 	// 	go producer(ch, done)
// 	// 	go consumer(ch)

// 	// 	<-done
// 	// 	fmt.Println("执行完毕")

// 	// 	//等待协程完成
// 	// 	// timeOut := time.After(10 * time.Second)
// 	// 	// for {
// 	// 	// 	select {
// 	// 	// 	case _, ok := <-ch:
// 	// 	// 		if !ok {
// 	// 	// 			fmt.Println("通道已关闭")
// 	// 	// 			return
// 	// 	// 		}
// 	// 	// 		fmt.Printf("主线程接收：%d\n", v)
// 	// 	// 	case <-timeOut:
// 	// 	// 		fmt.Println("超时")
// 	// 	// 		return
// 	// 	// 	default:
// 	// 	// 		time.Sleep(100 * time.Millisecond)
// 	// 	// 		fmt.Println("等待数据中...")
// 	// 	// 	}
// 	// 	// }

// 	//方式2，优雅关闭
// 	ch2 := make(chan int, 3)
// 	//创建一个等待组，用于等待所有携程完成
// 	wg2 := sync.WaitGroup{}

// 	//启动生产者协程
// 	//将当前携程放入增加计数
// 	wg2.Add(1)
// 	go producer2(ch2, &wg2)
// 	//启动消费者协程
// 	//将当前携程放入增加计数，注意必须要在启动携程之前增加计数，否则会导致等待组等待永远不会完成
// 	wg2.Add(1)
// 	go consumer2(ch2, &wg2)
// 	//等待所有携程完成
// 	wg2.Wait()
// 	fmt.Println("所有任务执行完毕")

// }

/**
题目 ：实现一个带有缓冲的通道，生产者协程向通道中发送100个整数，消费者协程从通道中接收这些整数并打印。
考察点 ：通道的缓冲机制。
*/

func producer3(ch chan<- int, wg *sync.WaitGroup) {
	//执行结束，销毁
	defer wg.Done()
	for i := 1; i <= 100; i++ {
		time.Sleep(100 * time.Microsecond)
		ch <- i
		// 添加这行，观察缓冲区的填充情况
		fmt.Printf("Sent: %d, Channel Status: [%d/%d]\n", i, len(ch), cap(ch))
	}

	close(ch)
}

func consumer3(ch <-chan int, wg *sync.WaitGroup) {
	//执行结束，销毁
	defer wg.Done()
	for v := range ch {
		time.Sleep(200 * time.Microsecond)
		fmt.Printf("接收到的整数：%d\n", v)
	}
}

// func main() {
// 	//创建一个通道，缓冲区大小为3
// 	ch3 := make(chan int, 5)
// 	//创建一个等待组，用于等待所有携程完成
// 	wg3 := sync.WaitGroup{}

// 	//启动生产者协程
// 	//将当前携程放入增加计数
// 	wg3.Add(1)
// 	go producer3(ch3, &wg3)
// 	//启动消费者协程
// 	//将当前携程放入增加计数，注意必须要在启动携程之前增加计数，否则会导致等待组等待永远不会完成
// 	wg3.Add(1)
// 	go consumer3(ch3, &wg3)
// 	//等待所有携程完成
// 	wg3.Wait()
// 	fmt.Println("所有任务执行完毕")
// }

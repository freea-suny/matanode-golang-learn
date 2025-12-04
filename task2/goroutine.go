/*
 * @Description:
 * @version: 1.0.0
 * @Author: sun.yong
 * @Date: 2025-12-03 17:08:31
 * @LastEditors: sun.yong
 * @LastEditTime: 2025-12-04 00:00:30
 */
package main

import (
	"fmt"
	"sync"
	"time"
)

/**
题目 ：编写一个程序，使用 go 关键字启动两个协程，一个协程打印从1到10的奇数，另一个协程打印从2到10的偶数。
考察点 ： go 关键字的使用、协程的并发执行。
*/

func printOdd() {
	for i := 1; i < 10; i += 2 {
		fmt.Println("奇数：", i)
		time.Sleep(time.Millisecond * 100)
	}
}

func printEven() {
	for i := 2; i <= 10; i += 2 {
		fmt.Println("偶数：", i)
		time.Sleep(time.Millisecond * 100)
	}
}

/*
*
题目 ：设计一个任务调度器，接收一组任务（可以用函数表示），并使用协程并发执行这些任务，同时统计每个任务的执行时间。
考察点 ：协程原理、并发任务调度。
*/

// 定义一个任务结构体
type task struct {
	name     string
	job      func()
	interval time.Duration
	nextRun  time.Time
}

// 定义一个调度器结构体，管理任务
type scheduler struct {
	tasks []*task
	mu    sync.Mutex
	wg    sync.WaitGroup
	Stop  chan bool //用于接收停止信号
}

// 向调度器添加任务 接收者
func (s *scheduler) addTask(t *task) {
	s.tasks = append(s.tasks, t)
	fmt.Println("任务%s已添加，将每隔%v执行一次\n", t.name, t.interval)
}

// 启动调度器
func (s *scheduler) start() {
	fmt.Println("调度器开始工作...")
	//定义一个定时器，每隔1秒检查一次任务列表
	ticker := time.NewTicker(1 * time.Second)
	//确保方法执行完后ticker停止
	defer ticker.Stop()

	//开始循环监听多路复用
	for {
		select {
		case <-s.Stop: //监听任务停止通道信息
			fmt.Println("收到停止信号，调度器正在停止...")
			s.wg.Wait() //等待所有正在执行的任务完成
			return
		case t := <-ticker.C: //监听定时器通道信息
			now := t
			//加锁
			s.mu.Lock()
			defer s.mu.Unlock()
			//遍历任务列表，确认是否有任务需要执行
			for _, ta := range s.tasks {
				if now.After(ta.nextRun) {
					//到达执行时间，启动一个协程执行任务
					s.wg.Add(1) //增加等待的计数
					//启动协程执行任务，不阻塞当前循环
					go func(t *task) {
						defer s.wg.Done() //任务执行完毕，减少计数器
						fmt.Println("[%s]开始执行%s 任务", time.Now().Format("23:26:10"), t.name)
						t.job() //执行真正的任务逻辑
						fmt.Println("[%s]任务 %s 执行完毕", time.Now().Format("12:12:12"), t.name)
					}(ta)
					//更新下一次任务执行时间
					ta.nextRun = now.Add(ta.interval)

				}
			}
		}
	}
}

// 结束调度器任务
func (s *scheduler) stop() {
	close(s.Stop) //关闭channel，通知start函数中的select case
}

func main() {
	//创建一个调度器实例
	scheduler := &scheduler{

		Stop: make(chan bool),
	}

	//添加示例任务
	scheduler.addTask(&task{
		name:     "心跳监测",
		interval: 1 * time.Second,
		job: func() {
			fmt.Println("心跳监测任务执行中...")
			time.Sleep(500 * time.Microsecond)
		},
	})

	scheduler.addTask(&task{
		name:     "数据同步",
		interval: 2 * time.Second,
		job: func() {
			fmt.Println("数据同步任务执行中...")
			time.Sleep(500 * time.Microsecond)
		},
	})

	//在后台启动调度器
	go scheduler.start()

	//主程序等待20s,让调度器运行一短时时间
	time.Sleep(20 * time.Second)

	//优雅地停止调度器
	scheduler.stop()
	fmt.Println("程序退出。")
}

package plus

import (
	"fmt"
	"sync"
	"time"
)

/**
 * go 中的异步操作总结
 *
 * @author wujohns
 * @date 18/5/31
 */
// 单独执行的任务
func simple1() {
	// 延迟 1s
	time.Sleep(time.Second)
	fmt.Println("simple")
}

// AsyncOp1 异步操作案例1
func AsyncOp1() {
	go simple1()
	go simple1()
	fmt.Println("Before run goruntine")

	// 人工阻塞 2s 保证 goruntine 中的操作跑完
	time.Sleep(time.Second * 2)
}

// AsyncOp2 异步操作案例2（waitgroup）
func AsyncOp2() {
	var wg sync.WaitGroup

	// 表示在需要在执行两次 wg.Done() 之后，才会让 wg.Wait() 不阻塞执行
	wg.Add(2)
	go func() {
		defer wg.Done()
		simple1()
	}()

	go func() {
		defer wg.Done()
		simple1()
	}()

	fmt.Println("Before run goruntine")
	wg.Wait()
	fmt.Println("After run goruntine")
}

// 单独执行的任务
func simple2(ch chan int) {
	ch <- 1
	ch <- 2
	ch <- 3

	close(ch)
}

// AsyncOp3 信道的使用（非缓冲）
func AsyncOp3() {
	channel := make(chan int)
	go simple2(channel)
	for num := range channel {
		fmt.Printf("%d", num)
	}
}

// AsyncOp4 信道的使用（缓冲）
func AsyncOp4() {
	channel := make(chan int, 5)
	simple2(channel)
	// go simple2(channel)
	for num := range channel {
		fmt.Printf("%d", num)
	}
}

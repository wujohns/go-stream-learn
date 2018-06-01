package plus

import (
	"fmt"
	"sync"
	"time"
)

/**
 * 对 go 中的 sync 包的使用
 *
 * @author wujohns
 * @date 18/6/1
 */

// WithOutMutex 不使用 Mutex（互斥锁）
func WithOutMutex() {
	num := 0
	for i := 0; i < 200; i++ {
		go func(idx int) {
			num++
		}(i)
	}
	time.Sleep(time.Second)
	fmt.Printf("num=%d\n", num)
}

// WithMutex 使用 Mutex（互斥锁）
func WithMutex() {
	var lock sync.Mutex
	num := 0
	for i := 0; i < 2000; i++ {
		go func(idx int) {
			lock.Lock()
			num++
			lock.Unlock()
		}(i)
	}
	time.Sleep(time.Second)
	fmt.Printf("num=%d\n", num)
}

// OnlyOnce 确保多次并发中某一个操作只被执行一次
func OnlyOnce() {
	var once sync.Once
	for i := 0; i < 10; i++ {
		go func() {
			once.Do(func() { fmt.Println("only once") })
		}()
	}
	time.Sleep(time.Second)
}

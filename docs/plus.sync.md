# go 中的 sync 包的使用
在 io.Pipe 的源码中有对 sync 包的使用，这里仅对在 io.Pipe 用到的特性部分做记录，更多关于 sync 包的使用可以参考 `详细参考` 中的链接。

## sync.Mutex
sync.Mutex 是互斥锁，可以保证在多个 goruntine 在同一时间只有一个在执行。放置并发操作同一变量时造成的一些问题。

不使用互斥锁：

```go
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
```

这里直接创建多个 `goruntine` 对 `num` 进行操作，逐渐增加创建的 `goruntine` 的个数，可以发现当 `goruntine` 超过一定个数时，最终获得的 `num` 的值并会小于 `goruntine` 的个数。  
原因就是在并发对 `num` 进行操作时，部分 `goruntine` 同时取到了 `num` 的值也相同（比如3），计算后将计算后的值写入到 `num` 时写入的值也相同，这样就导致了部分 `num` 最终值小于 `goruntine` 的个数 

使用互斥锁：

```go
// WithMutex 使用 Mutex（互斥锁）
func WithMutex() {
	var lock sync.Mutex
	num := 0
	for i := 0; i < 200; i++ {
		go func() {
			lock.Lock()
			num++
			lock.Unlock()
		}()
	}
	time.Sleep(time.Second)
	fmt.Printf("num=%d\n", num)
}
```

按照上述执行无论建立多少 `goruntine`，最终的得到的 `num` 都等于 `goruntine` 的个数。其机制是当 `lock.Lock()` 后，在 `lock.Unlock()` 之前，其他执行 `lock.Lock()` 的地方都会被阻塞，从而保证同一时间只有一个 `num++` 在执行。

## sync.Once
sync.Once 用于保证在多次执行中（无论是并行还是串行），一个函数只会执行一次，如下：

```go
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

// 结果只输出 only once 一次
```

这里使用 sync.Once 保证并发中的打印操作只执行了一次（sync.Once 是基于 sync.Mutex 封装的，使用了互斥锁的机制保证其只执行一次）

## sync/atomic 中的 atomic.Value
在处理并发读写的安全性上除了采用 `lock` 的方式（锁的部分即上述的 sync.Mutex 以及未提及的 sync.RWMutex）,另外一种方式就是使用 atomic.Value

atomic.Value 的策略是将读操作与写操作视为原子操作，避免在并发时同时对同一片数据进行操作而造成冲突，相比于自定义锁，atomic.Value 方便的地方在于它都封装好了。

主要方法为 Load 与 Save，在 io.Pipe 中对 err 的处理上有应用，主要是为了保证对 err 对象读写的安全性。

这里仅简单列出其接口的使用，具体用法可以参考 `详细参考` 中的 `go atomic.Value相关`

```go
var v atomic.Value
var value <Type>	// 需要存储的值，Type 为该值的类型

v.Store(value)		// 存储变量
v.Load().(<Type>)	// 加载变量，Type 为该值的类型
```

## 详细参考
[go sync包相关使用](https://deepzz.com/post/golang-sync-package-usage.html)
[go atomic.Value相关](https://my.oschina.net/u/222608/blog/881263)
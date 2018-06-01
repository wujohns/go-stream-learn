# go 中的异步处理
稍微深入一点 go 后，用 go 开发的大部分时间都不可避免会用到异步操作相关的东西，这里对 go 的异步操作部分做相关整理，便于以后参考使用。

## 概念相关
### goruntine
在 go 中进行异步操作时，可以使用 go 关键字开启一个 goruntine，让相应的操作在单独的 
goruntine 中运行而不阻塞后续的操作，ex：

```go
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
```

在 `AsyncOp1` 执行操作中，两次将 simple1 放置在 goruntine 中执行，虽然 simple1 中有阻塞的 1s，但是不会阻塞后续的执行。所以输出为：

```
Before run goruntine
simple
simple
```

## 使用 sync.WatiGroup 控制异步流程
在一些场景中我们需要等待异步中的操作执行完毕后进行后续操作，为了显示等待异步执行的任务完成，可以考虑使用 go 自带的 sync.WatiGroup：

```go
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
```

在 `AsyncOp2` 执行操作中，我们使用 `sync.WaitGroup` 创建了 wg 对象，并使用 `wg.Add(2)` 配置让只有在 `wg.Done()` 执行两次之后才会让 `wg.Wait()` 解除阻塞，从而实现了只有两次 `simple1` 执行完后才会让在 `wg.Wait()` 后面继续执行的操作。

`AsyncOp2` 运行结果如下：

```
Before run goruntine
simple
simple
After run goruntine
```

## channel 的使用
channel 是 go 的内置类型，通过 channel 可以发送或接受数据，可用于多个并发 goruntine 中的通讯。

由于网上关于 go channel 的讲解有很多，且都比较详细，这里不做展开讲解。详细部分可以参考 `详细参考` 部分的链接。其中对 channel 的使用以及重点的阻塞机制等有详细说明。

## 代码参考
[async.go](/src/plus/async.go)

## 详细参考
[go channel 相关](http://colobu.com/2016/04/14/Golang-Channels/)  
[go 并发相关](https://liushuchun.gitbooks.io/golang/content/go_concurence.html)  
备注：深入了解 channel，会发现可以用 channel 自己写一个简易版的 sync.WaitGroup
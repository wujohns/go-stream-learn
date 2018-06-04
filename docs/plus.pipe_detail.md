# go 中的 io.Pipe 源码解析
作为补充内容，这里将依据源码讲解 go 中的 pipe 的实现原理。

在开始这部分前，需要对 go 中的异步策略以及 `sync` 包的使用有一定了解。即建议先阅读： 

[go 中的异步处理](/docs/plus.async.md)  
[go 中的sync包使用](/docs/plus.sync.md)

## io.Pipe 原理概述
io.Pipe 是 go 自带的一个方法，使用 io.Pipe 可以构建一对相通的 `reader` 与 `writer`，其相通的原理正是使用了信道（channel）来做的。

该 `reader` 与 `writer` 的结构体中共用了同一个 `wrCh` (写信道) 与 `rdCh` (读信道)，当向 `writer` 中写入数据时，写入的数据同时会被传入到 `wrCh` 中，同时 `reader` 端则会从 `wrCh` 中读取数据并将读取的长度传入到 `rdCh` 中，`writer` 会依据 `rdCh` 判定这次的写入是否全部传入到了 `reader`。

## 自制 io.Pipe
参考 [src/plus/pipe_detail.go](/src/plus/pipe_detail.go)

这里写了一个简单版本的 pipe，便于理解其中的机制。

```go
type pipe struct {
	// 读写管理控制相关
	wrMu sync.Mutex
	wrCh chan []byte
	rdCh chan int

	// 读写完成控制相关
	once sync.Once
	done chan struct{}
}

func (p *pipe) Read(b []byte) (n int, err error) {
	select {
	case bw := <-p.wrCh:
		nr := copy(b, bw)
		p.rdCh <- nr
		return nr, nil
	case <-p.done:
		return 0, io.EOF
	}
}

func (p *pipe) Write(b []byte) (n int, err error) {
	select {
	case <-p.done:
		return 0, io.EOF
	default:
		// 如果写没有 close 则加锁后开始后续的写入操作
		p.wrMu.Lock()
		defer p.wrMu.Unlock()
	}

	for once := true; once || len(b) > 0; once = false {
		select {
		// 将数据写入到 wrCh 中
		case p.wrCh <- b:
			// 获取 rdCh，将上次未能写完的数据继续写入
			nw := <-p.rdCh
			b = b[nw:]
			n += nw
		}
	}
	return n, nil
}

// PipeReader 的封装
type PipeReader struct {
	p *pipe
}

func (r *PipeReader) Read(data []byte) (n int, err error) {
	return r.p.Read(data)
}

// PipeWriter 的封装
type PipeWriter struct {
	p *pipe
}

func (w *PipeWriter) Write(data []byte) (n int, err error) {
	return w.p.Write(data)
}

// Close 关闭 writer
func (w *PipeWriter) Close() {
	w.p.once.Do(func() { close(w.p.done) })
}

// CusPipe 自制简易 pipe
func CusPipe() (*PipeReader, *PipeWriter) {
	p := &pipe{
		wrCh: make(chan []byte),
		rdCh: make(chan int),
		done: make(chan struct{}),
	}
	return &PipeReader{p}, &PipeWriter{p}
}
```

在上述代码中，我们构建了一个 `pipe` 的类型，而 `PipeReader` 与 `PipeWriter` 均是依据此拓展得来的，其中:  
1. `Read` 方法会从 `wrCh` 中获取数据并返回此次读取的长度到 `rdCh`  
1. `Write` 方法会向 `wrCh` 中写入数据并通过 `rdCh` 确认所有数据均被 `Read` 所处理  
1. 在处理完数据后执行 `Close` 通过信道 `done` 终止掉 `Read` 与 `Write` 方法里的读写操作

## 对官方的 io.Pipe 的补充说明
这里参照的是 go 1.10.2。官方的版本可以视为增加了更为完善的 close 与 error 机制（或者说上述的自制版本中移除或简化了该部分）。
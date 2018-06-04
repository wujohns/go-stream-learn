package plus

import (
	"exs"
	"io"
	"sync"
)

/**
 * pipe 详细的实验
 */

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
		// 如果写没有 close 则继续
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

// CusPipeTest 测试该 CusPipe
func CusPipeTest() {
	readStream, writeStream := CusPipe()
	srcReader := &exs.BlankReader{
		Count:  30,
		Cursor: 0,
	}
	distWriter := exs.BlankWriter{}

	go func() {
		defer writeStream.Close()
		exs.OpA(srcReader, writeStream)
	}()
	exs.OpB(readStream, distWriter)
}

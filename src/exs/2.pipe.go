package exs

import (
	"io"
)

/**
 * go 中的 pipe 的应用
 *
 * @author wujohns
 * @date 18/5/31
 */

// OpA 对数据的处理A
func OpA(r io.Reader, w io.Writer) {
	buf := make([]byte, 10)
	wBuf := make([]byte, 10)

	// 读取并处理数据直到数据全部读完
	for {
		// 从 reader 中读取数据
		len, err := r.Read(buf)

		// 这里只是展示对数据的加工处理方式，实际是统一写入了同一字符
		if len > 0 {
			for n := 0; n < len; n++ {
				wBuf[n] = byte(100)
			}
			w.Write(wBuf)
		}
		if err != nil {
			break
		}
	}
}

// OpB 对数据的处理B
func OpB(r io.Reader, w io.Writer) {
	buf := make([]byte, 10)
	io.CopyBuffer(w, r, buf)
}

// PipeTest pipe 的使用实验
func PipeTest() {
	readStream, writeStream := io.Pipe()
	srcReader := &BlankReader{30, 0}
	distWriter := BlankWriter{}

	// 将 OpA 放在 goroutine 中的原因是如果直接执行，会使得 writerStream 被写满而阻塞
	go func() {
		defer writeStream.Close()
		OpA(srcReader, writeStream)
	}()
	OpB(readStream, distWriter)
}

// StdOpA 对 OpA 的封装
func StdOpA(r io.Reader) io.Reader {
	readStream, writeStream := io.Pipe()
	go func() {
		defer writeStream.Close()
		OpA(r, writeStream)
	}()
	return readStream
}

// StdPipeTest 标准化尝试
func StdPipeTest() {
	srcReader := &BlankReader{30, 0}
	distWriter := BlankWriter{}

	r := StdOpA(srcReader)
	OpB(r, distWriter)
}

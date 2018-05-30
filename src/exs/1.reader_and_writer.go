package exs

import (
	"io"
	"math/rand"
	"os"
	"time"
)

// BlankReader 自定义 reader
type BlankReader struct {
	count  int // 总量
	cursor int // 当前位置
}

// 对 Reader 接口中的 Read 方法的实现
func (r *BlankReader) Read(p []byte) (n int, err error) {
	rand.Seed(time.Now().Unix())

	// 在到达总量（count）且 p 未被填满的情况下持续将数据写入到 p
	for ; r.cursor < r.count && n < len(p); r.cursor++ {
		// 这里生成随机字符并写入到切片 p 中
		num := rand.Intn(57) + 65
		for {
			if num > 90 && num < 97 {
				num = rand.Intn(57) + 65
			} else {
				break
			}
		}
		p[n] = byte(num)
		n++
	}
	if n == 0 {
		// n 为 0，在上一次读取中达到了总量，本次读取返回结束标记
		return 0, io.EOF
	}
	// 返回读取的数据量（n）
	return n, nil
}

// WriteToFileTest 构建一个 blankReader，并将数据写入一个文件中
func WriteToFileTest() {
	r := &BlankReader{100, 0}
	w, _ := os.Create("files/form_blank_reader.txt")
	buf := make([]byte, 10)
	io.CopyBuffer(w, r, buf)
}

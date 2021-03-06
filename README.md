# go-stream-learn
go的流式处理笔记

## 简介
流式处理一般用于大量数据的处理中，采用边读取，边处理，边写入的策略可以将整个步骤使用的内存限定在一定范围内。这里对go中的流式处理策略做相应分析，方便理解其中的细节，以及作为之后使用流式处理时的参考。

备注：由于目前的主力语言是 nodejs，所以在相关的说明中会与 nodejs 的 stream 做类比（即如果对 nodejs 的 stream 有一定了解的话，这里的阅读会轻松很多）。

## 章节
(1) [go中的reader与writer](/docs/1.reader与writer.md)   
(2) [go中的io.Pipe](/docs/2.pipe.md)  

[plus go中的异步处理](/docs/plus.async.md)  
[plus go中的sync包的使用](/docs/plus.sync.md)  
[plus go中io.Pipe源码解析](/docs/plus.pipe_detail.md)

## 其他参考（非本仓库作者笔记）
[liushuchun的go笔记](https://liushuchun.gitbooks.io/golang/content/)
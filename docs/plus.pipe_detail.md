# go 中的 io.Pipe 源码解析
作为补充内容，这里将依据源码讲解 go 中的 pipe 的实现原理。

在开始这部分前，需要对 go 中的异步策略以及 `sync` 包的使用有一定了解。即建议先阅读：  
[go 中的异步处理](/docs/plus.async.md)  
[go 中的sync包使用](/docs/plus.sync.md)

## TODO
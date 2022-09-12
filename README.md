# GeeCache
[7天用Go从零实现分布式缓存GeeCache教程系列](https://geektutu.com/post/geecache-day1.html)

## 文件结构
```text
geecache/
    |--lru/
        |--lru.go   // lru 缓存淘汰策略
    |--byteview.go  // 缓存值的抽象与封装
    |--cache.go     // 并发控制
    |--geecache.go  // 负责与外部交互，控制缓存存储和获取的主流程
```

## Group 主体结构
Group 是 GeeCache 最核心的数据结构，负责与用户的交互，并且控制缓存值存储和获取的流程
```text
                           是
接收 key --> 检查是否被缓存 -----> 返回缓存值 ⑴
                |  否                        是
                |-----> 是否应当从远程节点获取 -----> 与远程节点交互 --> 返回缓存值 ⑵
                            |  否
                            |-----> 调用`回调函数`，获取值并添加到缓存 --> 返回缓存值 ⑶
```


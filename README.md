# GDM
A Multi-Thread Download Manager

使用示例：

```
go run main.go -url="目标url" -threadNum=目标线程数 -localAddress="本地存储目录"
go run main.go -url=https://iterm2.com/downloads/stable/iTerm2-3_3_12.zip -threadNum=5 -address=/Users/username/Downloads
```

**注**：如果不使用threadNum，默认线程数为6；

**目录结构**
1. conf——存放配置文件，限定了失败重复次数、临时文件前缀
2. genesis——初始化相关函数，初始化配置、日志工具等
3. http——下载函数
4. log——自动生成的日志文件保存路径，30天前日志文件会自动删除
5. logs——获取日志相关操作步骤
6. test——相关测试函数，可以删除
7. utils——相关工具函数
# GDM
A Multi-Thread Download Manager

使用示例：

```
go run main.go -url="目标url" -threadNum=目标线程数 -localAddress="本地存储目录"
```

如果不使用threadNum，默认线程数为6；
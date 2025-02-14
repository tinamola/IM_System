# IM_System
An online Instant messaging system

## Reference
This project follows the guide from the following video:  
[8小时转职Golang工程师 (如果你想低成本学习Go语言)](https://www.bilibili.com/video/BV1gf4y1r79E/)

## Getting Started

First, build the project:
```bash
go build -o server server.go main.go user.go
```

Next, run the server:

```bash
./server
```

Meanwhile, open other terminals and log in as different users with:

```bash
nc 127.0.0.1 8888
```


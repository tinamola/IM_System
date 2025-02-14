# IM_System
An online Instant messaging system

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


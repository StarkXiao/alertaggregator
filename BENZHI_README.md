# Structured Alert Aggregator

这是一个 Go 结构化日志告警聚合器，接收日志事件，按服务、环境、级别和归一化消息聚合告警，抑制重复通知，并追踪确认、恢复和通知历史。

## 本地校验

```bash
go test ./...
go test -race ./...
go vet ./...
go build ./...
```

## 启动服务

```bash
go run ./cmd/aggregator -addr :8080 -db ./data/alerts.json
```

## Docker 校验

```bash
./build_benzhi_docker.sh alertaggregator linux/amd64
./build_benzhi_docker.sh alertaggregator linux/arm64
docker run -it alertaggregator:latest
```

根目录是健康的基线模块。`base_bug_001` 至 `base_bug_010` 是从 `main` 导出的独立 Bug 复现源码目录，每个目录可单独执行其对应测试命令。

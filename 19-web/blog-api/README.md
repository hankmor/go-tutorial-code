# 项目架构与综合实战

本目录对应 Web 篇 510《项目架构设计》和 580《Web 开发测试与质量保障》。示例使用 Go 1.25.1、Gin 和 GORM，演示一个小型博客 API 如何把 HTTP 处理、业务规则和数据库访问拆开，并按职责验证各层行为。

## 运行

```bash
go run .
```

服务默认监听 `http://localhost:8086`。可使用 `curl` 访问 `/healthz`、`POST /posts` 和 `GET /posts`。

## 验证

```bash
go test .
go test -race .
go vet .
```

`app_test.go` 包含 Service 业务规则和错误传播测试、GORM Repository 行为测试、HTTP 请求与 JSON 响应契约测试，以及数据库不可用时的失败路径测试。测试使用独立的 SQLite 内存数据库，不要求预先启动服务或监听固定端口。

`model.go`、`repository.go`、`service.go`、`handler.go` 和 `app.go` 分别对应文章中的模型、数据访问、业务、HTTP 和依赖组装职责。580 章练习的参考实现集中在 `app_test.go`，可以通过 `go test . -count=1`、`go test . -race -count=1` 和 `go test . -shuffle=on -count=10` 验证。

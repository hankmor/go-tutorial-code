# 项目架构与综合实战

本目录对应 Web 篇 510《项目架构设计》。示例使用 Go 1.25.1、Gin 和 GORM，演示一个小型博客 API 如何把 HTTP 处理、业务规则和数据库访问拆开。

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

文章中的实质性代码练习是为 Service 增加文章标题长度规则和对应测试；当前实现已包含该规则，测试位于 `app_test.go`。`model.go`、`repository.go`、`service.go`、`handler.go` 和 `app.go` 分别对应文章中的模型、数据访问、业务、HTTP 和依赖组装职责。

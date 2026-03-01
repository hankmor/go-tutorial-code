# Go 模块管理示例代码

本目录包含 Go Modules 相关的示例代码。

## 目录结构

```
12-modules/
├── 01-basic/              # 基础示例
│   ├── go.mod
│   └── main.go
├── 02-dependencies/       # 依赖管理
│   ├── go.mod
│   └── main.go
├── 03-replace/           # replace 指令
│   ├── go.mod
│   ├── main.go
│   └── local-package/
├── 04-versions/          # 版本管理
│   ├── go.mod
│   └── main.go
├── 05-private/           # 私有仓库
│   ├── go.mod
│   └── main.go
└── README.md

## 运行示例

每个示例都是独立的模块，可以单独运行：

```bash
# 进入示例目录
cd 01-basic

# 下载依赖
go mod download

# 运行
go run main.go
```

## 常用命令

```bash
# 初始化模块
go mod init github.com/username/project

# 整理依赖
go mod tidy

# 查看依赖
go list -m all

# 升级依赖
go get -u ./...

# 清理缓存
go clean -modcache
```

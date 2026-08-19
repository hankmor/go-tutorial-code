# Web 篇示例

本目录对应 Go Web 教程，使用 Go 1.25.1。第一篇《Gin 入门与路由》对应 `gin-basic`，后续章节按主题映射到中间件、GORM、上传与 JWT、优雅关闭、可观测性和项目示例目录。

## 文章与代码映射

| 文章主题 | 示例目录 |
| --- | --- |
| Gin 入门与路由 | `gin-basic` |
| Gin 中间件、数据绑定 | `middleware` |
| Gin 错误处理与统一响应 | `error-response` |
| 文件上传与 JWT | `upload-jwt` |
| 跨域与优雅关闭 | `cors`、`graceful` |
| GORM 入门与 CRUD | `gorm-crud` |
| GORM 高级特性 | `gorm-advanced` |
| 工程实践：项目结构 | `../14-project-layout` |
| 工程实践：项目实战结构 | `../15-project-example` |
| 工程实践：环境变量与配置 | `../18-advanced/config` |
| 项目架构与综合实战 | `blog-api` |
| 部署 | `deploy` |
| 监控与可观测性 | `observability` |

```bash
go run ./gin-basic
JWT_SECRET='replace-with-a-long-random-secret' go run ./upload-jwt
```

验证全部 Web 示例：

```bash
go test ./...
go vet ./...
go test ./gin-basic -race
```

`gin-basic` 通过 `NewRouter` 暴露可测试的路由构造函数，测试覆盖基础响应、参数错误、静态路由优先级和路由分组。`middleware` 通过 `NewRouter` 覆盖 JSON 绑定、字段校验、请求 ID、认证拦截和错误响应。`error-response` 覆盖统一响应、业务错误、Panic 恢复和错误信息脱敏。`upload-jwt` 覆盖 JWT 签名算法和过期时间校验、Bearer 认证、上传大小、图片类型和服务端文件名。`cors` 覆盖允许来源、拒绝来源、预检状态码和 CORS 响应头，练习答案还包括多来源与凭据配置。`graceful` 通过 `NewServer` 验证服务启动、请求处理和 `Shutdown` 后的正常返回，练习答案覆盖慢请求超时后继续清理业务资源。`gorm-crud` 覆盖迁移、创建、查询、零值更新、软删除和唯一约束。完整的部署配置见 `deploy` 目录。

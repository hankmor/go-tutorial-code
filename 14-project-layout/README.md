# Go 项目结构示例

本目录包含不同规模的 Go 项目结构示例。

## 目录结构

```
13-project-layout/
├── 01-simple/          # 小型项目（单文件）
├── 02-basic/           # 基础项目（cmd + internal）
├── 03-standard/        # 标准项目（完整布局）
├── 04-multi-entry/     # 多入口项目
├── 05-microservices/   # 微服务项目
└── README.md
```

## 项目规模建议

- **小型项目（<1000 行）**：使用 01-simple 结构
- **中型项目（1000-10000 行）**：使用 02-basic 或 03-standard
- **大型项目（>10000 行）**：使用 03-standard 完整布局
- **微服务**：使用 05-microservices 结构

## 核心目录说明

- **cmd/**：程序入口，每个子目录对应一个可执行文件
- **internal/**：私有代码，只能被本项目导入
- **pkg/**：公共库，可以被外部项目导入（谨慎使用）
- **api/**：API 定义（OpenAPI/Proto）
- **configs/**：配置文件模板
- **scripts/**：构建和部署脚本

## 最佳实践

1. 从简单开始，按需扩展
2. 优先使用 internal/，除非确定要公开
3. main.go 要简洁，只负责启动
4. 配置文件不包含敏感信息
5. 使用 Makefile 管理常用命令

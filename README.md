# Go 语言完整教程 - 配套代码

[![Go Version](https://img.shields.io/badge/Go-1.23-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Tutorial](https://img.shields.io/badge/Tutorial-hankmo.com-orange.svg)](https://hankmo.com)

## 📚 教程结构

本代码库配套《极客老墨 GoLang 避坑实战教程》，分为三个部分：

### 第一部分：基础教程（17 章）

- **目标**：让零基础读者能写出简单的 Go 程序
- **状态**: ✅已完成

- [x] 第 0 章：开篇
- [x] 第 1 章：语言基础（核心语法）
- [x] 第 2 章：并发编程入门
- [x] 第 3 章：标准库与测试
- [x] 第 4 章：工程实践入门
- [x] 第 5 章：动手实战

### 第二部分：高级教程（21 章）

- **目标**：深入底层原理，掌握高级编程技巧
- **状态**: ⏰进行中

- [ ] 第 6 章：并发编程进阶
- [ ] 第 7 章：高级语言特性
- [ ] 第 8 章：内存管理与性能优化
- [ ] 第 9 章：CGO 与底层编程
- [ ] 第 10 章：工程实践进阶
- [ ] 第 11 章：测试进阶
- [ ] 第 12 章：依赖注入

### 第三部分：Web 开发实战（27 章）

- **目标**：从零到一构建企业级 Web 应用
- **状态**: 📅计划中

- [ ] 第 13 章：Gin 框架基础
- [ ] 第 14 章：GORM ORM 框架
- [ ] 第 15 章：企业级项目架构
- [ ] 第 16 章：依赖注入与解耦
- [ ] 第 17 章：测试与质量保障
- [ ] 第 18 章：部署与监控
- [ ] 第 19 章：综合实战项目

## 🚀 快速开始

### 环境要求

- Go 1.24 或更高版本
- Git

### 克隆代码

```bash
git clone https://github.com/hankmo/go-tutorial-code.git
cd go-tutorial-code
```

### 运行示例

每个章节都是独立的 Go 模块，可以单独运行：

```bash
# 进入某个章节目录
cd 01-intro

# 运行示例
go run main.go

# 运行测试
go test -v
```

## 📖 配套教程

完整教程请访问：[https://hankmo.com](https://hankmo.com)

## 🗂️ 目录结构

```
go-tutorial-code/
├── 00-setup/           # 安装
├── 01-intro/           # 简介
├── 02-style/           # 代码风格
├── ...
├── common/             # 公共代码库
├── docs/               # 文档
└── scripts/            # 工具脚本
```

## 💡 代码规范

所有代码遵循以下规范：

- ✅ 使用 Go 1.23 版本
- ✅ 通过 `golangci-lint` 检查
- ✅ 包含单元测试（覆盖率 > 60%）
- ✅ 代码有详细注释
- ✅ 每个示例都有 README 说明

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📄 许可证

MIT License

## 👨‍💻 关于作者

**极客老墨** - 一个热爱折腾的开发者

- 博客：[https://hankmo.com](https://hankmo.com)
- GitHub：[@hankmo](https://github.com/hankmor)
- 公众号：极客老墨

<img src="./docs/weixinqr.jpg" width="200">

---

**极客老墨，继续折腾！** 💪

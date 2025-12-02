# Backend (Go) — 快速启动

说明：这是为 `backend/` 目录初始化的 Go 后端服务，使用 Gin 框架，支持 GORM（PostgreSQL）、Redis、JWT、Viper 配置管理等。

假设与说明
- Module 路径：`github.com/yycy134679/school-secondhand-trading-system/backend`
- 如果你使用不同的仓库地址，请编辑 `go.mod` 中的 module 行。

## 前置要求

1. **安装 Go**
   - 官方下载：https://go.dev/dl/
   - 建议 Go 1.20 或更高（项目使用 Go 1.23）

2. **配置数据库（可选）**
   - PostgreSQL ≥ 13
   - 创建数据库：`CREATE DATABASE secondhand_dev;`
   - 在 `.env` 中配置 `DB_DSN`（见下）

3. **配置 Redis（可选）**
   - 若使用推荐/缓存功能，在 `.env` 中配置 `REDIS_ADDR`

## 快速启动（Windows PowerShell）

```powershell
# 1. 进入 backend 目录
cd d:/goland/gocode/school-secondhand-trading-system/backend

# 2. 下载依赖
go mod tidy

# 3. 运行应用（从 cmd/main.go 启动）
go run ./cmd

# 或构建后运行
go build -o backend.exe ./cmd
.\backend.exe
```

## 验证服务

打开浏览器或在 PowerShell 中测试：

```powershell
# 健康检查
Invoke-WebRequest -UseBasicParsing http://localhost:8080/health

# 示例 API（用户注册/登录等见路由）
Invoke-WebRequest -UseBasicParsing http://localhost:8080/api/v1/users/register
```

在 GoLand 中打开项目
- 打开 GoLand，选择 `Open` 并选择仓库根目录或 `backend` 目录。
- 确保在 Settings/Preferences -> Go -> GOROOT 中已配置本地 Go SDK（安装路径）。
- GoLand 会自动识别 `go.mod` 并启用 Go Modules。若没有，手动在右下角选择 "Enable Go Modules integration"。
- 创建 Run Configuration：Add -> Go Build / Go Run，设置工作目录为 `.../backend`，Main file 指向 `main.go`。

常用工具（可选）
- 格式化：`gofmt` / GoLand 自动格式化
- 静态检查：`golangci-lint`（推荐）
- 调试：使用 Delve，GoLand 已集成调试支持

后续建议
- 如果你计划用某个框架（Gin/Echo/Fiber 等），我可以帮你把 `main.go` 改为框架启动样板并添加示例路由与 Dockerfile。

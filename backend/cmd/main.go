// Package main 是应用程序的入口点
// 负责初始化配置、数据库连接、Redis连接，并启动HTTP服务器
package main

import (
	"fmt"
	"log"

	"github.com/yycy134679/school-secondhand-trading-system/backend/config"
	"github.com/yycy134679/school-secondhand-trading-system/backend/router"
)

// main 函数是程序的启动入口
// 执行流程：
// 1. 加载配置（从.env文件或环境变量）
// 2. 初始化数据库连接（PostgreSQL + GORM）
// 3. 初始化Redis连接（可选，用于缓存和推荐系统）
// 4. 设置路由和中间件
// 5. 启动HTTP服务器
func main() {
	// 步骤1: 加载应用配置
	// LoadConfig 会尝试从以下来源读取配置（优先级从高到低）：
	// - 环境变量
	// - .env 文件（位于backend目录下）
	cfg, err := config.LoadConfig()
	if err != nil {
		// 如果配置加载失败，记录致命错误并退出程序
		log.Fatalf("load config: %v", err)
	}

	// 打印已加载的配置信息（用于调试）
	// 注意：生产环境应避免打印敏感信息（如密码）
	log.Printf("Loaded config: DB_DSN=%s, HTTP_PORT=%d", cfg.DBDSN, cfg.HTTPPort)

	// 步骤2: 初始化数据库连接
	// 使用GORM（Go的ORM库）连接PostgreSQL数据库
	// 如果DSN为空字符串，NewDB会返回nil（允许在没有数据库的情况下运行）
	db, err := config.NewDB(cfg.DBDSN)
	if err != nil {
		// 数据库连接失败时，打印警告但不退出程序
		// 这允许在本地开发时即使没有数据库也能启动服务器
		log.Printf("warning: failed to init DB: %v (continuing without DB)", err)
	} else {
		// 数据库连接成功
		log.Println("DB connection established successfully")
	}

	// 步骤3: 设置路由和中间件
	// SetupRouter 会注册所有HTTP路由和中间件
	// 包括：用户模块、商品模块、分类标签模块等
	// 推荐系统现在使用内存缓存，不再依赖Redis
	r := router.SetupRouter(db, cfg)

	// 步骤4: 启动HTTP服务器
	// 构造监听地址（例如：:8080）
	addr := fmt.Sprintf(":%d", cfg.HTTPPort)
	log.Printf("starting server on %s", addr)

	// 启动Gin HTTP服务器
	// r.Run() 会阻塞，直到服务器关闭或发生错误
	if err := r.Run(addr); err != nil {
		// 服务器启动失败或运行时错误，记录致命错误并退出
		log.Fatalf("server error: %v", err)
	}
}

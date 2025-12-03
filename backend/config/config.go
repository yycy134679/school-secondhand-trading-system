// Package config 负责应用程序的配置管理
// 使用 Viper 库从.env文件和环境变量中读取配置
package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config 结构体保存应用程序的所有配置项
// 这些配置项来自.env文件或系统环境变量
type Config struct {
	AppEnv         string // 应用环境：development/production
	HTTPPort       int    // HTTP服务器监听端口，默认8080
	DBDSN          string // 数据库连接字符串（PostgreSQL）
	RedisAddr      string // Redis服务器地址，格式：host:port
	JWTSecret      string // JWT签名密钥，用于token的生成和验证
	FileStorageDir string // 文件上传存储目录，用于保存商品图片等
}

// LoadConfig 从配置源加载应用配置
// 配置读取优先级：
// 1. 系统环境变量（最高优先级）
// 2. .env 文件
// 3. 代码中定义的默认值（最低优先级）
//
// 返回值：
//   - *Config: 配置对象指针
//   - error: 配置加载或验证失败时返回错误
func LoadConfig() (*Config, error) {
	// 创建一个新的Viper实例用于配置管理
	v := viper.New()

	// 设置配置文件名称和类型
	v.SetConfigName(".env")   // 文件名（不包含扩展名）
	v.SetConfigType("dotenv") // 文件类型：dotenv格式（KEY=value）

	// 添加配置文件搜索路径（按顺序）
	v.AddConfigPath(".")         // 当前目录
	v.AddConfigPath("./backend") // backend子目录

	// 启用自动读取环境变量
	// 环境变量会覆盖.env文件中的配置
	v.AutomaticEnv()
	v.SetEnvPrefix("") // 不使用前缀，直接读取原始环境变量名

	// 尝试读取.env文件
	// 如果文件不存在或读取失败，不会中断程序，而是打印警告
	if err := v.ReadInConfig(); err != nil {
		fmt.Printf("Warning: .env file not found or error reading: %v\n", err)
	} else {
		// 成功读取时，打印配置文件的完整路径
		fmt.Printf("Loaded .env from: %s\n", v.ConfigFileUsed())
	}

	// 设置默认值（当配置项未在文件或环境变量中定义时使用）
	v.SetDefault("APP_ENV", "development")           // 默认开发环境
	v.SetDefault("HTTP_PORT", 8080)                  // 默认端口8080
	v.SetDefault("DB_DSN", "")                       // 默认无数据库连接
	v.SetDefault("REDIS_ADDR", "")                   // 默认无Redis连接
	v.SetDefault("JWT_SECRET", "please-change-this") // 默认JWT密钥（生产环境必须修改）
	v.SetDefault("FILE_STORAGE_DIR", "./uploads")    // 默认文件存储目录

	// 从Viper中读取配置值并构建Config对象
	cfg := &Config{
		AppEnv:         v.GetString("APP_ENV"),
		HTTPPort:       v.GetInt("HTTP_PORT"),
		DBDSN:          v.GetString("DB_DSN"),
		RedisAddr:      v.GetString("REDIS_ADDR"),
		JWTSecret:      v.GetString("JWT_SECRET"),
		FileStorageDir: v.GetString("FILE_STORAGE_DIR"),
	}

	// 配置验证：HTTP端口不能为0
	if cfg.HTTPPort == 0 {
		return nil, fmt.Errorf("invalid HTTP_PORT: 0")
	}

	return cfg, nil
}

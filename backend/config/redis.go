package config

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// NewRedis 初始化并返回一个Redis客户端实例
//
// Redis在本系统中的应用场景：
//  1. 推荐系统缓存：缓存用户推荐商品列表，减少数据库查询
//  2. 商品状态撤销：记录商品状态变更的3秒撤销窗口数据
//  3. 会话管理：存储用户会话信息（可选）
//
// 参数：
//   - addr: Redis服务器地址，格式为 "host:port"，例如 "localhost:6379"
//     如果addr为空字符串，函数返回(nil, nil)，允许应用在无Redis模式下运行
//
// 返回值：
//   - *redis.Client: Redis客户端实例
//   - error: 连接失败时返回错误
//
// 连接配置说明：
//   - DialTimeout: 连接超时时间，设为5秒
//   - Ping超时: 连接测试超时时间，设为3秒
func NewRedis(addr string) (*redis.Client, error) {
	// 如果地址为空，返回nil（允许无Redis模式）
	// 这样在开发环境或Redis不可用时，应用仍然可以启动
	if addr == "" {
		return nil, nil
	}

	// 创建Redis客户端实例
	// redis.Options 包含连接配置参数
	rdb := redis.NewClient(&redis.Options{
		Addr:        addr,            // Redis服务器地址
		DialTimeout: 5 * time.Second, // 连接建立超时时间
		// 其他可选参数（未设置则使用默认值）：
		// Password: "",                 // Redis密码（如果需要认证）
		// DB: 0,                        // 使用的数据库编号（0-15）
		// PoolSize: 10,                 // 连接池大小
		// MinIdleConns: 5,              // 最小空闲连接数
	})

	// 创建一个带超时的上下文，用于执行Ping命令
	// Ping用于验证Redis连接是否真正可用
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel() // 确保上下文资源被释放

	// 执行PING命令测试连接
	// Ping() 返回一个StatusCmd对象，调用Err()获取错误信息
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	// 连接成功，返回Redis客户端实例
	return rdb, nil
}

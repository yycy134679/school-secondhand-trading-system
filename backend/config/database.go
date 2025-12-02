package config

import (
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewDB 初始化并返回一个GORM数据库连接实例
//
// 功能说明：
//   - 使用PostgreSQL驱动连接数据库
//   - 配置连接池参数以优化性能
//   - 执行ping测试验证连接可用性
//
// 参数：
//   - dsn: 数据库连接字符串（DSN），格式示例：
//     "host=localhost user=postgres password=123456 dbname=mydb port=5432 sslmode=disable"
//     如果dsn为空字符串，函数会返回(nil, nil)，允许应用在无数据库模式下运行
//
// 返回值：
//   - *gorm.DB: GORM数据库实例，用于执行ORM操作
//   - error: 连接失败时返回错误
//
// 连接池配置说明：
//   - MaxIdleConns: 最大空闲连接数，设为10
//   - MaxOpenConns: 最大打开连接数，设为100
//   - ConnMaxLifetime: 连接最大生存时间，设为30分钟（防止长时间连接过期）
func NewDB(dsn string) (*gorm.DB, error) {
	// 如果DSN为空，返回nil（允许无数据库模式）
	if dsn == "" {
		return nil, nil
	}

	// 使用GORM打开PostgreSQL数据库连接
	// postgres.Open(dsn) 创建PostgreSQL方言的连接配置
	// &gorm.Config{} 使用默认的GORM配置
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// 获取底层的*sql.DB对象，用于配置连接池
	// GORM基于database/sql，这里获取原生DB对象进行底层配置
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// 配置连接池参数（这些是生产环境的合理默认值）

	// SetMaxIdleConns 设置空闲连接池中的最大连接数
	// 空闲连接可以立即被复用，无需重新建立连接，提升性能
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns 设置数据库的最大打开连接数
	// 限制并发连接数，防止数据库连接数过多导致资源耗尽
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime 设置连接的最大生存时间
	// 超过这个时间的连接会被关闭并重新创建
	// 这可以避免长时间连接因数据库端超时而失效
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	// 执行一次ping测试，验证数据库连接是否真正可用
	// 如果ping失败，会记录警告日志，但不会返回错误（容错处理）
	if err := sqlDB.Ping(); err != nil {
		log.Printf("db ping failed: %v", err)
	}

	return db, nil
}

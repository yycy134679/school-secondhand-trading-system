// Package cache 提供基于内存的缓存服务
// 用于替代Redis，实现推荐系统缓存和商品状态撤销功能
package cache

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// cacheItem 缓存项，包含值和过期时间
type cacheItem struct {
	value      interface{}
	expiration time.Time
}

// MemoryCache 内存缓存服务
// 线程安全的内存缓存实现，支持TTL（过期时间）
type MemoryCache struct {
	data   map[string]*cacheItem
	mu     sync.RWMutex
	stopCh chan struct{}
}

// NewMemoryCache 创建新的内存缓存实例
// 自动启动后台清理goroutine，每分钟清理一次过期数据
func NewMemoryCache() *MemoryCache {
	mc := &MemoryCache{
		data:   make(map[string]*cacheItem),
		stopCh: make(chan struct{}),
	}

	// 启动后台清理goroutine
	go mc.cleanupExpired()

	return mc
}

// Set 设置缓存项
// 参数：
//   - ctx: 上下文（用于保持与Redis接口一致，实际未使用）
//   - key: 缓存键
//   - value: 缓存值
//   - ttl: 过期时间，0表示永不过期
func (mc *MemoryCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	var expiration time.Time
	if ttl > 0 {
		expiration = time.Now().Add(ttl)
	}

	mc.data[key] = &cacheItem{
		value:      value,
		expiration: expiration,
	}

	return nil
}

// Get 获取缓存项
// 参数：
//   - ctx: 上下文
//   - key: 缓存键
//
// 返回值：
//   - interface{}: 缓存值，如果不存在或已过期返回nil
//   - error: 错误信息
func (mc *MemoryCache) Get(ctx context.Context, key string) (interface{}, error) {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	item, exists := mc.data[key]
	if !exists {
		return nil, fmt.Errorf("key not found")
	}

	// 检查是否过期
	if !item.expiration.IsZero() && time.Now().After(item.expiration) {
		return nil, fmt.Errorf("key expired")
	}

	return item.value, nil
}

// Delete 删除缓存项
// 参数：
//   - ctx: 上下文
//   - key: 缓存键
func (mc *MemoryCache) Delete(ctx context.Context, key string) error {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	delete(mc.data, key)
	return nil
}

// Exists 检查键是否存在且未过期
// 参数：
//   - ctx: 上下文
//   - key: 缓存键
//
// 返回值：
//   - bool: 键是否存在且有效
func (mc *MemoryCache) Exists(ctx context.Context, key string) bool {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	item, exists := mc.data[key]
	if !exists {
		return false
	}

	// 检查是否过期
	if !item.expiration.IsZero() && time.Now().After(item.expiration) {
		return false
	}

	return true
}

// cleanupExpired 定期清理过期的缓存项
// 每分钟执行一次清理操作
func (mc *MemoryCache) cleanupExpired() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			mc.mu.Lock()
			now := time.Now()
			for key, item := range mc.data {
				if !item.expiration.IsZero() && now.After(item.expiration) {
					delete(mc.data, key)
				}
			}
			mc.mu.Unlock()
		case <-mc.stopCh:
			return
		}
	}
}

// Close 关闭缓存服务
// 停止后台清理goroutine并清空所有数据
func (mc *MemoryCache) Close() error {
	close(mc.stopCh)

	mc.mu.Lock()
	defer mc.mu.Unlock()
	mc.data = make(map[string]*cacheItem)

	return nil
}

// Clear 清空所有缓存
func (mc *MemoryCache) Clear() {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	mc.data = make(map[string]*cacheItem)
}

// Size 返回当前缓存项数量（包括已过期但未清理的）
func (mc *MemoryCache) Size() int {
	mc.mu.RLock()
	defer mc.mu.RUnlock()
	return len(mc.data)
}

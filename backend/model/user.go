// Package model 定义了系统中的数据模型（实体类）
// 这些模型对应数据库表结构，用于GORM ORM操作
package model

import "time"

// User 用户模型，对应数据库中的 users 表
//
// 字段说明：
//   - ID: 用户唯一标识（主键，自增）
//   - Account: 登录账号（唯一，由字母和数字组成）
//   - Nickname: 用户昵称（可修改，但30天内只能修改一次）
//   - Password: 密码哈希值（bcrypt加密后的字符串，不会返回给前端）
//   - AvatarURL: 头像图片URL地址
//   - IsAdmin: 是否为管理员（true=管理员，false=普通用户）
//   - CreatedAt: 账号创建时间
//   - UpdatedAt: 最后更新时间（GORM自动维护）
//
// 数据库表结构对应：
//   CREATE TABLE users (
//     id BIGSERIAL PRIMARY KEY,
//     account VARCHAR(50) UNIQUE NOT NULL,
//     nickname VARCHAR(50) NOT NULL,
//     password_hash VARCHAR(255) NOT NULL,
//     avatar_url VARCHAR(500),
//     wechat_id VARCHAR(64),
//     is_admin BOOLEAN DEFAULT FALSE,
//     last_nickname_changed_at TIMESTAMP,
//     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
//     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
//   );
//
// 注意事项：
//   - Password字段使用 json:"-" 标签，表示JSON序列化时忽略此字段
//     这确保密码哈希永远不会被返回给前端
//   - 实际数据库表中还有 wechat_id 和 last_nickname_changed_at 字段
//     后续需要根据需求补充到模型中
type User struct {
	ID                    int64      `json:"id" gorm:"primaryKey"`                         // 主键ID
	Account               string     `json:"account" gorm:"uniqueIndex;size:50"`           // 登录账号（唯一索引）
	Nickname              string     `json:"nickname" gorm:"size:50"`                      // 用户昵称
	Password              string     `json:"-" gorm:"column:password_hash;size:255"`       // 密码哈希（不返回给前端）
	AvatarUrl             string     `json:"avatar_url" gorm:"column:avatar_url;size:500"` // 头像URL
	IsAdmin               bool       `json:"is_admin" gorm:"default:false"`                // 是否管理员
	WechatID              string     `json:"wechat_id" gorm:"size:64"`                     // 微信号
	LastNicknameChangedAt *time.Time `json:"last_nickname_changed_at" gorm:"index"`        // 最后昵称修改时间
	CreatedAt             time.Time  `json:"created_at" gorm:"autoCreateTime"`             // 创建时间
	UpdatedAt             time.Time  `json:"updated_at" gorm:"autoUpdateTime"`             // 更新时间
}

// TableName 指定User模型对应的数据库表名
// GORM默认会使用结构体名的复数形式作为表名（User -> users）
// 这里显式指定是为了确保一致性，避免命名问题
func (User) TableName() string {
	return "users"
}

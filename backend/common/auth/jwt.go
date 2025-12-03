// Package auth 提供JWT（JSON Web Token）相关的工具函数
// JWT用于实现无状态的用户认证
package auth

import (
	"time"

	jwt "github.com/golang-jwt/jwt/v5"

	"github.com/yycy134679/school-secondhand-trading-system/backend/config"
)

// JWT（JSON Web Token）说明：
// JWT是一种用于在双方之间安全传输信息的开放标准（RFC 7519）
// 特点：
// 1. 无状态：服务器不需要存储session，token包含所有必要信息
// 2. 跨域友好：适合前后端分离和微服务架构
// 3. 可扩展：可以在payload中存储自定义信息
//
// JWT结构（三部分，用.分隔）：
// Header.Payload.Signature
// 示例：eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOjEyM30.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c
//
// Header（头部）：包含算法和token类型
// {
//   "alg": "HS256",  // 签名算法（HMAC SHA256）
//   "typ": "JWT"     // token类型
// }
//
// Payload（负载）：包含用户信息和其他声明
// {
//   "userId": 123,           // 自定义字段：用户ID
//   "isAdmin": false,        // 自定义字段：是否管理员
//   "exp": 1735300000,       // 标准字段：过期时间（Unix时间戳）
//   "iat": 1735296400,       // 标准字段：签发时间
//   "iss": "secondhand-app"  // 标准字段：签发者
// }
//
// Signature（签名）：用于验证token的完整性
// 计算方式：HMACSHA256(base64(header) + "." + base64(payload), secret)

// GenerateToken 生成JWT token（待实现）
//
// 功能说明：
//   - 根据用户ID和其他信息生成JWT token
//   - 使用配置中的JWT_SECRET作为签名密钥
//   - 支持普通登录（短期token）和"记住我"（长期token）
//
// 参数：
//   - userID: 用户ID
//   - isAdmin: 是否管理员
//   - rememberMe: 是否记住登录（影响token有效期）
//
// 返回值：
//   - string: JWT token字符串
//   - error: 生成失败时返回错误
//
// Token有效期：
//   - 普通登录：1小时（3600秒）
//   - 记住我：7天（604800秒）
//
// 使用示例：
//
//	token, err := GenerateToken(user.ID, user.IsAdmin, true)
//	if err != nil {
//	    return err
//	}
//	// 将token返回给前端
//
// TODO: 使用 github.com/golang-jwt/jwt/v5 实现
//  1. 定义Claims结构体（包含userID、isAdmin、exp等字段）
//  2. 创建token对象：jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
//  3. 使用密钥签名：token.SignedString([]byte(jwtSecret))
func GenerateToken(userID int64) (string, error) {
	// 创建claims
	expiration := time.Now().Add(time.Hour)
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     expiration.Unix(),
		"iat":     time.Now().Unix(),
	}

	// 创建token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名并返回token字符串
	// 从配置中获取JWT密钥，如果未设置则使用默认值
	jwtSecret := "please-change-this" // 默认密钥
	if config, err := config.LoadConfig(); err == nil {
		jwtSecret = config.JWTSecret
	}
	return token.SignedString([]byte(jwtSecret))
}

// ParseToken 解析并验证JWT token（待实现）
//
// 功能说明：
//   - 解析JWT token字符串
//   - 验证签名是否正确
//   - 检查token是否过期
//   - 提取用户信息（userID、isAdmin）
//
// 参数：
//   - token: JWT token字符串（不包含"Bearer "前缀）
//
// 返回值：
//   - int64: 用户ID
//   - bool: 是否管理员
//   - error: 解析或验证失败时返回错误
//
// 可能的错误：
//   - token格式错误
//   - 签名验证失败（token被篡改）
//   - token已过期
//   - Claims格式不正确
//
// 使用示例：
//
//	userID, isAdmin, err := ParseToken(tokenString)
//	if err != nil {
//	    // token无效
//	    return errors.New("token无效或已过期")
//	}
//	// 使用userID和isAdmin进行后续处理
//
// TODO: 使用 github.com/golang-jwt/jwt/v5 实现
//  1. 解析token：jwt.Parse(tokenString, keyFunc)
//  2. 在keyFunc中返回验证密钥：[]byte(config.JWTSecret)
//  3. 验证token有效性：token.Valid
//  4. 提取Claims并返回userID和isAdmin
func ParseToken(token string) (int64, error) {
	// TODO: 实现JWT解析逻辑
	// 伪代码示例：
	//
	// // 1. 解析token
	// parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
	//     // 验证签名算法
	//     if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
	//         return nil, fmt.Errorf("unexpected signing method")
	//     }
	//     // 返回验证密钥
	//     return []byte(config.JWTSecret), nil
	// })
	//
	// if err != nil {
	//     return 0, false, err
	// }
	//
	// // 2. 验证token有效性
	// if !parsedToken.Valid {
	//     return 0, false, fmt.Errorf("invalid token")
	// }
	//
	// // 3. 提取Claims
	// claims, ok := parsedToken.Claims.(jwt.MapClaims)
	// if !ok {
	//     return 0, false, fmt.Errorf("invalid claims")
	// }
	//
	// // 4. 提取userID和isAdmin
	// userID := int64(claims["userId"].(float64))
	// isAdmin := claims["isAdmin"].(bool)
	//
	// return userID, isAdmin, nil
	return 0, nil
}

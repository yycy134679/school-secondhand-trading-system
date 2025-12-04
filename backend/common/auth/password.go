// Package auth 提供认证相关的工具函数
// 包括密码哈希、JWT生成和验证等
package auth

import "golang.org/x/crypto/bcrypt"

// bcrypt算法说明：
// bcrypt是一种基于Blowfish加密算法的密码哈希函数
// 特点：
// 1. 单向哈希：无法从哈希值反推原始密码
// 2. 自动加盐：每次哈希都会生成随机盐值，相同密码的哈希结果也不同
// 3. 计算成本可调：可以通过cost参数控制哈希计算时间，防止暴力破解
// 4. 安全性高：即使数据库泄露，攻击者也很难破解密码
//
// DefaultCost = 10，表示2^10次迭代，哈希一次约需100ms
// 这个时间对用户来说几乎无感，但对暴力破解者来说成本很高

// HashPassword 对明文密码进行bcrypt哈希加密
//
// 使用场景：
//   - 用户注册时，对原始密码进行加密存储
//   - 用户修改密码时，对新密码进行加密
//
// 参数：
//   - pw: 明文密码字符串
//
// 返回值：
//   - string: bcrypt哈希值（包含盐值和哈希结果），长度约60字符
//            格式示例：$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy
//   - error: 哈希失败时返回错误（一般不会失败）
//
// 安全性说明：
//   - 哈希值包含了算法版本、cost参数、盐值和最终哈希
//   - 同一个密码多次哈希，结果都不同（因为盐值随机）
//   - 存储在数据库中的是哈希值，而不是明文密码
//
// 使用示例：
//   hash, err := HashPassword("mypassword123")
//   if err != nil {
//       return err
//   }
//   // 将hash存入数据库的password_hash字段
func HashPassword(pw string) (string, error) {
	// GenerateFromPassword 执行以下步骤：
	// 1. 生成随机盐值
	// 2. 将密码和盐值组合
	// 3. 执行2^cost次bcrypt哈希迭代
	// 4. 返回包含所有信息的哈希字符串
	b, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	return string(b), err
}

// ComparePassword 验证明文密码是否与哈希值匹配
//
// 使用场景：
//   - 用户登录时，验证输入的密码是否正确
//   - 用户修改密码时，验证旧密码是否正确
//
// 参数：
//   - hash: 存储在数据库中的bcrypt哈希值
//   - pw: 用户输入的明文密码
//
// 返回值：
//   - error: 如果密码匹配，返回nil；如果不匹配，返回错误
//
// 工作原理：
//   1. 从hash中提取盐值和cost参数
//   2. 使用相同的盐值和cost对输入密码进行哈希
//   3. 比较两个哈希值是否相同
//
// 使用示例：
//   err := ComparePassword(user.PasswordHash, "inputPassword")
//   if err != nil {
//       // 密码错误
//       return errors.New("用户名或密码错误")
//   }
//   // 密码正确，允许登录
func ComparePassword(hash, pw string) error {
	// CompareHashAndPassword 是常量时间比较
	// 即使密码不匹配，也会花费相同的时间
	// 这可以防止时序攻击（timing attack）
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pw))
}

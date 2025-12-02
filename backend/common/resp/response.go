// Package resp 提供统一的API响应格式和辅助函数
// 确保所有API接口返回格式一致，便于前端解析和错误处理
package resp

import "github.com/gin-gonic/gin"

// 统一响应格式说明：
// {
//   "code": 0,           // 业务状态码：0表示成功，非0表示各种业务错误
//   "message": "ok",     // 响应消息：成功时为"ok"，失败时为错误描述
//   "data": {...}        // 响应数据：成功时包含业务数据，失败时为null
// }
//
// 注意：HTTP状态码统一为200，实际的业务成功/失败通过code字段区分
// 这种设计的优点：
// - 前端可以统一处理HTTP 200响应
// - 通过code字段精确区分不同的业务错误
// - 避免HTTP状态码与业务错误码混淆

// Success 返回成功响应
//
// 使用场景：
//   - API调用成功并需要返回数据时
//   - 例如：查询用户信息、获取商品列表等
//
// 参数：
//   - c: Gin上下文对象
//   - data: 要返回的业务数据，可以是任意类型（结构体、map、数组等）
//
// 响应格式示例：
//   {
//     "code": 0,
//     "message": "ok",
//     "data": {
//       "id": 123,
//       "username": "张三"
//     }
//   }
//
// 使用示例：
//   user := &User{ID: 123, Username: "张三"}
//   resp.Success(c, user)
func Success(c *gin.Context, data interface{}) {
	c.JSON(200, gin.H{
		"code":    0,    // 成功状态码固定为0
		"message": "ok", // 成功消息固定为"ok"
		"data":    data, // 返回的业务数据
	})
}

// Error 返回错误响应
//
// 使用场景：
//   - API调用失败时
//   - 例如：参数验证失败、用户未登录、权限不足等
//
// 参数：
//   - c: Gin上下文对象
//   - code: 业务错误码，定义在 common/errors/codes.go 中
//     常见错误码：
//     1001 - 参数错误
//     1002 - 未认证（未登录）
//     1003 - 无权限
//     2001 - 用户不存在
//     3001 - 商品不存在
//   - message: 错误描述信息，向用户说明错误原因
//
// 响应格式示例：
//   {
//     "code": 1001,
//     "message": "用户名不能为空",
//     "data": null
//   }
//
// 使用示例：
//   resp.Error(c, errors.CodeInvalidParams, "用户名不能为空")
func Error(c *gin.Context, code int, message string) {
	c.JSON(200, gin.H{
		"code":    code,    // 业务错误码
		"message": message, // 错误描述信息
		"data":    nil,     // 错误时data字段固定为null
	})
}

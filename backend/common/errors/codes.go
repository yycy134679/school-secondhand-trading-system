// Package errors 定义了系统中所有的业务错误码
// 错误码用于在API响应中标识不同类型的业务错误
package errors

// 业务错误码分类说明：
// 1xxx - 通用错误（参数、权限等）
// 2xxx - 用户模块错误
// 3xxx - 商品模块错误
// 4xxx - 分类/标签模块错误
// 5xxx - 推荐系统错误
//
// 设计原则：
// - 错误码应该具有语义化，便于前端和开发者理解
// - 每个错误码对应一种明确的错误场景
// - 与API文档中定义的错误码保持一致

// ============ 通用错误码（1xxx）============

const (
	// CodeInvalidParams 表示请求参数错误
	// 使用场景：
	// - 必填参数缺失（如用户名、密码为空）
	// - 参数格式错误（如邮箱格式不正确）
	// - 参数值不合法（如价格为负数）
	// 示例消息："用户名不能为空"、"邮箱格式不正确"
	CodeInvalidParams = 1001

	// CodeUnauthenticated 表示用户未登录或token无效
	// 使用场景：
	// - 访问需要登录的接口但未提供token
	// - token已过期
	// - token签名验证失败
	// 示例消息："请先登录"、"登录已过期，请重新登录"
	CodeUnauthenticated = 1002

	// CodeForbidden 表示用户无权限访问
	// 使用场景：
	// - 普通用户访问管理员接口
	// - 用户尝试编辑他人的商品
	// - 用户尝试删除他人的评论
	// 示例消息："无权限访问"、"只能编辑自己的商品"
	CodeForbidden = 1003

	// CodeInternal 表示服务器内部错误
	// 使用场景：
	// - 服务层或数据库发生未预期错误
	// - 外部依赖调用失败
	// 示例消息："服务器内部错误"
	CodeInternal = 1009
)

// ============ 用户模块错误码（2xxx）============
// TODO: 根据需求添加更多用户相关错误码
// 示例：
// const (
//     CodeUserNotFound      = 2001 // 用户不存在
//     CodeUserAlreadyExists = 2002 // 用户已存在
//     CodeWrongPassword     = 2003 // 密码错误
//     CodeNicknameLimit     = 2004 // 昵称修改频率限制
// )

// ============ 商品模块错误码（3xxx）============
// TODO: 根据需求添加更多商品相关错误码
// 示例：
// const (
//     CodeProductNotFound    = 3001 // 商品不存在
//     CodeProductNotForSale  = 3002 // 商品未在售
//     CodeInvalidStatus      = 3003 // 非法的状态转换
//     CodeStatusIsFinal      = 3004 // 状态是终态，不可修改
//     CodeUndoTimeout        = 3005 // 撤销超时（超过3秒窗口）
// )

// ============ 分类/标签模块错误码（4xxx）============
// TODO: 根据需求添加分类和标签相关错误码
// 示例：
// const (
//     CodeCategoryInUse = 4001 // 分类被商品引用，不可删除
//     CodeTagInUse      = 4002 // 标签被商品引用，不可删除
// )

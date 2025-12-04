package category

import "errors"

// 错误码定义
const (
	ErrCodeCategoryHasProducts = 4001 // 分类下有商品，无法删除
)

// 错误定义
var (
	ErrCategoryHasProducts = errors.New("category has products, cannot delete")
)

package tag

import "errors"

// 错误码定义
const (
	ErrCodeTagHasProducts = 4002 // 标签下有商品，无法删除
)

// 错误定义
var (
	ErrTagHasProducts = errors.New("tag has products, cannot delete")
)

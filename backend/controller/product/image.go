package product

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/yycy134679/school-secondhand-trading-system/backend/common/resp"
	"github.com/yycy134679/school-secondhand-trading-system/backend/service/product"
)

// ImageController 图片控制器
type ImageController struct {
	productService *product.ProductService
}

// NewImageController 创建图片控制器实例
func NewImageController(productService *product.ProductService) *ImageController {
	return &ImageController{
		productService: productService,
	}
}

// UploadProductImage 上传商品图片
// POST /api/v1/products/:id/images
func (ic *ImageController) UploadProductImage(c *gin.Context) {
	// 从上下文中获取用户ID
	userIDStr, exists := c.Get("user_id")
	if !exists {
		resp.Error(c, 401, "用户未登录")
		return
	}

	userID, err := strconv.ParseInt(userIDStr.(string), 10, 64)
	if err != nil {
		resp.Error(c, 400, "无效的用户ID")
		return
	}

	// 获取商品ID
	productIDStr := c.Param("id")
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		resp.Error(c, 400, "无效的商品ID")
		return
	}

	// 获取上传的文件
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		resp.Error(c, 400, "请上传有效的图片文件")
		return
	}
	defer file.Close()

	// 调用服务层方法上传图片
	imageDTO, err := ic.productService.AddProductImage(c.Request.Context(), userID, productID, file, header)
	if err != nil {
		// 根据错误类型返回不同的错误码
		if err.Error() == "商品不存在" {
			resp.Error(c, 404, err.Error())
		} else if err.Error() == "无权限操作该商品" {
			resp.Error(c, 403, err.Error())
		} else if err.Error() == "已售出的商品不能修改" {
			resp.Error(c, 400, err.Error())
		} else {
			resp.Error(c, 400, err.Error())
		}
		return
	}

	resp.Success(c, imageDTO)
}

// SetPrimaryImage 设置主图
// PUT /api/v1/products/:id/images/:imageId/primary
func (ic *ImageController) SetPrimaryImage(c *gin.Context) {
	// 从上下文中获取用户ID
	userIDStr, exists := c.Get("user_id")
	if !exists {
		resp.Error(c, 401, "用户未登录")
		return
	}

	userID, err := strconv.ParseInt(userIDStr.(string), 10, 64)
	if err != nil {
		resp.Error(c, 400, "无效的用户ID")
		return
	}

	// 获取商品ID和图片ID
	productIDStr := c.Param("id")
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		resp.Error(c, 400, "无效的商品ID")
		return
	}

	imageIDStr := c.Param("imageId")
	imageID, err := strconv.ParseInt(imageIDStr, 10, 64)
	if err != nil {
		resp.Error(c, 400, "无效的图片ID")
		return
	}

	// 调用服务层方法设置主图
	err = ic.productService.SetPrimaryImage(c.Request.Context(), userID, productID, imageID)
	if err != nil {
		// 根据错误类型返回不同的错误码
		if err.Error() == "商品不存在" || err.Error() == "图片不存在" {
			resp.Error(c, 404, err.Error())
		} else if err.Error() == "无权限操作该商品" {
			resp.Error(c, 403, err.Error())
		} else if err.Error() == "已售出的商品不能修改" {
			resp.Error(c, 400, err.Error())
		} else {
			resp.Error(c, 400, err.Error())
		}
		return
	}

	resp.Success(c, gin.H{"message": "设置主图成功"})
}

// UpdateImageSortOrder 更新图片排序
// PATCH /api/v1/products/:id/images/:imageId
func (ic *ImageController) UpdateImageSortOrder(c *gin.Context) {
	// 从上下文中获取用户ID
	userIDStr, exists := c.Get("user_id")
	if !exists {
		resp.Error(c, 401, "用户未登录")
		return
	}

	userID, err := strconv.ParseInt(userIDStr.(string), 10, 64)
	if err != nil {
		resp.Error(c, 400, "无效的用户ID")
		return
	}

	// 获取商品ID和图片ID
	productIDStr := c.Param("id")
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		resp.Error(c, 400, "无效的商品ID")
		return
	}

	imageIDStr := c.Param("imageId")
	imageID, err := strconv.ParseInt(imageIDStr, 10, 64)
	if err != nil {
		resp.Error(c, 400, "无效的图片ID")
		return
	}

	// 解析请求体中的sortOrder
	type SortOrderRequest struct {
		SortOrder int `json:"sortOrder" binding:"required,min=1"`
	}
	var req SortOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.Error(c, 400, "请提供有效的排序序号（正整数）")
		return
	}

	// 调用服务层方法更新排序
	err = ic.productService.UpdateImageSortOrder(c.Request.Context(), userID, productID, imageID, req.SortOrder)
	if err != nil {
		// 根据错误类型返回不同的错误码
		if err.Error() == "商品不存在" || err.Error() == "图片不存在" {
			resp.Error(c, 404, err.Error())
		} else if err.Error() == "无权限操作该商品" {
			resp.Error(c, 403, err.Error())
		} else if err.Error() == "已售出的商品不能修改" {
			resp.Error(c, 400, err.Error())
		} else {
			resp.Error(c, 400, err.Error())
		}
		return
	}

	resp.Success(c, gin.H{"message": "更新排序成功"})
}

// DeleteProductImage 删除商品图片
// DELETE /api/v1/products/:id/images/:imageId
func (ic *ImageController) DeleteProductImage(c *gin.Context) {
	// 从上下文中获取用户ID
	userIDStr, exists := c.Get("user_id")
	if !exists {
		resp.Error(c, 401, "用户未登录")
		return
	}

	userID, err := strconv.ParseInt(userIDStr.(string), 10, 64)
	if err != nil {
		resp.Error(c, 400, "无效的用户ID")
		return
	}

	// 获取商品ID和图片ID
	productIDStr := c.Param("id")
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		resp.Error(c, 400, "无效的商品ID")
		return
	}

	imageIDStr := c.Param("imageId")
	imageID, err := strconv.ParseInt(imageIDStr, 10, 64)
	if err != nil {
		resp.Error(c, 400, "无效的图片ID")
		return
	}

	// 调用服务层方法删除图片
	err = ic.productService.DeleteProductImage(c.Request.Context(), userID, productID, imageID)
	if err != nil {
		// 根据错误类型返回不同的错误码
		if err.Error() == "商品不存在" || err.Error() == "图片不存在" {
			resp.Error(c, 404, err.Error())
		} else if err.Error() == "无权限操作该商品" {
			resp.Error(c, 403, err.Error())
		} else if err.Error() == "已售出的商品不能修改" {
			resp.Error(c, 400, err.Error())
		} else if err.Error() == "至少保留一张图片" {
			resp.Error(c, 400, err.Error())
		} else {
			resp.Error(c, 400, err.Error())
		}
		return
	}

	resp.Success(c, gin.H{"message": "删除图片成功"})
}

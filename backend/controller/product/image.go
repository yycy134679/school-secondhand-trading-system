// Package product 提供商品相关控制器
package product

import (
	"context"
	"net/http"
	"strconv"

	"github.com/yycy134679/school-secondhand-trading-system/backend/common/resp"
	"github.com/yycy134679/school-secondhand-trading-system/backend/middleware"

	"mime/multipart"

	"github.com/gin-gonic/gin"
)

// ProductService 商品服务接口
type ProductService interface {
	// UploadProductImage 上传商品图片
	UploadProductImage(ctx context.Context, userID, productID int64, file multipart.File, header *multipart.FileHeader) (imageID int64, url string, err error)
	// SetPrimaryImage 设置商品主图
	SetPrimaryImage(ctx context.Context, userID, productID, imageID int64) error
	// UpdateImageSort 更新图片排序
	UpdateImageSort(ctx context.Context, userID, productID, imageID int64, sortOrder int) error
	// DeleteProductImage 删除商品图片
	DeleteProductImage(ctx context.Context, userID, productID, imageID int64) error
}

// ImageController 商品图片控制器
type ImageController struct {
	productService ProductService
}

// NewImageController 创建图片控制器实例
func NewImageController(productService ProductService) *ImageController {
	return &ImageController{
		productService: productService,
	}
}

// RegisterRoutes 注册图片相关路由
func (ic *ImageController) RegisterRoutes(router *gin.RouterGroup) {
	// 图片管理相关接口都需要登录
	router.Use(middleware.AuthMiddleware())

	// 上传图片
	router.POST("/products/:id/images", ic.UploadImage)
	// 设置主图
	router.PUT("/products/:id/images/:imageId/primary", ic.SetPrimaryImage)
	// 更新排序
	router.PATCH("/products/:id/images/:imageId", ic.UpdateImageSort)
	// 删除图片
	router.DELETE("/products/:id/images/:imageId", ic.DeleteImage)
}

// UploadImage 上传商品图片
// @Summary 上传商品图片
// @Description 上传商品图片，自动处理主图设置
// @Tags 商品图片管理
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "商品ID"
// @Param image formData file true "图片文件"
// @Success 200 {object} resp.Response{data=map[string]interface{}}
// @Failure 400 {object} resp.Response
// @Failure 401 {object} resp.Response
// @Failure 403 {object} resp.Response
// @Router /api/v1/products/{id}/images [post]
func (ic *ImageController) UploadImage(c *gin.Context) {
	// 获取商品ID
	productID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		resp.Error(c, http.StatusBadRequest, "无效的商品ID")
		return
	}

	// 获取当前用户ID
	userID := c.GetInt64("userID")

	// 获取上传的文件
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		resp.Error(c, http.StatusBadRequest, "文件上传失败")
		return
	}
	defer file.Close()

	// 调用服务层处理图片上传
	imageID, url, err := ic.productService.UploadProductImage(c.Request.Context(), userID, productID, file, header)
	if err != nil {
		resp.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	// 返回成功响应
	resp.Success(c, gin.H{
		"imageId": imageID,
		"url":     url,
	})
}

// SetPrimaryImage 设置主图
// @Summary 设置商品主图
// @Description 将指定图片设置为商品主图
// @Tags 商品图片管理
// @Produce json
// @Param id path int true "商品ID"
// @Param imageId path int true "图片ID"
// @Success 200 {object} resp.Response
// @Failure 400 {object} resp.Response
// @Failure 401 {object} resp.Response
// @Failure 403 {object} resp.Response
// @Router /api/v1/products/{id}/images/{imageId}/primary [put]
func (ic *ImageController) SetPrimaryImage(c *gin.Context) {
	// 获取商品ID和图片ID
	productID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		resp.Error(c, http.StatusBadRequest, "无效的商品ID")
		return
	}

	imageID, err := strconv.ParseInt(c.Param("imageId"), 10, 64)
	if err != nil {
		resp.Error(c, http.StatusBadRequest, "无效的图片ID")
		return
	}

	// 获取当前用户ID
	userID := c.GetInt64("userID")

	// 调用服务层设置主图
	if err := ic.productService.SetPrimaryImage(c.Request.Context(), userID, productID, imageID); err != nil {
		resp.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	// 返回成功响应
	resp.Success(c, nil)
}

// UpdateImageSort 更新图片排序
// @Summary 更新商品图片排序
// @Description 更新指定图片的排序顺序
// @Tags 商品图片管理
// @Accept json
// @Produce json
// @Param id path int true "商品ID"
// @Param imageId path int true "图片ID"
// @Param body body map[string]int true "排序信息"
// @Success 200 {object} resp.Response
// @Failure 400 {object} resp.Response
// @Failure 401 {object} resp.Response
// @Failure 403 {object} resp.Response
// @Router /api/v1/products/{id}/images/{imageId} [patch]
func (ic *ImageController) UpdateImageSort(c *gin.Context) {
	// 获取商品ID和图片ID
	productID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		resp.Error(c, http.StatusBadRequest, "无效的商品ID")
		return
	}

	imageID, err := strconv.ParseInt(c.Param("imageId"), 10, 64)
	if err != nil {
		resp.Error(c, http.StatusBadRequest, "无效的图片ID")
		return
	}

	// 获取当前用户ID
	userID := c.GetInt64("userID")

	// 解析请求体
	var req struct {
		SortOrder int `json:"sortOrder" binding:"required,min=0"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.Error(c, http.StatusBadRequest, "排序参数无效")
		return
	}

	// 调用服务层更新排序
	if err := ic.productService.UpdateImageSort(c.Request.Context(), userID, productID, imageID, req.SortOrder); err != nil {
		resp.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	// 返回成功响应
	resp.Success(c, nil)
}

// DeleteImage 删除图片
// @Summary 删除商品图片
// @Description 删除指定商品图片，自动处理主图更换
// @Tags 商品图片管理
// @Produce json
// @Param id path int true "商品ID"
// @Param imageId path int true "图片ID"
// @Success 200 {object} resp.Response
// @Failure 400 {object} resp.Response
// @Failure 401 {object} resp.Response
// @Failure 403 {object} resp.Response
// @Router /api/v1/products/{id}/images/{imageId} [delete]
func (ic *ImageController) DeleteImage(c *gin.Context) {
	// 获取商品ID和图片ID
	productID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		resp.Error(c, http.StatusBadRequest, "无效的商品ID")
		return
	}

	imageID, err := strconv.ParseInt(c.Param("imageId"), 10, 64)
	if err != nil {
		resp.Error(c, http.StatusBadRequest, "无效的图片ID")
		return
	}

	// 获取当前用户ID
	userID := c.GetInt64("userID")

	// 调用服务层删除图片
	if err := ic.productService.DeleteProductImage(c.Request.Context(), userID, productID, imageID); err != nil {
		resp.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	// 返回成功响应
	resp.Success(c, nil)
}

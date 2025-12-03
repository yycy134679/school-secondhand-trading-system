package admin

import (
	"strconv"

	"github.com/yycy134679/school-secondhand-trading-system/backend/common/resp"
	"github.com/yycy134679/school-secondhand-trading-system/backend/service/admin"

	"github.com/gin-gonic/gin"
)

// ProductController 商品管理控制器
type ProductController struct {
	adminService *admin.AdminService
}

// NewProductController 创建商品管理控制器
func NewProductController(adminService *admin.AdminService) *ProductController {
	return &ProductController{
		adminService: adminService,
	}
}

// ListProducts 商品列表接口
// GET /api/v1/admin/products
func (pc *ProductController) ListProducts(c *gin.Context) {
	// 获取查询参数
	status := c.Query("status")
	sellerIDStr := c.Query("sellerId")
	keyword := c.Query("keyword")
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "10")

	// 参数转换
	page, pageErr := strconv.Atoi(pageStr)
	if pageErr != nil || page < 1 {
		page = 1
	}

	pageSize, sizeErr := strconv.Atoi(pageSizeStr)
	if sizeErr != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// 卖家ID转换
	var sellerID int64 = 0
	if sellerIDStr != "" {
		id, sellerErr := strconv.ParseInt(sellerIDStr, 10, 64)
		if sellerErr == nil && id > 0 {
			sellerID = id
		}
	}

	// 调用服务层方法
	response, serviceErr := pc.adminService.ListProductsAdmin(c.Request.Context(), status, sellerID, keyword, page, pageSize)
	if serviceErr != nil {
		resp.Error(c, 500, "获取商品列表失败: "+serviceErr.Error())
		return
	}

	// 返回成功响应
	resp.Success(c, response)
}

// UpdateProduct 更新商品接口
// PUT /api/v1/admin/products/:id
func (pc *ProductController) UpdateProduct(c *gin.Context) {
	// 获取商品ID
	idStr := c.Param("id")
	productID, idErr := strconv.ParseInt(idStr, 10, 64)
	if idErr != nil || productID <= 0 {
		resp.Error(c, 400, "无效的商品ID")
		return
	}

	// 绑定请求体
	var req admin.UpdateProductRequest
	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		resp.Error(c, 400, "请求参数无效: "+bindErr.Error())
		return
	}

	// 调用服务层方法
	updateErr := pc.adminService.UpdateProductAsAdmin(c.Request.Context(), productID, req)
	if updateErr != nil {
		// 检查是否是禁止修改状态的错误
		errMsg := updateErr.Error()
		if len(errMsg) >= 4 && errMsg[:4] == "3004" {
			resp.Error(c, 3004, errMsg[4:])
			return
		}
		resp.Error(c, 500, "更新商品失败: "+errMsg)
		return
	}

	// 返回成功响应
	resp.Success(c, gin.H{"message": "商品更新成功"})
}

package main

import (
	"github.com/gin-gonic/gin"

	"school-secondhand-trading-system/controller/product"
	"school-secondhand-trading-system/router"
	productservice "school-secondhand-trading-system/service/product"
)

func main() {
	// 初始化服务
	productService := productservice.NewProductService()

	// 初始化控制器
	productController := product.NewProductController(productService)
	imageController := product.NewImageController(productService)

	// 创建Gin引擎
	engine := gin.Default()

	// 设置路由
	router.SetupProductRoutes(engine, productController, imageController)
}

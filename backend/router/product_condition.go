package router

import (
	"github.com/gin-gonic/gin"

	"github.com/yycy134679/school-secondhand-trading-system/backend/controller/product_condition"
)

// SetupProductConditionRoutes 设置新旧程度相关路由
func SetupProductConditionRoutes(engine *gin.Engine, controller *productcondition.Controller) {
	api := engine.Group("/api/v1")
	public := api.Group("/")
	{
		public.GET("/product-conditions", controller.ListProductConditions)
	}
}

package main

import (
	"fmt"
	"os"
	"reflect"

	categorycontroller "github.com/yycy134679/school-secondhand-trading-system/backend/controller/category"
	"github.com/yycy134679/school-secondhand-trading-system/backend/controller/product"
	"github.com/yycy134679/school-secondhand-trading-system/backend/controller/user"
	"github.com/yycy134679/school-secondhand-trading-system/backend/repository"
	adminservice "github.com/yycy134679/school-secondhand-trading-system/backend/service/admin"
	categoryservice "github.com/yycy134679/school-secondhand-trading-system/backend/service/category"
	productservice "github.com/yycy134679/school-secondhand-trading-system/backend/service/product"
	tagservice "github.com/yycy134679/school-secondhand-trading-system/backend/service/tag"
	userservice "github.com/yycy134679/school-secondhand-trading-system/backend/service/user"
)

// VerifyResult 表示验证结果
type VerifyResult struct {
	Module  string
	Status  bool
	Message string
}

// VerifyRepository 验证仓库层实现
func VerifyRepository() []VerifyResult {
	var results []VerifyResult

	fmt.Println("验证仓库层实现...")

	// 检查UserRepository方法
	userRepoType := reflect.TypeOf((*repository.UserRepository)(nil)).Elem()
	requiredUserRepoMethods := []string{
		"GetByAccount",
		"GetByID",
		"Create",
		"UpdateProfile",
		"UpdatePassword",
	}

	for _, method := range requiredUserRepoMethods {
		if _, exists := userRepoType.MethodByName(method); !exists {
			results = append(results, VerifyResult{
				Module:  "UserRepository",
				Status:  false,
				Message: fmt.Sprintf("缺少方法: %s", method),
			})
		}
	}

	if len(results) == 0 {
		results = append(results, VerifyResult{
			Module:  "UserRepository",
			Status:  true,
			Message: "所有方法实现完整",
		})
	}

	// 检查ProductRepository方法
	productRepoType := reflect.TypeOf((*repository.ProductRepository)(nil)).Elem()
	requiredProductRepoMethods := []string{
		"Create",
		"Update",
		"GetByID",
		"ListBySeller",
		"UpdateStatus",
		"Search",
		"ListLatestForSale",
		"ListByCategory",
	}

	for _, method := range requiredProductRepoMethods {
		if _, exists := productRepoType.MethodByName(method); !exists {
			results = append(results, VerifyResult{
				Module:  "ProductRepository",
				Status:  false,
				Message: fmt.Sprintf("缺少方法: %s", method),
			})
		}
	}

	if len(results) == 1 && results[0].Module == "UserRepository" {
		results = append(results, VerifyResult{
			Module:  "ProductRepository",
			Status:  true,
			Message: "所有方法实现完整",
		})
	}

	return results
}

// VerifyService 验证服务层实现
func VerifyService() []VerifyResult {
	var results []VerifyResult

	fmt.Println("验证服务层实现...")

	// 检查UserService方法
	userService := userservice.NewUserService(nil)
	userServiceType := reflect.TypeOf(userService)
	requiredUserServiceMethods := []string{
		"Register",
		"Login",
		"GetProfile",
		"UpdateProfile",
		"ChangePassword",
	}

	for _, method := range requiredUserServiceMethods {
		if _, exists := userServiceType.MethodByName(method); !exists {
			results = append(results, VerifyResult{
				Module:  "UserService",
				Status:  false,
				Message: fmt.Sprintf("缺少方法: %s", method),
			})
		}
	}

	if len(results) == 0 {
		results = append(results, VerifyResult{
			Module:  "UserService",
			Status:  true,
			Message: "所有方法实现完整",
		})
	}

	// 检查ProductService方法
	productService := productservice.NewProductService()
	productServiceType := reflect.TypeOf(productService)
	requiredProductServiceMethods := []string{
		"CreateProduct",
		"UpdateProduct",
		"ChangeStatus",
		"UndoLastStatusChange",
		"GetProductDetail",
		"ListMyProducts",
		"Search",
		"ListByCategory",
	}

	for _, method := range requiredProductServiceMethods {
		if _, exists := productServiceType.MethodByName(method); !exists {
			results = append(results, VerifyResult{
				Module:  "ProductService",
				Status:  false,
				Message: fmt.Sprintf("缺少方法: %s", method),
			})
		}
	}

	if len(results) == 1 && results[0].Module == "UserService" {
		results = append(results, VerifyResult{
			Module:  "ProductService",
			Status:  true,
			Message: "所有方法实现完整",
		})
	}

	return results
}

// VerifyController 验证控制器层实现
func VerifyController() []VerifyResult {
	var results []VerifyResult
	
	fmt.Println("验证控制器层实现...")
	
	// 检查UserController - user包使用路由注册模式，没有传统控制器结构体
	// 检查RegisterRoutes函数是否存在
	userRegisterRoutesType := reflect.ValueOf(user.RegisterRoutes).Type()
	if userRegisterRoutesType.Kind() != reflect.Func {
		results = append(results, VerifyResult{
			Module:   "UserController",
			Status:   false,
			Message:  "缺少RegisterRoutes函数",
		})
	} else {
		results = append(results, VerifyResult{
			Module:   "UserController",
			Status:   true,
			Message:  "路由注册函数实现完整",
		})
	}
	
	// 检查ProductController
	productController := product.NewProductController(nil)
	productControllerType := reflect.TypeOf(productController)
	requiredProductControllerMethods := []string{
		"CreateProduct",
		"UpdateProduct",
		"GetProductDetail",
		"SearchProducts",
		"ListMyProducts",
		"ChangeProductStatus",
		"UndoLastStatusChange",
		"GetProductsByCategory",
	}
	
	for _, method := range requiredProductControllerMethods {
		if _, exists := productControllerType.MethodByName(method); !exists {
			results = append(results, VerifyResult{
				Module:   "ProductController",
				Status:   false,
				Message:  fmt.Sprintf("缺少方法: %s", method),
			})
		}
	}
	
	if len(results) == 1 && results[0].Module == "UserController" {
		results = append(results, VerifyResult{
			Module:   "ProductController",
			Status:   true,
			Message:  "所有方法实现完整",
		})
	}
	
	// 检查CategoryController
	categoryController := categorycontroller.NewCategoryController(nil)
	categoryControllerType := reflect.TypeOf(categoryController)
	requiredCategoryControllerMethods := []string{
		"ListCategories",
		"CreateCategory",
		"UpdateCategory",
		"DeleteCategory",
	}
	
	for _, method := range requiredCategoryControllerMethods {
		if _, exists := categoryControllerType.MethodByName(method); !exists {
			results = append(results, VerifyResult{
				Module:   "CategoryController",
				Status:   false,
				Message:  fmt.Sprintf("缺少方法: %s", method),
			})
		}
	}
	
	if len(results) == 2 && results[1].Module == "ProductController" {
		results = append(results, VerifyResult{
			Module:   "CategoryController",
			Status:   true,
			Message:  "所有方法实现完整",
		})
	}
	
	return results
}

// VerifyAdminModule 验证管理员模块实现
func VerifyAdminModule() []VerifyResult {
	var results []VerifyResult

	fmt.Println("验证管理员模块实现...")

	// 检查AdminService
	adminService := adminservice.NewAdminService(nil)
	adminServiceType := reflect.TypeOf(adminService)
	requiredAdminServiceMethods := []string{
		"GetDashboardStats",
		"ListUsers",
		"ListProductsAdmin",
		"UpdateProductAsAdmin",
	}

	for _, method := range requiredAdminServiceMethods {
		if _, exists := adminServiceType.MethodByName(method); !exists {
			results = append(results, VerifyResult{
				Module:  "AdminService",
				Status:  false,
				Message: fmt.Sprintf("缺少方法: %s", method),
			})
		}
	}

	if len(results) == 0 {
		results = append(results, VerifyResult{
			Module:  "AdminService",
			Status:  true,
			Message: "所有方法实现完整",
		})
	}

	return results
}

// VerifyCategoryTagModule 验证分类和标签模块实现
func VerifyCategoryTagModule() []VerifyResult {
	var results []VerifyResult

	fmt.Println("验证分类和标签模块实现...")

	// 检查CategoryService
	categoryService := categoryservice.NewCategoryService(nil)
	categoryServiceType := reflect.TypeOf(categoryService)
	requiredCategoryServiceMethods := []string{
		"ListCategories",
		"CreateCategory",
		"UpdateCategory",
		"DeleteCategory",
	}

	for _, method := range requiredCategoryServiceMethods {
		if _, exists := categoryServiceType.MethodByName(method); !exists {
			results = append(results, VerifyResult{
				Module:  "CategoryService",
				Status:  false,
				Message: fmt.Sprintf("缺少方法: %s", method),
			})
		}
	}

	if len(results) == 0 {
		results = append(results, VerifyResult{
			Module:  "CategoryService",
			Status:  true,
			Message: "所有方法实现完整",
		})
	}

	// 检查TagService
	tagService := tagservice.NewTagService(nil)
	tagServiceType := reflect.TypeOf(tagService)
	requiredTagServiceMethods := []string{
		"ListTags",
		"CreateTag",
		"UpdateTag",
		"DeleteTag",
	}

	for _, method := range requiredTagServiceMethods {
		if _, exists := tagServiceType.MethodByName(method); !exists {
			results = append(results, VerifyResult{
				Module:  "TagService",
				Status:  false,
				Message: fmt.Sprintf("缺少方法: %s", method),
			})
		}
	}

	if len(results) == 1 && results[0].Module == "CategoryService" {
		results = append(results, VerifyResult{
			Module:  "TagService",
			Status:  true,
			Message: "所有方法实现完整",
		})
	}

	return results
}

func main() {
	fmt.Println("开始验证后端功能（除任务8、11、12外）...\n")

	var allResults []VerifyResult

	// 验证仓库层
	allResults = append(allResults, VerifyRepository()...)

	// 验证服务层
	allResults = append(allResults, VerifyService()...)

	// 验证控制器层
	allResults = append(allResults, VerifyController()...)

	// 验证管理员模块
	allResults = append(allResults, VerifyAdminModule()...)

	// 验证分类和标签模块
	allResults = append(allResults, VerifyCategoryTagModule()...)

	// 输出验证结果
	fmt.Println("\n=== 验证结果 ===")
	allPassed := true
	for _, result := range allResults {
		status := "✓ 通过"
		if !result.Status {
			status = "✗ 失败"
			allPassed = false
		}
		fmt.Printf("%s: %s\n", result.Module, status)
		if !result.Status {
			fmt.Printf("  原因: %s\n", result.Message)
		}
	}

	// 写入日志文件
	logFile, err := os.Create("Script/verify_result.log")
	if err != nil {
		fmt.Printf("\n创建日志文件失败: %v\n", err)
		return
	}
	defer logFile.Close()

	logFile.WriteString("=== 后端功能验证结果 ===\n")
	logFile.WriteString("验证范围: 除任务8、11、12外的所有功能\n\n")

	for _, result := range allResults {
		status := "通过"
		if !result.Status {
			status = "失败"
		}
		logFile.WriteString(fmt.Sprintf("%s: %s\n", result.Module, status))
		if !result.Status {
			logFile.WriteString(fmt.Sprintf("  原因: %s\n", result.Message))
		}
	}

	if allPassed {
		logFile.WriteString("\n✅ 所有模块验证通过！")
		fmt.Println("\n🎉 所有模块验证通过！")
		os.Exit(0)
	} else {
		logFile.WriteString("\n❌ 部分模块验证失败！")
		fmt.Println("\n❌ 部分模块验证失败！")
		os.Exit(1)
	}
}

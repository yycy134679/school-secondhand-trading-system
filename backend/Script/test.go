package main

import (
	"fmt"
	"os"
	"reflect"

	"github.com/yycy134679/school-secondhand-trading-system/backend/common/auth"
	"github.com/yycy134679/school-secondhand-trading-system/backend/common/errors"
	"github.com/yycy134679/school-secondhand-trading-system/backend/config"
	"github.com/yycy134679/school-secondhand-trading-system/backend/middleware"
	"github.com/yycy134679/school-secondhand-trading-system/backend/model"
)

// TestResult 表示测试结果

// TestConfig 测试配置管理
func TestConfig() bool {
	fmt.Println("测试配置管理...")
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("配置加载失败: %v\n", err)
		return false
	}
	fmt.Printf("配置加载成功: %+v\n", cfg)
	return true
}

// TestAuth 测试认证功能
func TestAuth() bool {
	fmt.Println("测试认证功能...")

	// 测试密码加密与验证
	password := "test123456"
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		fmt.Printf("密码加密失败: %v\n", err)
		return false
	}

	if err := auth.ComparePassword(hashedPassword, password); err != nil {
		fmt.Printf("密码验证失败: %v\n", err)
		return false
	}

	// 测试JWT生成
	token, err := auth.GenerateToken(1)
	if err != nil {
		fmt.Printf("JWT生成失败: %v\n", err)
		return false
	}
	fmt.Printf("JWT生成成功: %s\n", token)

	return true
}

// TestCommon 测试通用模块
func TestCommon() bool {
	fmt.Println("测试通用模块...")

	// 测试错误码定义
	errCodeMap := map[string]int{
		"CodeInvalidParams":   1001,
		"CodeUnauthenticated": 1002,
		"CodeForbidden":       1003,
	}

	for name, expected := range errCodeMap {
		actual := getErrorCode(name)
		if actual != expected {
			fmt.Printf("错误码 %s 不匹配: 期望 %d, 实际 %d\n", name, expected, actual)
			return false
		}
	}

	fmt.Println("通用模块测试通过")
	return true
}

// TestMiddleware 测试中间件
func TestMiddleware() bool {
	fmt.Println("测试中间件...")

	// 测试AuthMiddleware
	authMW := middleware.AuthMiddleware()
	if authMW == nil {
		fmt.Println("AuthMiddleware 创建失败")
		return false
	}

	// 测试AdminMiddleware
	adminMW := middleware.AdminMiddleware()
	if adminMW == nil {
		fmt.Println("AdminMiddleware 创建失败")
		return false
	}

	fmt.Println("中间件测试通过")
	return true
}

// TestModel 测试模型定义
func TestModel() bool {
	fmt.Println("测试模型定义...")

	// 检查User模型字段
	user := model.User{}
	userType := reflect.TypeOf(user)
	requiredUserFields := []string{"ID", "Account", "Nickname", "Password", "AvatarUrl", "WechatID", "IsAdmin", "CreatedAt", "UpdatedAt"}

	for _, field := range requiredUserFields {
		if _, exists := userType.FieldByName(field); !exists {
			fmt.Printf("User模型缺少字段: %s\n", field)
			return false
		}
	}

	// 检查Product模型字段
	product := model.Product{}
	productType := reflect.TypeOf(product)
	requiredProductFields := []string{"ID", "SellerID", "Title", "Description", "Price", "ConditionID", "CategoryID", "Status", "CreatedAt", "UpdatedAt"}
	
	for _, field := range requiredProductFields {
		if _, exists := productType.FieldByName(field); !exists {
			fmt.Printf("Product模型缺少字段: %s\n", field)
			return false
		}
	}
	
	// 检查ProductImage模型字段
	productImage := model.ProductImage{}
	productImageType := reflect.TypeOf(productImage)
	requiredImageFields := []string{"ID", "ProductID", "URL", "IsPrimary", "SortOrder"}
	
	for _, field := range requiredImageFields {
		if _, exists := productImageType.FieldByName(field); !exists {
			fmt.Printf("ProductImage模型缺少字段: %s\n", field)
			return false
		}
	}
	
	// 检查DTO模型字段
	productCardDTO := model.ProductCardDTO{}
	productCardDTOType := reflect.TypeOf(productCardDTO)
	requiredCardFields := []string{"ID", "Title", "Price", "MainImage", "Status"}
	
	for _, field := range requiredCardFields {
		if _, exists := productCardDTOType.FieldByName(field); !exists {
			fmt.Printf("ProductCardDTO模型缺少字段: %s\n", field)
			return false
		}
	}
	
	productDetailDTO := model.ProductDetailDTO{}
	productDetailDTOType := reflect.TypeOf(productDetailDTO)
	requiredDetailFields := []string{"ID", "Title", "Description", "Price", "CategoryID", "ConditionID", "Images", "Tags"}
	
	for _, field := range requiredDetailFields {
		if _, exists := productDetailDTOType.FieldByName(field); !exists {
			fmt.Printf("ProductDetailDTO模型缺少字段: %s\n", field)
			return false
		}
	}
	
	fmt.Println("模型定义测试通过")
	return true
}

// 辅助函数：获取错误码
func getErrorCode(name string) int {
	switch name {
	case "CodeInvalidParams":
		return errors.CodeInvalidParams
	case "CodeUnauthenticated":
		return errors.CodeUnauthenticated
	case "CodeForbidden":
		return errors.CodeForbidden
	default:
		return 0
	}
}

func main() {
	fmt.Println("开始验证后端功能（除任务8、11、12外）...\n")

	results := make(map[string]bool)
	results["配置管理"] = TestConfig()
	results["认证功能"] = TestAuth()
	results["通用模块"] = TestCommon()
	results["中间件"] = TestMiddleware()
	results["模型定义"] = TestModel()

	fmt.Println("\n=== 测试结果 ===")
	allPassed := true
	for testName, passed := range results {
		status := "✓ 通过"
		if !passed {
			status = "✗ 失败"
			allPassed = false
		}
		fmt.Printf("%s: %s\n", testName, status)
	}

	if allPassed {
		fmt.Println("\n🎉 所有测试通过！")
		os.Exit(0)
	} else {
		fmt.Println("\n❌ 部分测试失败！")
		os.Exit(1)
	}
}

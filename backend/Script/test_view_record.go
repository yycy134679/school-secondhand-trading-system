// 浏览记录功能测试脚本
// 用于验证浏览记录相关功能的实现

package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/yycy134679/school-secondhand-trading-system/backend/repository"
	"github.com/yycy134679/school-secondhand-trading-system/backend/service/recommend"
)

func main() {
	// 连接数据库 (需要根据实际环境修改连接串)
	dsn := "host=localhost user=postgres password=postgres dbname=secondhand_dev port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	fmt.Println("✓ 数据库连接成功")

	// 初始化仓库
	viewRecordRepo := repository.NewViewRecordRepository(db)
	productRepo := repository.NewProductRepository(db)

	// 初始化推荐服务 (不使用Redis)
	recommendService := recommend.NewRecommendService(
		viewRecordRepo,
		productRepo,
		db,
		nil, // Redis可选
	)

	fmt.Println("✓ 服务初始化成功")

	ctx := context.Background()

	// 测试1: 添加浏览记录
	fmt.Println("\n--- 测试1: 添加浏览记录 ---")
	testUserID := int64(1)
	testProductID := int64(1)

	err = recommendService.RecordView(ctx, testUserID, testProductID)
	if err != nil {
		log.Printf("✗ 添加浏览记录失败: %v", err)
	} else {
		fmt.Printf("✓ 成功添加浏览记录: userID=%d, productID=%d\n", testUserID, testProductID)
	}

	// 添加更多浏览记录以测试
	for i := int64(2); i <= 5; i++ {
		err = recommendService.RecordView(ctx, testUserID, i)
		if err != nil {
			log.Printf("✗ 添加浏览记录失败: %v", err)
		} else {
			fmt.Printf("✓ 成功添加浏览记录: userID=%d, productID=%d\n", testUserID, i)
		}
		time.Sleep(100 * time.Millisecond) // 稍微延迟以确保时间戳不同
	}

	// 测试2: 获取最近浏览记录
	fmt.Println("\n--- 测试2: 获取最近浏览记录 ---")
	views, err := viewRecordRepo.ListRecentViews(ctx, testUserID, 10)
	if err != nil {
		log.Printf("✗ 获取浏览记录失败: %v", err)
	} else {
		fmt.Printf("✓ 成功获取 %d 条浏览记录\n", len(views))
		for i, view := range views {
			fmt.Printf("  [%d] ProductID: %d, ViewedAt: %s\n",
				i+1, view.ProductID, view.ViewedAt.Format("2006-01-02 15:04:05"))
		}
	}

	// 测试3: 获取推荐商品
	fmt.Println("\n--- 测试3: 获取推荐商品 ---")
	recommendations, err := recommendService.GetRecommendations(ctx, testUserID, 5)
	if err != nil {
		log.Printf("✗ 获取推荐失败: %v", err)
	} else {
		fmt.Printf("✓ 成功获取 %d 条推荐\n", len(recommendations))
		for i, product := range recommendations {
			fmt.Printf("  [%d] ID: %d, Title: %s, Price: %.2f\n",
				i+1, product.ID, product.Title, product.Price)
		}
	}

	// 测试4: 获取首页数据
	fmt.Println("\n--- 测试4: 获取首页数据 ---")
	homeData, err := recommendService.GetHomeData(ctx, &testUserID, 1, 10)
	if err != nil {
		log.Printf("✗ 获取首页数据失败: %v", err)
	} else {
		fmt.Printf("✓ 成功获取首页数据\n")
		fmt.Printf("  推荐商品数量: %d\n", len(homeData.Recommendations))
		fmt.Printf("  最新商品数量: %d\n", len(homeData.Latest))
		fmt.Printf("  商品总数: %d\n", homeData.TotalCount)
	}

	// 测试5: 获取浏览记录并关联商品
	fmt.Println("\n--- 测试5: 获取浏览记录并关联商品 ---")
	viewsWithProducts, err := recommendService.GetRecentViewsWithProducts(ctx, testUserID, 5)
	if err != nil {
		log.Printf("✗ 获取浏览记录失败: %v", err)
	} else {
		fmt.Printf("✓ 成功获取 %d 条浏览记录（含商品信息）\n", len(viewsWithProducts))
		for i, item := range viewsWithProducts {
			fmt.Printf("  [%d] ViewedAt: %s, Product: %s (ID: %d, Price: %.2f)\n",
				i+1,
				item.ViewedAt.Format("2006-01-02 15:04:05"),
				item.Product.Title,
				item.Product.ID,
				item.Product.Price)
		}
	}

	fmt.Println("\n=== 所有测试完成 ===")
}

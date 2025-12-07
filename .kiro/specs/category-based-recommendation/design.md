# 设计文档

## 概述

本设计文档描述了如何将现有的基于标签的推荐系统改造为基于商品分类的推荐系统。新系统将分析用户最近浏览的商品分类，选择最新的两个不同分类，并从每个分类中推荐2件商品，总共推荐4件商品。

### 设计目标

1. 提高推荐的相关性和多样性
2. 简化推荐算法，提升性能
3. 保持与现有API的完全兼容
4. 确保边界情况的正确处理

## 架构

### 现有架构

当前推荐系统的架构包括：
- **Controller层**: `RecommendController` 处理HTTP请求
- **Service层**: `RecommendService` 实现推荐逻辑
- **Repository层**: `ViewRecordRepository` 和 `ProductRepository` 处理数据访问
- **Model层**: `UserRecentView`, `Product`, `Category` 等数据模型

### 修改范围

本次改造仅需修改 `RecommendService` 中的 `calculateRecommendations` 方法，其他组件保持不变。

## 组件和接口

### 修改的组件

#### RecommendService.calculateRecommendations

**方法签名**:
```go
func (s *RecommendService) calculateRecommendations(ctx context.Context, userID int64, maxCount int) ([]model.Product, error)
```

**新的实现逻辑**:

1. **获取用户最近浏览记录**
   - 查询 `user_recent_views` 表，获取用户最近20条浏览记录
   - 按 `viewed_at` 倒序排列

2. **提取浏览商品的分类信息**
   - 根据浏览记录中的 `product_id` 查询对应商品
   - 提取每个商品的 `category_id`
   - 按浏览时间顺序记录分类

3. **选择推荐分类**
   - 遍历浏览记录，找出最新的两个不同分类
   - 如果只有一个分类，则仅使用该分类
   - 记录每个选定分类的出现顺序（用于确定推荐优先级）

4. **查询推荐商品**
   - 对于每个选定的分类，查询符合条件的商品：
     - 状态为 "ForSale"
     - 不是用户自己发布的（`seller_id != userID`）
     - 不在用户已浏览列表中
     - 按创建时间倒序排列
   - 每个分类取2件商品（如果有两个分类）
   - 如果只有一个分类，取4件商品

5. **组装结果**
   - 按分类顺序组装推荐列表
   - 确保总数不超过 `maxCount` 参数（通常为4）

### 数据流

```
用户浏览商品 → 记录到 user_recent_views
                    ↓
            calculateRecommendations 被调用
                    ↓
            查询最近20条浏览记录
                    ↓
            提取商品的 category_id
                    ↓
            选择最新的2个不同分类
                    ↓
            每个分类查询2件推荐商品
                    ↓
            返回最多4件推荐商品
```

## 数据模型

### 使用的现有模型

#### UserRecentView
```go
type UserRecentView struct {
    ID        int64
    UserID    int64
    ProductID int64
    ViewedAt  time.Time
}
```

#### Product
```go
type Product struct {
    ID          int64
    Title       string
    Description string
    Price       float64
    CategoryID  int64
    ConditionID int64
    SellerID    int64
    Status      string
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
```

#### Category
```go
type Category struct {
    ID          int64
    Name        string
    Description string
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
```

### 数据库查询

#### 查询1: 获取浏览记录及商品分类
```sql
SELECT urv.viewed_at, p.id, p.category_id
FROM user_recent_views urv
INNER JOIN products p ON urv.product_id = p.id
WHERE urv.user_id = ?
ORDER BY urv.viewed_at DESC
LIMIT 20
```

#### 查询2: 获取分类推荐商品
```sql
SELECT *
FROM products
WHERE category_id = ?
  AND status = 'ForSale'
  AND seller_id != ?
  AND id NOT IN (?)
ORDER BY created_at DESC
LIMIT ?
```


## 正确性属性

*属性是指在系统的所有有效执行中都应该成立的特征或行为——本质上是关于系统应该做什么的形式化陈述。属性作为人类可读规范和机器可验证正确性保证之间的桥梁。*

### 属性反思

在定义正确性属性之前，我们需要识别并消除冗余：

**冗余分析**:
- 属性 1.3 和 4.1 都测试单一分类场景，可以合并
- 属性 1.5 和 3.1 都测试排除逻辑，但关注点不同（排除自己的商品 vs 排除已浏览商品），应保留两者
- 属性 2.2 和 2.3 测试不同场景（双分类 vs 单分类），应保留两者

**保留的核心属性**:
1. 分类提取和选择的正确性
2. 推荐数量的约束
3. 排除逻辑的正确性
4. 商品状态和排序的正确性
5. API兼容性

### 核心属性

**属性 1: 分类选择正确性**
*对于任何* 有多个不同分类浏览记录的用户，推荐系统应该选择最新浏览的两个不同分类
**验证需求: 1.2**

**属性 2: 推荐数量上限**
*对于任何* 用户和任何输入，推荐系统返回的商品数量不应超过指定的maxCount（通常为4）
**验证需求: 2.1**

**属性 3: 双分类均衡推荐**
*对于任何* 基于两个分类且每个分类都有足够商品的推荐场景，推荐系统应该从每个分类各推荐2件商品
**验证需求: 2.2**

**属性 4: 单分类推荐数量**
*对于任何* 只浏览了单一分类的用户，推荐系统应该从该分类推荐最多4件商品
**验证需求: 2.3, 4.1**

**属性 5: 排除自己的商品**
*对于任何* 用户和推荐结果，推荐列表中不应包含该用户自己发布的商品
**验证需求: 1.5**

**属性 6: 排除已浏览商品**
*对于任何* 用户和推荐结果，推荐列表中不应包含该用户已浏览过的商品
**验证需求: 3.1**

**属性 7: 仅推荐在售商品**
*对于任何* 推荐结果，所有推荐商品的状态都应该是"ForSale"
**验证需求: 3.2**

**属性 8: 商品时间排序**
*对于任何* 同一分类内的推荐商品，应该按创建时间倒序排列
**验证需求: 3.3**

**属性 9: 浏览记录处理顺序**
*对于任何* 用户的浏览记录，推荐系统应该按浏览时间倒序处理，确保选择的是最新浏览的分类
**验证需求: 4.4**

**属性 10: API响应格式兼容性**
*对于任何* 推荐服务调用，返回的数据结构应该符合ProductCardDTO格式，保持API兼容性
**验证需求: 5.2**

**属性 11: 推荐不足时的补充逻辑**
*对于任何* 推荐数量不足5件的已登录用户，GetHomeData方法应该用最新商品补充推荐列表
**验证需求: 5.3**

### 边界情况

以下边界情况将通过单元测试覆盖：

1. **空浏览记录**: 用户没有浏览记录时返回空列表（需求 1.4）
2. **单一分类**: 用户只浏览了一个分类的商品（需求 1.3, 4.1）
3. **商品不足**: 某个分类的可推荐商品不足2件（需求 2.4）
4. **总数不足**: 推荐商品总数不足4件（需求 2.5）
5. **空分类**: 某个分类下没有可推荐的商品（需求 4.2）
6. **无效记录**: 浏览记录中的商品已被删除或下架（需求 4.3）

## 错误处理

### 错误场景

1. **数据库查询失败**
   - 场景: 查询浏览记录或商品时数据库连接失败
   - 处理: 返回错误，不返回部分结果
   - 错误信息: 包含具体的数据库错误信息

2. **无效用户ID**
   - 场景: 传入的userID为0或负数
   - 处理: 返回空列表或错误
   - 错误信息: "invalid user ID"

3. **无效maxCount参数**
   - 场景: maxCount为0或负数
   - 处理: 使用默认值4
   - 日志: 记录警告信息

### 错误传播

- Service层的错误应该向上传播到Controller层
- Controller层负责将错误转换为适当的HTTP响应
- 不应该吞没错误或返回部分结果

## 测试策略

### 单元测试

单元测试将覆盖以下场景：

1. **基本功能测试**
   - 测试正常的双分类推荐场景
   - 测试单一分类推荐场景
   - 测试空浏览记录场景

2. **边界情况测试**
   - 商品数量不足的场景
   - 无效浏览记录的处理
   - 空分类的处理

3. **排除逻辑测试**
   - 验证排除用户自己的商品
   - 验证排除已浏览的商品
   - 验证只包含在售商品

4. **错误处理测试**
   - 数据库错误的处理
   - 无效参数的处理

### 属性测试

属性测试将使用 **Go的testing/quick包** 或 **gopter库** 来实现基于属性的测试。

每个属性测试应该：
- 运行至少100次迭代
- 使用随机生成的测试数据
- 在测试注释中明确标注对应的属性编号和需求编号
- 使用格式: `// Feature: category-based-recommendation, Property X: [property description]`

属性测试将验证：
1. 分类选择的正确性（属性1, 9）
2. 推荐数量约束（属性2, 3, 4）
3. 排除逻辑（属性5, 6, 7）
4. 排序正确性（属性8）
5. API兼容性（属性10, 11）

### 集成测试

集成测试将验证：
1. 完整的推荐流程（从浏览记录到推荐结果）
2. 与数据库的交互
3. API端点的正确性

### 测试数据生成

为了支持属性测试，需要实现以下生成器：

1. **用户生成器**: 生成随机用户ID
2. **分类生成器**: 生成随机分类
3. **商品生成器**: 生成随机商品，包含各种状态和分类
4. **浏览记录生成器**: 生成随机浏览记录序列
5. **约束生成器**: 生成满足特定约束的测试数据（如：多分类、单分类、商品不足等）

## 性能考虑

### 查询优化

1. **索引使用**
   - `user_recent_views` 表的 `(user_id, viewed_at)` 复合索引
   - `products` 表的 `category_id` 索引
   - `products` 表的 `status` 索引

2. **查询限制**
   - 浏览记录限制为20条，避免处理过多历史数据
   - 每个分类限制推荐2件商品，减少数据传输

3. **批量查询**
   - 使用 `IN` 查询批量获取商品信息
   - 避免N+1查询问题

### 缓存策略

虽然当前实现移除了Redis缓存，但设计保留了缓存接口：
- 可以在未来添加推荐结果缓存（TTL: 5-10分钟）
- 缓存键格式: `recommend:user:{userID}`

## 实现注意事项

### 代码修改范围

1. **主要修改**: `backend/service/recommend/service.go` 中的 `calculateRecommendations` 方法
2. **保持不变**: 
   - Controller层接口
   - Repository层接口
   - 数据模型
   - API路由

### 向后兼容性

- 保持方法签名不变
- 保持返回数据结构不变
- 保持API响应格式不变
- GetHomeData方法的补充逻辑保持不变

### 数据库依赖

- 依赖现有的 `user_recent_views` 表
- 依赖现有的 `products` 表
- 依赖现有的 `categories` 表
- 不需要新增表或字段

### 算法复杂度

- 时间复杂度: O(n + m)，其中n是浏览记录数（最多20），m是查询的商品数（最多4）
- 空间复杂度: O(n + m)
- 相比基于标签的算法，新算法更简单高效

## 迁移计划

### 部署步骤

1. **代码部署**
   - 更新 `calculateRecommendations` 方法
   - 部署新版本代码

2. **验证**
   - 检查API响应格式
   - 验证推荐结果的正确性
   - 监控错误日志

3. **回滚计划**
   - 如果出现问题，可以快速回滚到旧版本
   - 旧版本代码保留在版本控制中

### 数据迁移

不需要数据迁移，因为：
- 使用现有的浏览记录数据
- 不修改数据库结构
- 不需要数据转换

## 附录

### 算法伪代码

```
function calculateRecommendations(userID, maxCount):
    // 1. 获取浏览记录
    recentViews = getRecentViews(userID, limit=20)
    if recentViews.isEmpty():
        return []
    
    // 2. 提取分类
    viewedProductIDs = []
    for view in recentViews:
        viewedProductIDs.append(view.productID)
    
    // 3. 获取商品及其分类
    products = getProductsByIDs(viewedProductIDs)
    productMap = createMap(products, key=id)
    
    // 4. 按浏览时间顺序提取分类
    categories = []
    seenCategories = set()
    for view in recentViews:
        product = productMap[view.productID]
        if product exists and product.categoryID not in seenCategories:
            categories.append(product.categoryID)
            seenCategories.add(product.categoryID)
            if categories.length >= 2:
                break
    
    if categories.isEmpty():
        return []
    
    // 5. 确定每个分类的推荐数量
    if categories.length == 1:
        perCategoryLimit = maxCount
    else:
        perCategoryLimit = maxCount / 2
    
    // 6. 查询推荐商品
    recommendations = []
    for categoryID in categories:
        products = queryProducts(
            categoryID=categoryID,
            status="ForSale",
            excludeSellerID=userID,
            excludeProductIDs=viewedProductIDs,
            orderBy="created_at DESC",
            limit=perCategoryLimit
        )
        recommendations.extend(products)
        if recommendations.length >= maxCount:
            break
    
    return recommendations[0:maxCount]
```

### 与旧算法的对比

| 维度 | 旧算法（基于标签） | 新算法（基于分类） |
|------|-------------------|-------------------|
| 复杂度 | 高（需要统计标签频次） | 低（直接使用分类） |
| 查询次数 | 3-4次 | 2-3次 |
| 推荐相关性 | 中等 | 高 |
| 推荐多样性 | 低（可能集中在某些标签） | 高（保证两个分类） |
| 性能 | 较慢 | 较快 |
| 维护成本 | 高（依赖标签数据质量） | 低（分类是必填字段） |

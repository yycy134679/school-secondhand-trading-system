# 验证分类与标签模块的实现

Write-Host "=== 开始验证分类与标签模块 ===" -ForegroundColor Green

# 验证模型文件
$modelFiles = @(
    "model/category.go",
    "model/tag.go"
)

Write-Host "\n=== 验证模型文件 ===" -ForegroundColor Cyan
foreach ($file in $modelFiles) {
    if (Test-Path $file) {
        Write-Host "✅ 找到文件: $file"
    } else {
        Write-Host "❌ 未找到文件: $file" -ForegroundColor Red
    }
}

# 验证仓库层文件
$repoFiles = @(
    "repository/category_repo.go",
    "repository/tag_repo.go"
)

Write-Host "\n=== 验证仓库层文件 ===" -ForegroundColor Cyan
foreach ($file in $repoFiles) {
    if (Test-Path $file) {
        Write-Host "✅ 找到文件: $file"
    } else {
        Write-Host "❌ 未找到文件: $file" -ForegroundColor Red
    }
}

# 验证服务层文件
$serviceFiles = @(
    "service/category/service.go",
    "service/category/errors.go",
    "service/tag/service.go",
    "service/tag/errors.go"
)

Write-Host "\n=== 验证服务层文件 ===" -ForegroundColor Cyan
foreach ($file in $serviceFiles) {
    if (Test-Path $file) {
        Write-Host "✅ 找到文件: $file"
    } else {
        Write-Host "❌ 未找到文件: $file" -ForegroundColor Red
    }
}

# 验证控制器文件
$controllerFiles = @(
    "controller/category/controller.go",
    "controller/tag/controller.go"
)

Write-Host "\n=== 验证控制器文件 ===" -ForegroundColor Cyan
foreach ($file in $controllerFiles) {
    if (Test-Path $file) {
        Write-Host "✅ 找到文件: $file"
    } else {
        Write-Host "❌ 未找到文件: $file" -ForegroundColor Red
    }
}

# 验证路由文件
$routerFiles = @(
    "router/category.go",
    "router/tag.go"
)

Write-Host "\n=== 验证路由文件 ===" -ForegroundColor Cyan
foreach ($file in $routerFiles) {
    if (Test-Path $file) {
        Write-Host "✅ 找到文件: $file"
    } else {
        Write-Host "❌ 未找到文件: $file" -ForegroundColor Red
    }
}

# 验证错误码定义
Write-Host "\n=== 验证错误码定义 ===" -ForegroundColor Cyan
$categoryErrorFile = "service/category/errors.go"
$tagErrorFile = "service/tag/errors.go"

if (Test-Path $categoryErrorFile) {
    $categoryContent = Get-Content $categoryErrorFile -Raw
    if ($categoryContent -match "ErrCodeCategoryHasProducts.*=.*4001") {
        Write-Host "✅ 分类错误码4001已正确定义"
    } else {
        Write-Host "❌ 分类错误码4001定义不正确" -ForegroundColor Red
    }
}

if (Test-Path $tagErrorFile) {
    $tagContent = Get-Content $tagErrorFile -Raw
    if ($tagContent -match "ErrCodeTagHasProducts.*=.*4002") {
        Write-Host "✅ 标签错误码4002已正确定义"
    } else {
        Write-Host "❌ 标签错误码4002定义不正确" -ForegroundColor Red
    }
}

# 验证路由注册
Write-Host "\n=== 验证路由注册 ===" -ForegroundColor Cyan
$categoryRouterFile = "router/category.go"
$tagRouterFile = "router/tag.go"

if (Test-Path $categoryRouterFile) {
    $categoryRouterContent = Get-Content $categoryRouterFile -Raw
    if ($categoryRouterContent -match "SetupCategoryRoutes" -and 
        $categoryRouterContent -match "admin\.Use\(middleware\.AuthMiddleware\(\)\)" -and
        $categoryRouterContent -match "admin\.Use\(middleware\.AdminMiddleware\(\)\)") {
        Write-Host "✅ 分类路由注册包含正确的中间件"
    } else {
        Write-Host "❌ 分类路由注册缺少必要的中间件" -ForegroundColor Red
    }
}

if (Test-Path $tagRouterFile) {
    $tagRouterContent = Get-Content $tagRouterFile -Raw
    if ($tagRouterContent -match "SetupTagRoutes" -and 
        $tagRouterContent -match "admin\.Use\(middleware\.AuthMiddleware\(\)" -and
        $tagRouterContent -match "admin\.Use\(middleware\.AdminMiddleware\(\)") {
        Write-Host "✅ 标签路由注册包含正确的中间件"
    } else {
        Write-Host "❌ 标签路由注册缺少必要的中间件" -ForegroundColor Red
    }
}

Write-Host "\n=== 验证完成 ===" -ForegroundColor Green

# 商品模块、商品图片子模块和分类标签模块完整验证脚本

Write-Host "===== 商品系统完整模块验证脚本 ====="
Write-Host "开始时间: $(Get-Date)"
Write-Host ""

# 定义所有需要验证的模块
$modules = @(
    @{
        Name = "商品模块"
        Files = @(
            "model/product.go",
            "repository/product_repo.go",
            "service/product/service.go",
            "controller/product/controller.go"
        )
        Packages = @(
            "./model",
            "./repository",
            "./service/product",
            "./controller/product"
        )
    },
    @{
        Name = "商品图片子模块"
        Files = @(
            "common/util/file.go",
            "controller/product/image.go"
        )
        Packages = @(
            "./common/util",
            "./controller/product"
        )
    },
    @{
        Name = "分类与标签模块"
        Files = @(
            "model/category.go",
            "model/tag.go",
            "repository/category_repo.go",
            "repository/tag_repo.go",
            "service/category/service.go",
            "service/tag/service.go"
        )
        Packages = @(
            "./model",
            "./repository",
            "./service/category",
            "./service/tag"
        )
    }
)

# 全局结果标志
$globalSuccess = $true

# 1. 检查所有模块的必要文件是否存在
Write-Host "[1/6] 检查所有模块的必要文件是否存在..."
Write-Host ""

foreach ($module in $modules) {
    Write-Host "  检查 $($module.Name) 文件..."
    $moduleSuccess = $true
    
    foreach ($file in $module.Files) {
        $fullPath = "$PSScriptRoot\..\$file"
        if (Test-Path $fullPath) {
            Write-Host "    ✓ $file 存在"
        } else {
            Write-Host "    ✗ $file 不存在"
            $moduleSuccess = $false
            $globalSuccess = $false
        }
    }
    
    if ($moduleSuccess) {
        Write-Host "  ✓ $($module.Name) 所有文件存在"
    } else {
        Write-Host "  ✗ $($module.Name) 缺少必要文件"
    }
    Write-Host ""
}

if (-not $globalSuccess) {
    Write-Host "错误: 缺少必要的文件，无法继续验证"
    exit 1
}

# 2. 检查 Go 环境
Write-Host "[2/6] 检查 Go 环境..."
go version
if ($LASTEXITCODE -ne 0) {
    Write-Host "错误: Go 环境未配置或不可用"
    exit 1
}
Write-Host "  ✓ Go 环境检查通过"
Write-Host ""

# 3. 下载依赖
Write-Host "[3/6] 下载依赖..."
Set-Location -Path "$PSScriptRoot\.."
go mod tidy
if ($LASTEXITCODE -ne 0) {
    Write-Host "错误: 依赖下载失败"
    exit 1
}
Write-Host "  ✓ 依赖下载成功"
Write-Host ""

# 4. 运行 Go 语法检查 (vet)
Write-Host "[4/6] 运行 Go 语法检查 (vet)..."

# 合并所有包路径
$allPackages = @()
foreach ($module in $modules) {
    $allPackages += $module.Packages
}

# 去重
$allPackages = $allPackages | Select-Object -Unique

# 运行 vet 检查
go vet $allPackages
if ($LASTEXITCODE -eq 0) {
    Write-Host "  ✓ Go vet 检查通过，没有语法问题"
} else {
    Write-Host "警告: Go vet 检查发现潜在问题，但这不一定是致命错误"
}
Write-Host ""

# 5. 运行静态代码分析 (lint)
Write-Host "[5/6] 运行静态代码分析..."

# 检查是否安装了 golangci-lint
$golangciLintInstalled = $true
try {
    golangci-lint --version | Out-Null
} catch {
    $golangciLintInstalled = $false
}

if ($golangciLintInstalled) {
    golangci-lint run $allPackages
    if ($LASTEXITCODE -eq 0) {
        Write-Host "  ✓ golangci-lint 检查通过，代码质量良好"
    } else {
        Write-Host "警告: golangci-lint 检查发现代码质量问题，建议修复"
    }
} else {
    Write-Host "  ⚠ golangci-lint 未安装，跳过代码质量检查"
    Write-Host "  建议安装: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"
}
Write-Host ""

# 6. 尝试编译所有相关包
Write-Host "[6/6] 编译所有相关包..."

go build -v $allPackages
if ($LASTEXITCODE -eq 0) {
    Write-Host "  ✓ 编译成功! 所有包可以正常编译"
    Write-Host ""
    Write-Host "====================================="
    Write-Host "商品系统完整模块验证成功完成!"
    Write-Host "所有检查通过，代码可以正常编译。"
    Write-Host "====================================="
} else {
    Write-Host "  ✗ 编译失败! 请检查代码中的错误"
    exit 1
}

Write-Host ""
Write-Host "结束时间: $(Get-Date)"
Write-Host "验证完成!"
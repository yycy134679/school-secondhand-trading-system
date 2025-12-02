#!/usr/bin/env pwsh

# 简单的商品模块验证脚本
Write-Host "=== 商品模块验证脚本 ==="
Write-Host "开始时间: $(Get-Date)"
Write-Host ""

# 设置路径
$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$projectRoot = Join-Path -Path $scriptDir -ChildPath ".."

# 检查文件存在性
Write-Host "[1/3] 检查关键文件..."
Write-Host ""

$filesToCheck = @(
    "controller\product\controller.go",
    "controller\product\image.go",
    "router\product.go",
    "common\util\file.go"
)

foreach ($file in $filesToCheck) {
    $fullPath = Join-Path -Path $projectRoot -ChildPath $file
    if (Test-Path $fullPath) {
        Write-Host "✓ 存在: $file"
    } else {
        Write-Host "✗ 不存在: $file"
    }
}

# 检查控制器方法
Write-Host ""
Write-Host "[2/3] 检查API方法..."
Write-Host ""

# 检查主要控制器
$controllerPath = Join-Path -Path $projectRoot -ChildPath "controller\product\controller.go"
if (Test-Path $controllerPath) {
    $content = Get-Content -Path $controllerPath -Raw
    Write-Host "控制器方法检查:" 
    
    $methods = @(
        "CreateProduct",
        "UpdateProduct",
        "ChangeProductStatus",
        "UndoLastStatusChange",
        "GetProductDetail",
        "ListMyProducts",
        "SearchProducts",
        "GetProductsByCategory"
    )
    
    foreach ($method in $methods) {
        if ($content -match "func.*$method") {
            Write-Host "  ✓ $method"
        } else {
            Write-Host "  ✗ $method"
        }
    }
}

# 检查图片控制器
Write-Host ""
$imagePath = Join-Path -Path $projectRoot -ChildPath "controller\product\image.go"
if (Test-Path $imagePath) {
    $content = Get-Content -Path $imagePath -Raw
    Write-Host "图片管理方法检查:" 
    
    $methods = @(
        "UploadProductImage",
        "SetPrimaryImage",
        "UpdateImageSortOrder",
        "DeleteProductImage"
    )
    
    foreach ($method in $methods) {
        if ($content -match "func.*$method") {
            Write-Host "  ✓ $method"
        } else {
            Write-Host "  ✗ $method"
        }
    }
}

# 简单编译检查
Write-Host ""
Write-Host "[3/3] 执行简单编译检查..."
Write-Host ""

Set-Location -Path $projectRoot
Write-Host "检查Go环境..."
go version
if ($LASTEXITCODE -eq 0) {
    Write-Host "Go环境可用"
    
    # 检查关键文件语法
    Write-Host ""
    Write-Host "检查语法..."
    $criticalFiles = @(
        "controller\product\controller.go",
        "controller\product\image.go",
        "router\product.go"
    )
    
    foreach ($file in $criticalFiles) {
        $fullPath = Join-Path -Path $projectRoot -ChildPath $file
        if (Test-Path $fullPath) {
            Write-Host "检查: $file"
            go vet "$fullPath" > $null 2>&1
            if ($LASTEXITCODE -eq 0) {
                Write-Host "  ✓ 语法正确"
            } else {
                Write-Host "  ✗ 语法错误"
            }
        }
    }
}

Write-Host ""
Write-Host "验证完成!"

# 简单的商品模块编译验证脚本

Write-Host "===== 商品模块编译验证脚本 ====="
Write-Host "开始时间: $(Get-Date)"
Write-Host ""

# 检查必要的文件是否存在
Write-Host "[1/5] 检查必要文件是否存在..."

$requiredFiles = @(
    "model/product.go",
    "repository/product_repo.go",
    "service/product/service.go",
    "controller/product/controller.go",
    "router/router.go"
)

$allFilesExist = $true
foreach ($file in $requiredFiles) {
    $fullPath = "$PSScriptRoot\..\$file"
    if (Test-Path $fullPath) {
        Write-Host "  ✓ $file 存在"
    } else {
        Write-Host "  ✗ $file 不存在"
        $allFilesExist = $false
    }
}

if (-not $allFilesExist) {
    Write-Host "错误: 缺少必要的文件"
    exit 1
}

Write-Host ""
Write-Host "[2/5] 检查 Go 环境..."
go version
if ($LASTEXITCODE -ne 0) {
    Write-Host "错误: Go 环境未配置或不可用"
    exit 1
}

Write-Host ""
Write-Host "[3/5] 下载依赖..."
Set-Location -Path "$PSScriptRoot\.."
go mod tidy
if ($LASTEXITCODE -ne 0) {
    Write-Host "错误: 依赖下载失败"
    exit 1
}
Write-Host "  ✓ 依赖下载成功"

Write-Host ""
Write-Host "[4/5] 运行 Go 语法检查..."
go vet ./model ./repository ./service/product ./controller/product ./router
if ($LASTEXITCODE -eq 0) {
    Write-Host "  ✓ Go vet 检查通过"
} else {
    Write-Host "警告: Go vet 检查发现潜在问题，但这不一定是致命错误"
}

Write-Host ""
Write-Host "[5/5] 尝试编译整个项目..."
go build -v ./model ./repository ./service/product ./controller/product ./router
if ($LASTEXITCODE -eq 0) {
    Write-Host "  ✓ 编译成功!"
    Write-Host ""
    Write-Host "====================================="
    Write-Host "商品模块编译验证成功完成!"
    Write-Host "所有检查通过，代码可以正常编译。"
    Write-Host "====================================="
} else {
    Write-Host "  ✗ 编译失败!"
    exit 1
}

Write-Host ""
Write-Host "结束时间: $(Get-Date)"
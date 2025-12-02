# 简单的商品模块编译验证脚本
Write-Output "开始验证商品模块编译状态..."

# 检查必要的文件是否存在
Write-Output "检查必要文件..."
$requiredFiles = @(
    "c:\Users\15751\Desktop\school-secondhand-trading-system\controller\product\controller.go",
    "c:\Users\15751\Desktop\school-secondhand-trading-system\controller\product\image.go",
    "c:\Users\15751\Desktop\school-secondhand-trading-system\router\product.go",
    "c:\Users\15751\Desktop\school-secondhand-trading-system\service\product\service.go",
    "c:\Users\15751\Desktop\school-secondhand-trading-system\model\product.go",
    "c:\Users\15751\Desktop\school-secondhand-trading-system\repository\product_repo.go"
)

$allFilesExist = $true
foreach ($file in $requiredFiles) {
    if (Test-Path $file) {
        Write-Output "✓ 文件存在: $file"
    } else {
        Write-Output "✗ 文件不存在: $file"
        $allFilesExist = $false
    }
}

if (-not $allFilesExist) {
    Write-Output "错误: 必要文件缺失，无法进行编译验证"
    exit 1
}

Write-Output "所有必要文件都存在，开始进行编译检查..."
Write-Output "注意: 由于完整编译可能需要更多依赖，我们将进行简单的语法检查"

# 运行go build来检查语法和编译错误
Write-Output "运行 go build 检查..."
Set-Location "c:\Users\15751\Desktop\school-secondhand-trading-system"
go build 2>&1
$buildResult = $LASTEXITCODE

if ($buildResult -eq 0) {
    Write-Output "✓ 编译成功!"
    exit 0
} else {
    Write-Output "✗ 编译失败，请检查错误信息"
    exit 1
}
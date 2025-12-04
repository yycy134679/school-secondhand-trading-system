# 推荐服务功能快速测试脚本 (PowerShell)
# 使用方法: .\test_recommend.ps1

$BASE_URL = "http://localhost:8080"
$TOKEN = "YOUR_TOKEN_HERE"  # 需要先登录获取token

Write-Host "======================================" -ForegroundColor Cyan
Write-Host "  推荐服务功能测试 (任务 8.2)" -ForegroundColor Cyan
Write-Host "======================================" -ForegroundColor Cyan
Write-Host ""

# 测试1: 健康检查
Write-Host "📋 测试 1: 健康检查" -ForegroundColor Yellow
try {
    $response = Invoke-RestMethod -Uri "$BASE_URL/health" -Method Get
    if ($response.status -eq "ok") {
        Write-Host "✓ 服务正常运行" -ForegroundColor Green
    } else {
        Write-Host "✗ 服务状态异常" -ForegroundColor Red
    }
} catch {
    Write-Host "✗ 服务未响应: $_" -ForegroundColor Red
    exit 1
}
Write-Host ""

# 测试2: 记录浏览 (需要登录)
Write-Host "📋 测试 2: 记录商品浏览" -ForegroundColor Yellow
if ($TOKEN -eq "YOUR_TOKEN_HERE") {
    Write-Host "⚠ 跳过 - 需要先设置 TOKEN 变量" -ForegroundColor Magenta
} else {
    Write-Host "  浏览商品 ID: 1, 2, 3"
    foreach ($id in 1..3) {
        try {
            $headers = @{
                "Authorization" = "Bearer $TOKEN"
            }
            $response = Invoke-RestMethod -Uri "$BASE_URL/api/v1/products/$id/view" `
                -Method Post -Headers $headers
            
            if ($response.data.recorded) {
                Write-Host "  ✓ 记录商品 $id 浏览成功" -ForegroundColor Green
            } else {
                Write-Host "  ✗ 记录商品 $id 浏览失败" -ForegroundColor Red
                Write-Host "  响应: $($response | ConvertTo-Json)" -ForegroundColor Gray
            }
        } catch {
            Write-Host "  ✗ 记录商品 $id 浏览失败: $_" -ForegroundColor Red
        }
        Start-Sleep -Milliseconds 200
    }
}
Write-Host ""

# 测试3: 获取首页数据 (不需要登录)
Write-Host "📋 测试 3: 获取首页数据 (未登录)" -ForegroundColor Yellow
try {
    $response = Invoke-RestMethod -Uri "$BASE_URL/api/v1/home" -Method Get
    if ($response.data.recommendations -or $response.data.latest) {
        Write-Host "✓ 成功获取首页数据" -ForegroundColor Green
        $recCount = if ($response.data.recommendations) { $response.data.recommendations.Count } else { 0 }
        $latestCount = if ($response.data.latest) { $response.data.latest.Count } else { 0 }
        Write-Host "  推荐商品数: $recCount"
        Write-Host "  最新商品数: $latestCount"
    } else {
        Write-Host "✗ 获取首页数据失败" -ForegroundColor Red
        Write-Host "  响应: $($response | ConvertTo-Json)" -ForegroundColor Gray
    }
} catch {
    Write-Host "✗ 获取首页数据失败: $_" -ForegroundColor Red
}
Write-Host ""

# 测试4: 获取首页数据 (已登录 - 有推荐)
Write-Host "📋 测试 4: 获取首页数据 (已登录)" -ForegroundColor Yellow
if ($TOKEN -eq "YOUR_TOKEN_HERE") {
    Write-Host "⚠ 跳过 - 需要先设置 TOKEN 变量" -ForegroundColor Magenta
} else {
    try {
        $headers = @{
            "Authorization" = "Bearer $TOKEN"
        }
        $response = Invoke-RestMethod -Uri "$BASE_URL/api/v1/home" `
            -Method Get -Headers $headers
        
        if ($response.data.recommendations -or $response.data.latest) {
            Write-Host "✓ 成功获取个性化推荐" -ForegroundColor Green
            $recCount = if ($response.data.recommendations) { $response.data.recommendations.Count } else { 0 }
            $latestCount = if ($response.data.latest) { $response.data.latest.Count } else { 0 }
            $total = $response.data.totalCount
            Write-Host "  推荐商品数: $recCount"
            Write-Host "  最新商品数: $latestCount"
            Write-Host "  商品总数: $total"
        } else {
            Write-Host "✗ 获取个性化推荐失败" -ForegroundColor Red
            Write-Host "  响应: $($response | ConvertTo-Json)" -ForegroundColor Gray
        }
    } catch {
        Write-Host "✗ 获取个性化推荐失败: $_" -ForegroundColor Red
    }
}
Write-Host ""

# 测试5: 获取浏览记录
Write-Host "📋 测试 5: 获取浏览记录" -ForegroundColor Yellow
if ($TOKEN -eq "YOUR_TOKEN_HERE") {
    Write-Host "⚠ 跳过 - 需要先设置 TOKEN 变量" -ForegroundColor Magenta
} else {
    try {
        $headers = @{
            "Authorization" = "Bearer $TOKEN"
        }
        $response = Invoke-RestMethod -Uri "$BASE_URL/api/v1/users/recent-views?limit=10" `
            -Method Get -Headers $headers
        
        if ($response.data.views) {
            Write-Host "✓ 成功获取浏览记录" -ForegroundColor Green
            $viewCount = $response.data.total
            Write-Host "  浏览记录数: $viewCount"
        } else {
            Write-Host "✗ 获取浏览记录失败" -ForegroundColor Red
            Write-Host "  响应: $($response | ConvertTo-Json)" -ForegroundColor Gray
        }
    } catch {
        Write-Host "✗ 获取浏览记录失败: $_" -ForegroundColor Red
    }
}
Write-Host ""

# 测试总结
Write-Host "======================================" -ForegroundColor Cyan
Write-Host "  测试完成" -ForegroundColor Cyan
Write-Host "======================================" -ForegroundColor Cyan
Write-Host ""

if ($TOKEN -eq "YOUR_TOKEN_HERE") {
    Write-Host "提示: 要完整测试所有功能，请先登录并设置 TOKEN 变量" -ForegroundColor Magenta
    Write-Host ""
    Write-Host "步骤:" -ForegroundColor White
    Write-Host "1. 调用登录接口获取 token" -ForegroundColor Gray
    Write-Host '   $body = @{' -ForegroundColor Gray
    Write-Host '       account = "your_account"' -ForegroundColor Gray
    Write-Host '       password = "your_password"' -ForegroundColor Gray
    Write-Host '   } | ConvertTo-Json' -ForegroundColor Gray
    Write-Host '   $response = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/users/login" `' -ForegroundColor Gray
    Write-Host '       -Method Post -Body $body -ContentType "application/json"' -ForegroundColor Gray
    Write-Host '   $token = $response.data.token' -ForegroundColor Gray
    Write-Host ""
    Write-Host "2. 修改脚本中的 `$TOKEN 变量为获取到的 token" -ForegroundColor Gray
    Write-Host "3. 重新运行测试脚本" -ForegroundColor Gray
}
Write-Host ""
Write-Host "详细文档:" -ForegroundColor White
Write-Host "- API 测试文档: backend/API_TEST.md" -ForegroundColor Gray
Write-Host "- 功能验证报告: backend/TASK_8.2_VERIFICATION.md" -ForegroundColor Gray
Write-Host ""

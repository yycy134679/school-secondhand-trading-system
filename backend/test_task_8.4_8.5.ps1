# 任务 8.4 & 8.5 测试脚本
# 测试首页数据和浏览记录接口

Write-Host "=================================" -ForegroundColor Cyan
Write-Host "任务 8.4 & 8.5 接口测试" -ForegroundColor Cyan
Write-Host "=================================" -ForegroundColor Cyan

$baseUrl = "http://localhost:8080"
$token = "" # 在这里填入你的访问令牌

# ========================================
# 任务 8.4 测试：首页数据接口
# ========================================

Write-Host "`n【任务 8.4】首页数据接口测试" -ForegroundColor Green
Write-Host "=====================================" -ForegroundColor Green

# 测试 1: 未登录访问首页
Write-Host "`n测试 1: 未登录访问首页" -ForegroundColor Yellow
Write-Host "-------------------------------" -ForegroundColor Gray
try {
    $response = Invoke-RestMethod -Uri "$baseUrl/api/v1/home" -Method Get
    Write-Host "✓ 状态码: $($response.code)" -ForegroundColor Green
    Write-Host "  - 推荐数量: $($response.data.recommendations.Count)" -ForegroundColor White
    Write-Host "  - 最新商品数量: $($response.data.latest.Count)" -ForegroundColor White
    Write-Host "  - 总商品数: $($response.data.totalCount)" -ForegroundColor White
    
    if ($response.data.recommendations.Count -eq 0) {
        Write-Host "✓ 未登录用户推荐列表为空（正确）" -ForegroundColor Green
    }
} catch {
    Write-Host "✗ 请求失败: $_" -ForegroundColor Red
}

# 测试 2: 分页查询
Write-Host "`n测试 2: 分页查询 (page=1, pageSize=5)" -ForegroundColor Yellow
Write-Host "-------------------------------" -ForegroundColor Gray
try {
    $response = Invoke-RestMethod -Uri "$baseUrl/api/v1/home?page=1&pageSize=5" -Method Get
    Write-Host "✓ 状态码: $($response.code)" -ForegroundColor Green
    Write-Host "  - 返回最新商品数: $($response.data.latest.Count)" -ForegroundColor White
    Write-Host "  - 总商品数: $($response.data.totalCount)" -ForegroundColor White
} catch {
    Write-Host "✗ 请求失败: $_" -ForegroundColor Red
}

# 如果有 token，测试已登录场景
if ($token -ne "") {
    Write-Host "`n测试 3: 已登录访问首页" -ForegroundColor Yellow
    Write-Host "-------------------------------" -ForegroundColor Gray
    try {
        $headers = @{ "Authorization" = "Bearer $token" }
        $response = Invoke-RestMethod -Uri "$baseUrl/api/v1/home" -Method Get -Headers $headers
        Write-Host "✓ 状态码: $($response.code)" -ForegroundColor Green
        Write-Host "  - 推荐数量: $($response.data.recommendations.Count)" -ForegroundColor White
        Write-Host "  - 最新商品数量: $($response.data.latest.Count)" -ForegroundColor White
        
        if ($response.data.recommendations.Count -gt 0) {
            Write-Host "`n  推荐商品列表:" -ForegroundColor Cyan
            foreach ($product in $response.data.recommendations) {
                Write-Host "    - [$($product.id)] $($product.title) - ¥$($product.price)" -ForegroundColor White
            }
        }
    } catch {
        Write-Host "✗ 请求失败: $_" -ForegroundColor Red
    }
}

# ========================================
# 任务 8.5 测试：浏览记录接口
# ========================================

Write-Host "`n`n【任务 8.5】浏览记录接口测试" -ForegroundColor Green
Write-Host "=====================================" -ForegroundColor Green

# 测试 4: 未登录获取浏览记录（应该失败）
Write-Host "`n测试 4: 未登录获取浏览记录（应返回 401）" -ForegroundColor Yellow
Write-Host "-------------------------------" -ForegroundColor Gray
try {
    $response = Invoke-RestMethod -Uri "$baseUrl/api/v1/users/recent-views" -Method Get
    Write-Host "✗ 错误：应该返回 401 未授权" -ForegroundColor Red
} catch {
    $statusCode = $_.Exception.Response.StatusCode.value__
    if ($statusCode -eq 401) {
        Write-Host "✓ 正确返回 401 未授权" -ForegroundColor Green
    } else {
        Write-Host "✗ 返回状态码: $statusCode（期望 401）" -ForegroundColor Red
    }
}

# 如果有 token，测试已登录场景
if ($token -ne "") {
    # 测试 5: 记录一些浏览行为
    Write-Host "`n测试 5: 记录浏览行为" -ForegroundColor Yellow
    Write-Host "-------------------------------" -ForegroundColor Gray
    $headers = @{ "Authorization" = "Bearer $token" }
    $productIds = @(1, 2, 3, 5, 7)
    
    foreach ($id in $productIds) {
        try {
            $response = Invoke-RestMethod -Uri "$baseUrl/api/v1/products/$id/view" -Method Post -Headers $headers
            if ($response.data.recorded) {
                Write-Host "  ✓ 记录浏览商品 $id" -ForegroundColor Green
            }
            Start-Sleep -Milliseconds 200 # 短暂延迟
        } catch {
            Write-Host "  ✗ 记录商品 $id 失败" -ForegroundColor Red
        }
    }
    
    # 测试 6: 获取浏览记录
    Write-Host "`n测试 6: 获取浏览记录" -ForegroundColor Yellow
    Write-Host "-------------------------------" -ForegroundColor Gray
    try {
        $response = Invoke-RestMethod -Uri "$baseUrl/api/v1/users/recent-views" -Method Get -Headers $headers
        Write-Host "✓ 状态码: $($response.code)" -ForegroundColor Green
        Write-Host "  - 浏览记录总数: $($response.data.total)" -ForegroundColor White
        
        if ($response.data.views.Count -gt 0) {
            Write-Host "`n  浏览记录列表:" -ForegroundColor Cyan
            foreach ($view in $response.data.views) {
                $viewTime = ([DateTime]$view.viewedAt).ToString("HH:mm:ss")
                Write-Host "    - [$viewTime] $($view.product.title) - ¥$($view.product.price) ($($view.product.status))" -ForegroundColor White
            }
        } else {
            Write-Host "  (无浏览记录)" -ForegroundColor Gray
        }
    } catch {
        Write-Host "✗ 请求失败: $_" -ForegroundColor Red
    }
    
    # 测试 7: 限制返回数量
    Write-Host "`n测试 7: 限制返回数量 (limit=3)" -ForegroundColor Yellow
    Write-Host "-------------------------------" -ForegroundColor Gray
    try {
        $response = Invoke-RestMethod -Uri "$baseUrl/api/v1/users/recent-views?limit=3" -Method Get -Headers $headers
        Write-Host "✓ 状态码: $($response.code)" -ForegroundColor Green
        Write-Host "  - 返回记录数: $($response.data.total)" -ForegroundColor White
        
        if ($response.data.total -le 3) {
            Write-Host "✓ 限制生效（返回 ≤ 3 条）" -ForegroundColor Green
        } else {
            Write-Host "✗ 限制未生效（返回 > 3 条）" -ForegroundColor Red
        }
    } catch {
        Write-Host "✗ 请求失败: $_" -ForegroundColor Red
    }
    
    # 测试 8: 验证推荐效果
    Write-Host "`n测试 8: 浏览后再次访问首页（验证推荐）" -ForegroundColor Yellow
    Write-Host "-------------------------------" -ForegroundColor Gray
    try {
        $response = Invoke-RestMethod -Uri "$baseUrl/api/v1/home" -Method Get -Headers $headers
        Write-Host "✓ 状态码: $($response.code)" -ForegroundColor Green
        Write-Host "  - 推荐数量: $($response.data.recommendations.Count)" -ForegroundColor White
        
        if ($response.data.recommendations.Count -gt 0) {
            Write-Host "`n  个性化推荐列表:" -ForegroundColor Cyan
            foreach ($product in $response.data.recommendations) {
                Write-Host "    - [$($product.id)] $($product.title) - ¥$($product.price)" -ForegroundColor White
            }
            Write-Host "✓ 推荐系统生效" -ForegroundColor Green
        } else {
            Write-Host "  (暂无推荐商品)" -ForegroundColor Gray
        }
    } catch {
        Write-Host "✗ 请求失败: $_" -ForegroundColor Red
    }
} else {
    Write-Host "`n⚠ 跳过需要登录的测试（未提供 token）" -ForegroundColor Yellow
    Write-Host "  请在脚本开头的 `$token 变量中填入访问令牌" -ForegroundColor Gray
}

# ========================================
# 测试总结
# ========================================

Write-Host "`n`n=================================" -ForegroundColor Cyan
Write-Host "测试完成！" -ForegroundColor Cyan
Write-Host "=================================" -ForegroundColor Cyan

Write-Host "`n测试覆盖项目:" -ForegroundColor White
Write-Host "  ✓ 任务 8.4: GET /api/v1/home" -ForegroundColor Green
Write-Host "    - 未登录访问" -ForegroundColor Gray
Write-Host "    - 分页参数" -ForegroundColor Gray
if ($token -ne "") {
    Write-Host "    - 已登录访问（个性化推荐）" -ForegroundColor Gray
}
Write-Host "  ✓ 任务 8.5: GET /api/v1/users/recent-views" -ForegroundColor Green
Write-Host "    - 未登录访问（401）" -ForegroundColor Gray
if ($token -ne "") {
    Write-Host "    - 记录浏览行为" -ForegroundColor Gray
    Write-Host "    - 获取浏览记录" -ForegroundColor Gray
    Write-Host "    - limit 参数限制" -ForegroundColor Gray
    Write-Host "    - 推荐效果验证" -ForegroundColor Gray
}

Write-Host "`n如需完整测试，请:" -ForegroundColor Yellow
Write-Host "  1. 先登录获取 token" -ForegroundColor White
Write-Host "  2. 在脚本开头设置 `$token 变量" -ForegroundColor White
Write-Host "  3. 重新运行脚本" -ForegroundColor White

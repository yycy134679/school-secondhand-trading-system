#!/bin/bash
# 推荐服务功能快速测试脚本

BASE_URL="http://localhost:8080"
TOKEN="YOUR_TOKEN_HERE"  # 需要先登录获取token

echo "======================================"
echo "  推荐服务功能测试 (任务 8.2)"
echo "======================================"
echo ""

# 颜色定义
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 测试1: 健康检查
echo "📋 测试 1: 健康检查"
response=$(curl -s "$BASE_URL/health")
if echo "$response" | grep -q "ok"; then
    echo -e "${GREEN}✓ 服务正常运行${NC}"
else
    echo -e "${RED}✗ 服务未响应${NC}"
    exit 1
fi
echo ""

# 测试2: 记录浏览 (需要登录)
echo "📋 测试 2: 记录商品浏览"
if [ "$TOKEN" = "YOUR_TOKEN_HERE" ]; then
    echo -e "${YELLOW}⚠ 跳过 - 需要先设置 TOKEN 变量${NC}"
else
    echo "  浏览商品 ID: 1, 2, 3"
    for id in 1 2 3; do
        response=$(curl -s -X POST "$BASE_URL/api/v1/products/$id/view" \
            -H "Authorization: Bearer $TOKEN")
        if echo "$response" | grep -q "recorded"; then
            echo -e "  ${GREEN}✓ 记录商品 $id 浏览成功${NC}"
        else
            echo -e "  ${RED}✗ 记录商品 $id 浏览失败${NC}"
            echo "  响应: $response"
        fi
        sleep 0.2
    done
fi
echo ""

# 测试3: 获取首页数据 (不需要登录)
echo "📋 测试 3: 获取首页数据 (未登录)"
response=$(curl -s "$BASE_URL/api/v1/home")
if echo "$response" | grep -q "recommendations"; then
    echo -e "${GREEN}✓ 成功获取首页数据${NC}"
    # 解析数据
    rec_count=$(echo "$response" | jq -r '.data.recommendations | length' 2>/dev/null || echo "N/A")
    latest_count=$(echo "$response" | jq -r '.data.latest | length' 2>/dev/null || echo "N/A")
    echo "  推荐商品数: $rec_count"
    echo "  最新商品数: $latest_count"
else
    echo -e "${RED}✗ 获取首页数据失败${NC}"
    echo "  响应: $response"
fi
echo ""

# 测试4: 获取首页数据 (已登录 - 有推荐)
echo "📋 测试 4: 获取首页数据 (已登录)"
if [ "$TOKEN" = "YOUR_TOKEN_HERE" ]; then
    echo -e "${YELLOW}⚠ 跳过 - 需要先设置 TOKEN 变量${NC}"
else
    response=$(curl -s "$BASE_URL/api/v1/home" \
        -H "Authorization: Bearer $TOKEN")
    if echo "$response" | grep -q "recommendations"; then
        echo -e "${GREEN}✓ 成功获取个性化推荐${NC}"
        rec_count=$(echo "$response" | jq -r '.data.recommendations | length' 2>/dev/null || echo "N/A")
        latest_count=$(echo "$response" | jq -r '.data.latest | length' 2>/dev/null || echo "N/A")
        total=$(echo "$response" | jq -r '.data.totalCount' 2>/dev/null || echo "N/A")
        echo "  推荐商品数: $rec_count"
        echo "  最新商品数: $latest_count"
        echo "  商品总数: $total"
    else
        echo -e "${RED}✗ 获取个性化推荐失败${NC}"
        echo "  响应: $response"
    fi
fi
echo ""

# 测试5: 获取浏览记录
echo "📋 测试 5: 获取浏览记录"
if [ "$TOKEN" = "YOUR_TOKEN_HERE" ]; then
    echo -e "${YELLOW}⚠ 跳过 - 需要先设置 TOKEN 变量${NC}"
else
    response=$(curl -s "$BASE_URL/api/v1/users/recent-views?limit=10" \
        -H "Authorization: Bearer $TOKEN")
    if echo "$response" | grep -q "views"; then
        echo -e "${GREEN}✓ 成功获取浏览记录${NC}"
        view_count=$(echo "$response" | jq -r '.data.total' 2>/dev/null || echo "N/A")
        echo "  浏览记录数: $view_count"
    else
        echo -e "${RED}✗ 获取浏览记录失败${NC}"
        echo "  响应: $response"
    fi
fi
echo ""

# 测试总结
echo "======================================"
echo "  测试完成"
echo "======================================"
echo ""
if [ "$TOKEN" = "YOUR_TOKEN_HERE" ]; then
    echo -e "${YELLOW}提示: 要完整测试所有功能，请先登录并设置 TOKEN 变量${NC}"
    echo ""
    echo "步骤:"
    echo "1. 调用登录接口获取 token"
    echo "   curl -X POST $BASE_URL/api/v1/users/login \\"
    echo "     -H 'Content-Type: application/json' \\"
    echo "     -d '{\"account\":\"your_account\",\"password\":\"your_password\"}'"
    echo ""
    echo "2. 修改脚本中的 TOKEN 变量"
    echo "3. 重新运行测试脚本"
fi
echo ""
echo "详细文档:"
echo "- API 测试文档: backend/API_TEST.md"
echo "- 功能验证报告: backend/TASK_8.2_VERIFICATION.md"
echo ""

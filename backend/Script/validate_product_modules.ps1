#!/usr/bin/env pwsh

Write-Host "===== School Second-hand Trading System Backend Validation Script =====" -ForegroundColor Cyan
Write-Host "Validating: Product Module, Category & Tag Module Implementation against Requirements"
Write-Host "`n"

# Define project root directory
$backendRoot = "c:\Users\15751\Desktop\school-secondhand-trading-system\backend"

# 1. Check file existence
Write-Host "1. Checking required files existence" -ForegroundColor Green

$requiredFiles = @(
    @{"Path" = "$backendRoot\model\product.go"; "Module" = "Product Model"}
    @{"Path" = "$backendRoot\repository\product_repo.go"; "Module" = "Product Repository"}
    @{"Path" = "$backendRoot\service\product\service.go"; "Module" = "Product Service"}
    @{"Path" = "$backendRoot\controller\product\controller.go"; "Module" = "Product Controller"}
    @{"Path" = "$backendRoot\common\util\file.go"; "Module" = "File Storage Service"}
    @{"Path" = "$backendRoot\controller\product\image.go"; "Module" = "Product Image Controller"}
    @{"Path" = "$backendRoot\model\category.go"; "Module" = "Category Model"}
    @{"Path" = "$backendRoot\repository\category_repo.go"; "Module" = "Category Repository"}
    @{"Path" = "$backendRoot\model\tag.go"; "Module" = "Tag Model"}
    @{"Path" = "$backendRoot\repository\tag_repo.go"; "Module" = "Tag Repository"}
    @{"Path" = "$backendRoot\router\router.go"; "Module" = "Router Configuration"}
)

$fileFailures = 0
foreach ($file in $requiredFiles) {
    if (Test-Path $file.Path) {
        Write-Host "✓ File exists: $($file.Path)" -ForegroundColor Green
    } else {
        Write-Host "✗ File missing: $($file.Path)" -ForegroundColor Red
        $fileFailures += 1
    }
}

# 2. Check product API endpoints
Write-Host "`n2. Product API Endpoints" -ForegroundColor Green
Write-Host "The following endpoints should be implemented:"
$productEndpoints = @(
    "POST /api/v1/products (Authenticated)"
    "PUT /api/v1/products/:id (Authenticated)"
    "POST /api/v1/products/:id/status (Authenticated)"
    "POST /api/v1/products/:id/status/undo (Authenticated)"
    "GET /api/v1/products/:id (Public)"
    "GET /api/v1/products/my (Authenticated)"
    "GET /api/v1/products/search (Public)"
    "GET /api/v1/products/category/:categoryId (Public)"
)
foreach ($endpoint in $productEndpoints) {
    Write-Host "   - $endpoint"
}

# 3. Check image management API endpoints
Write-Host "`n3. Product Image Management API Endpoints" -ForegroundColor Green
Write-Host "The following endpoints should be implemented:"
$imageEndpoints = @(
    "POST /api/v1/products/:id/images (Authenticated)"
    "PUT /api/v1/products/:id/images/:imageId/primary (Authenticated)"
    "PATCH /api/v1/products/:id/images/:imageId (Authenticated)"
    "DELETE /api/v1/products/:id/images/:imageId (Authenticated)"
)
foreach ($endpoint in $imageEndpoints) {
    Write-Host "   - $endpoint"
}

# 4. Check category and tag API endpoints
Write-Host "`n4. Category & Tag API Endpoints" -ForegroundColor Green
Write-Host "The following endpoints should be implemented:"
$categoryTagEndpoints = @(
    "GET /api/v1/categories (Public)"
    "POST /api/v1/admin/categories (Admin)"
    "PUT /api/v1/admin/categories/:id (Admin)"
    "DELETE /api/v1/admin/categories/:id (Admin)"
    "GET /api/v1/tags (Public)"
    "POST /api/v1/admin/tags (Admin)"
    "PUT /api/v1/admin/tags/:id (Admin)"
    "DELETE /api/v1/admin/tags/:id (Admin)"
)
foreach ($endpoint in $categoryTagEndpoints) {
    Write-Host "   - $endpoint"
}

# 5. Check core functionality requirements
Write-Host "`n5. Core Functionality Requirements" -ForegroundColor Green
Write-Host "Key implementations to verify in code review:"
$coreRequirements = @(
    "ProductService.CreateProduct: Field validation, category/tag/condition validation, image upload, status setting"
    "ProductService.UpdateProduct: Permission check, status restrictions, image replacement/addition"
    "ProductService.ChangeStatus: Status transition logic, Redis cache for undo records"
    "ProductService.UndoLastStatusChange: Redis cache validation, status rollback logic"
    "FileService.SaveImage: Image format validation, size limit (2MB), URL generation"
    "CategoryService.DeleteCategory: Check product references, return error code 4001"
    "TagService.DeleteTag: Check product references, return error code 4002"
)
foreach ($req in $coreRequirements) {
    Write-Host "   - $req"
}

# 6. Summary report
Write-Host "`n===== Validation Summary =====" -ForegroundColor Cyan
Write-Host "Total files checked: $($requiredFiles.Count)"
Write-Host "Files missing: $fileFailures"
if ($fileFailures -gt 0) {
    Write-Host "`nMissing files need to be created:" -ForegroundColor Red
    foreach ($file in $requiredFiles) {
        if (-not (Test-Path $file.Path)) {
            Write-Host "   - $($file.Module): $($file.Path)" -ForegroundColor Red
        }
    }
} else {
    Write-Host "`nAll required files exist. Proceed with code review to verify implementation completeness." -ForegroundColor Green
}

# 7. Critical findings
Write-Host "`n===== Critical Implementation Notes =====" -ForegroundColor Yellow
Write-Host "1. Product Module must implement:"
Write-Host "   - Product & ProductImage models with correct fields"
Write-Host "   - Repository methods with transaction support"
Write-Host "   - Service layer with complete business logic and validations"
Write-Host "   - Controller layer with all 8 API endpoints"

Write-Host "`n2. Image Management must implement:"
Write-Host "   - File service with image validation (JPG/PNG, max 2MB)"
Write-Host "   - 4 image management API endpoints"

Write-Host "`n3. Category & Tag Module must implement:"
Write-Host "   - Complete models and repository layers"
Write-Host "   - Reference checking before deletion"
Write-Host "   - Proper error codes (4001 for categories, 4002 for tags)"

Write-Host "`nValidation script execution completed. Please fix any issues identified above." -ForegroundColor Green

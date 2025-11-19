# API 规格文档（`api.md`）

> **项目**：校园二手交易平台（前后端分离 · RESTful）
> **统一前缀**：`/api/v1`
> **说明**：本文定义**学生端 + 管理后台**所有对外 HTTP 接口及请求/响应契约，作为开发与联调的唯一依据。本平台为课程作业，一切接口仅面向教学场景。核心业务规则（**一物一件**、**商品状态机**、**卖家微信号展示规则**、**推荐与最近浏览** 等）均以当前 PRD/设计与数据库 schema 为准。

---

## 1. 概述（Overview）

* 系统采用 **前后端分离 + RESTful** 风格，所有资源路径以 `/api/v1` 为统一前缀，数据格式为 `application/json`（图片上传使用 `multipart/form-data`）。
* **商品状态机**：`ForSale`（在售）↔ `Delisted`（已下架），以及单向 `ForSale → Sold`（已售，**终态**、仅禁止状态字段反向变更；管理员可对已售商品**非状态字段**做纠错/数据清洗）。前台列表/搜索/推荐仅展示 `ForSale`。该约束在**业务**与**数据库触发器**双层落地。
* **一物一件**：每条 `products` 记录代表**一件实物**；无库存字段。 

---

## 2. 通用约定（General Conventions）

### 2.1 认证与权限

* **认证方式**：登录后颁发 token（建议 JWT），前端在后续请求头携带 `Authorization: Bearer <token>`。
* **权限判定**：

  * 需要登录的接口统一由鉴权中间件校验；
  * 管理端接口需 `isAdmin=true`（用户表含 `is_admin` 字段）；
  * 管理端登录与学生端登录共用登录接口，后端在登录响应中返回 `isAdmin`。

### 2.2 统一响应结构

```json
{
  "code": 0,
  "message": "ok",
  "data": {} 
}
```

* `code`：业务码；`0` 成功，其它为错误。
* `message`：人类可读提示。
* `data`：业务数据（对象/数组/分页对象）。 

> **时间/数字格式**
>
> * 所有时间字段均使用 ISO8601（UTC 或服务端时区，一致即可）：例 `2025-10-23T12:00:00Z`。
> * 金额以 `number`（精度见 `NUMERIC(10,2)`），入参与出参保留两位小数。 

### 2.3 分页与排序约定

* 通用 Query：`page`（默认 `1`），`pageSize`（默认 `20`，最大 `100`）。
* 排序字段按各接口定义；不提供前端可自由组合的任意字段排序。
* 分页响应统一结构：

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "items": [],
    "page": 1,
    "pageSize": 20,
    "total": 0
  }
}
```

### 2.4 统一错误码（节选）

| code | 语义                      |
| ---- | ----------------------- |
| 0    | 成功                      |
| 1001 | 参数校验失败 / 频率限制           |
| 1002 | 未登录或 token 无效           |
| 1003 | 权限不足（需要管理员或非本人的越权操作）    |
| 2001 | 账号已存在                   |
| 2002 | 账号或密码错误                 |
| 3001 | 商品不存在                   |
| 3002 | 非发布者，无权操作该商品            |
| 3003 | 商品状态不允许此操作（非法流转）        |
| 3004 | 已售商品为终态：禁止状态变更；普通卖家不可编辑 |
| 3005 | 撤销失败（超时或状态冲突）           |
| 4001 | 分类被引用，无法删除              |
| 4002 | 标签被引用，无法删除              |

> 依据详细设计的错误码框架与状态机规则。 

---

## 3. 数据字段与业务规则要点

### 3.1 状态与新旧程度

* **商品状态**：`ForSale` / `Delisted` / `Sold`（枚举型 `product_status`），并有触发器**强约束**合法流转；`Sold` 为终态，禁止任何状态回退（但允许在不改状态的前提下更正其它字段用于管理用途）。 
* **新旧程度**：使用 `product_conditions` 表（唯一事实来源）；API 采用 `conditionId`（可返回 `id`、`code`、`name` 供展示）。 

### 3.2 商品主图与图片

* 主图来源于 `product_images` 表中 `is_primary = true` 的记录（数据库唯一索引保证每商品最多一张主图）；接口可在卡片/详情中返回 `mainImageUrl` 字段用于快捷展示。

### 3.3 联系卖家（微信号）

* 卖家联系方式是**用户级**字段 `users.wechat_id`（非商品级字段）。
* “联系卖家”返回规则：

  * **未登录**：不返回微信号（`null`），提示需先登录；
  * **登录且查看者=卖家本人**：不展示“联系卖家”入口；
  * **登录且查看者≠卖家**：返回完整微信号（若无历史数据则 `null` 并返回降级文案）。

### 3.4 最近浏览与推荐

* 系统在用户查看商品详情时记录浏览（`user_recent_views`），**每用户保留最近 20 条**（DB 触发器自动裁剪）；首页“猜你喜欢”基于最近浏览商品的**标签 Top3~5** 在在售商品中做匹配，并与“最新发布”结果**去重**。

---

## 4. 模块与接口

> 下列每个接口均给出：**方法 + 路径、功能描述、认证要求、请求参数（Path/Query/Body）、响应数据（含示例）、错误情况**。字段命名采用 **lowerCamelCase**。

### 4.1 用户与认证模块

#### 4.1.1 注册

* **方法 + 路径**：`POST /api/v1/users/register`
* **功能**：创建新用户并自动登录（返回 token）。
* **认证**：不需要。
* **Request Body（JSON）**

  | 字段              | 类型     | 必填 | 说明         |
  | --------------- | ------ | -- | ---------- |
  | account         | string | 是  | 登录账号（字母数字） |
  | nickname        | string | 是  | 昵称         |
  | password        | string | 是  | 密码（≥8 位）   |
  | confirmPassword | string | 是  | 确认密码       |
  | wechatId        | string | 否  | 微信号（可选，注册时不强制，发布商品前需完善） |
* **Response (`data`)**

  ```json
  {
    "user": {
      "id": 1,
      "account": "alice01",
      "nickname": "Alice",
      "avatarUrl": null,
      "isAdmin": false,
      "createdAt": "2025-10-23T12:00:00Z"
    },
    "token": "eyJhbGciOi..."
  }
  ```
* **错误**：`2001` 账号已存在；`1001` 参数校验失败。 

#### 4.1.2 登录

* **方法 + 路径**：`POST /api/v1/users/login`
* **功能**：账号密码登录，返回 token；响应中附带 `isAdmin`。
* **认证**：不需要。
* **Request Body**

  | 字段         | 类型      | 必填 | 说明         |
  | ---------- | ------- | -- | ---------- |
  | account    | string  | 是  | 登录账号       |
  | password   | string  | 是  | 密码         |
  | rememberMe | boolean | 否  | 记住我（延长有效期） |
* **Response（同注册）**
* **错误**：`2002` 账号或密码错误。 

> **说明**：若实现服务端会话黑名单，可结合 `sessions` 表（可选），但本项目以 JWT 为主；`users.is_admin` 决定是否有管理端权限。 

#### 4.1.3 获取当前用户信息

* **方法 + 路径**：`GET /api/v1/users/profile`
* **功能**：获取当前登录用户的基础信息。
* **认证**：需要（登录）。
* **Response（示例）**

  ```json
  {
    "code": 0,
    "message": "ok",
    "data": {
      "id": 1,
      "account": "alice01",
      "nickname": "Alice",
      "avatarUrl": null,
      "wechatId": "alice_wx",
      "isAdmin": false,
      "createdAt": "2025-10-23T12:00:00Z",
      "updatedAt": "2025-10-23T12:00:00Z"
    }
  }
  ```
* **错误**：`1002` 未登录。 

#### 4.1.4 更新个人资料

* **方法 + 路径**：`PUT /api/v1/users/profile`
* **功能**：更新头像、昵称、（及可选）微信号；昵称支持**30 天频控**。
* **认证**：需要。
* **Request Body（任意字段可选）**

  | 字段        | 类型     | 必填 | 说明                          |
  | --------- | ------ | -- | --------------------------- |
  | nickname  | string | 否  | 昵称（可能受 30 天频控）              |
  | avatarUrl | string | 否  | 头像 URL                      |
  | wechatId  | string | 否  | 微信号（建议 4~64 字符，允许字母/数字/_/-） |
* **错误**：`1001` 参数/频控不通过。 

#### 4.1.5 修改密码

* **方法 + 路径**：`PUT /api/v1/users/password`
* **功能**：验证旧密码后设置新密码（当前登录态保持有效）。
* **认证**：需要。
* **Request Body**

  | 字段              | 类型     | 必填 | 说明        |
  | --------------- | ------ | -- | --------- |
  | oldPassword     | string | 是  | 旧密码       |
  | newPassword     | string | 是  | 新密码（≥8 位） |
  | confirmPassword | string | 是  | 确认新密码     |
* **错误**：`1001` 参数不合法或旧密码错误。 

---

### 4.2 商品模块（学生端）

#### 4.2.1 发布商品

* **方法 + 路径**：`POST /api/v1/products`
* **功能**：发布一件商品（默认 `ForSale`），支持多图上传。
* **认证**：需要（登录且用户 `wechatId` 已填写）。
* **Content-Type**：`multipart/form-data`
* **Form 字段**

  | 字段          | 类型     | 必填 | 说明                               |
  | ----------- | ------ | -- | -------------------------------- |
  | title       | string | 是  | 标题                               |
  | description | string | 是  | 描述                               |
  | price       | number | 是  | 价格（>0）                           |
  | categoryId  | number | 是  | 分类 ID                            |
  | tagIds      | string | 否  | 逗号分隔的标签 ID（如 `"1,2,3"`）          |
  | conditionId | number | 是  | 新旧程度 ID（来自 `product_conditions`） |
  | images      | file[] | 否  | 商品图片（多图）                         |
* **Response（简要）**

  ```json
  {
    "code": 0,
    "message": "ok",
    "data": { "id": 101, "status": "ForSale", "mainImageUrl": "..." }
  }
  ```
* **错误**：`1001`（字段/图片校验失败）；`3003`（状态异常——理论不应触发）；**未填写 wechatId 时**拒绝发布并引导用户完善资料。
* **说明**：发布时若包含图片，默认将第一张设为主图，并同步更新 `products.main_image_url`。

#### 4.2.2 编辑商品

* **方法 + 路径**：`PUT /api/v1/products/{id}`
* **功能**：更新标题/描述/价格/分类/标签/新旧程度/图片（覆盖式或由图片模块单独维护，见 4.3）；不改变状态。
* **认证**：需要（发布者本人；**普通卖家不可编辑 `Sold` 终态**）。
* **Request Body（JSON，任意字段可选）**

  | 字段          | 类型       | 说明                     |
  | ----------- | -------- | ---------------------- |
  | title       | string   |                        |
  | description | string   |                        |
  | price       | number   | >0                     |
  | categoryId  | number   |                        |
  | tagIds      | number[] |                        |
  | conditionId | number   |                        |
  | imageUrls   | string[] | 可选：覆盖式替换（若采用图片子模块可不使用） |
* **错误**：`3001` 商品不存在；`3002` 非发布者；`3004` 已售商品（普通卖家）不可编辑。

> **管理员例外**：后台接口允许在不变更 `status` 的前提下编辑已售商品的**非状态字段**（详见 4.7.4）。

#### 4.2.3 上/下架/标记已售

* **方法 + 路径**：`POST /api/v1/products/{id}/status`
* **功能**：在售↔下架、在售→已售；成功后支持**3 秒撤销**（仅上/下架）。
* **认证**：需要（发布者本人）。
* **Request Body**

  | 字段     | 类型     | 必填 | 说明                           |
  | ------ | ------ | -- | ---------------------------- |
  | action | string | 是  | `delist` / `relist` / `sold` |
* **Response（示例）**

  ```json
  {"code":0,"message":"ok","data":{"id":101,"status":"Delisted"}}
  ```
* **错误**：`3004` Sold 终态禁止任何状态变更；`3003` 非法流转。

> **撤销窗口**：上/下架成功后，服务端缓存记录最近一次状态（TTL≈3s）；`sold` 不可撤销。 

#### 4.2.4 撤销上/下架

* **方法 + 路径**：`POST /api/v1/products/{id}/status/undo`
* **功能**：在 3 秒窗口内撤销最近一次上/下架（若当前状态未被他操作改变）。
* **认证**：需要（发布者本人）。
* **错误**：`3005` 撤销超时/冲突；`3004` 若目标为 `Sold` 不允许撤销。 

#### 4.2.5 获取单个商品详情

* **方法 + 路径**：`GET /api/v1/products/{id}`
* **功能**：返回详情（含卖家基础信息、图片、标签、`mainImageUrl`）；若带登录态，将记录最近浏览。
* **认证**：可匿名（若带 token 将记录浏览）。
* **Response（示意）**

  ```json
  {
    "code": 0,
    "message": "ok",
    "data": {
      "id": 101,
      "title": "iPad 2021",
      "description": "...",
      "price": 1599.00,
      "conditionId": 2,
      "conditionName": "九成新",
      "categoryId": 10,
      "status": "ForSale",
      "mainImageUrl": "https://img/xxx.jpg",
      "images": [
        {"id": 1, "url": "https://img/xxx.jpg", "sortOrder": 0, "isPrimary": true},
        {"id": 2, "url": "https://img/yyy.jpg", "sortOrder": 1, "isPrimary": false}
      ],
      "tagIds": [3, 5],
      "seller": {
        "id": 7, "nickname": "Tom", "avatarUrl": "https://..."
      },
      "viewerIsSeller": false,
      "sellerWechat": null  // 见“联系卖家”接口
    }
  }
  ```

> 记录浏览写入 `user_recent_views`（并保留最近 20 条）；主图由 `product_images.is_primary=true` 推导。

#### 4.2.6 我发布的商品列表

* **方法 + 路径**：`GET /api/v1/products/my`
* **功能**：当前登录用户的所有商品（含所有状态），支持关键词过滤。
* **认证**：需要。
* **Query**

  | 名称       | 类型     | 必填 | 说明     |
  | -------- | ------ | -- | ------ |
  | keyword  | string | 否  | 标题模糊搜索 |
  | page     | number | 否  | 默认 1   |
  | pageSize | number | 否  | 默认 20  |
* **Response**：分页结构，`items` 含 `id/title/price/status/createdAt/mainImageUrl`。 

#### 4.2.7 搜索商品

* **方法 + 路径**：`GET /api/v1/products/search`
* **功能**：仅查询 `ForSale` 商品；支持关键词、价格区间、时间范围、新旧程度、分类/标签、排序。
* **认证**：无需。
* **Query**

  | 名称                 | 类型     | 必填 | 说明                                     |
  | ------------------ | ------ | -- | -------------------------------------- |
  | q                  | string | 否  | 关键词（标题/描述/分类/标签）                       |
  | categoryId         | number | 否  | 分类过滤                                   |
  | tagIds             | string | 否  | 逗号分隔标签 ID                              |
  | conditionIds       | string | 否  | 逗号分隔新旧程度 ID                            |
  | minPrice           | number | 否  | 最小价                                    |
  | maxPrice           | number | 否  | 最大价                                    |
  | publishedTimeRange | string | 否  | `all`（默认）/`last_7_days`/`last_30_days` |
  | sort               | string | 否  | `latest`（默认）/`priceAsc`/`priceDesc`    |
  | page/pageSize      | number | 否  | 分页                                     |
* **Response**：分页商品卡片列表（含 `mainImageUrl`）。

> 关键词可匹配标题/描述/分类/标签；时间范围基于 `products.created_at` 过滤。

#### 4.2.8 分类下商品列表

* **方法 + 路径**：`GET /api/v1/products/category/{categoryId}`
* **功能**：浏览某分类下在售商品；支持价格与排序。
* **认证**：无需。
* **Query**：`minPrice` / `maxPrice` / `sort` / `page` / `pageSize`。 

---

### 4.3 商品图片模块

> `product_images`：`id/productId/url/sortOrder/isPrimary`；每商品最多一张 `isPrimary=true`（唯一索引保障）。 

#### 4.3.1 上传商品图片（追加）

* **方法 + 路径**：`POST /api/v1/products/{id}/images`
* **功能**：为指定商品追加图片。
* **认证**：需要（发布者本人；`Sold` 终态普通卖家不可变更图片）。
* **Content-Type**：`multipart/form-data`
* **Form 字段**：`images`（file[]，支持多图）。
* **Response**：返回本次新增的图片元信息数组。

#### 4.3.2 设置主图

* **方法 + 路径**：`PUT /api/v1/products/{id}/images/{imageId}/primary`
* **功能**：将指定图片设为主图（DB 层保证唯一）。
* **认证**：需要（发布者本人；`Sold` 终态普通卖家不可编辑）。
* **Response**：`{ id, isPrimary: true }`
* **说明**：实现时需先将该商品下所有图片置为非主图（`is_primary=false`），再设置新主图，并同步更新 `products.main_image_url` 字段。

#### 4.3.3 调整排序

* **方法 + 路径**：`PATCH /api/v1/products/{id}/images/{imageId}`
* **功能**：更新 `sortOrder`。
* **认证**：需要。
* **Request Body**：`{ "sortOrder": 2 }`

#### 4.3.4 删除图片

* **方法 + 路径**：`DELETE /api/v1/products/{id}/images/{imageId}`
* **功能**：删除某张图片。
* **认证**：需要（发布者本人；`Sold` 终态普通卖家不可编辑）。

---

### 4.4 分类与标签模块

#### 4.4.1 获取分类列表

* **方法 + 路径**：`GET /api/v1/categories`
* **功能**：列出所有分类（用于发布/筛选）。
* **认证**：无需。
* **Response**：`[{ id, name, description }]`。 

#### 4.4.2 获取标签列表

* **方法 + 路径**：`GET /api/v1/tags`
* **功能**：列出所有标签。
* **认证**：无需。
* **Response**：`[{ id, name }]`。 

> **管理员维护接口**：见 4.7 分类/标签管理；删除前需检查是否被商品引用，否则返回 `4001/4002`。 

---

### 4.5 新旧程度枚举模块

#### 4.5.1 获取新旧程度列表

* **方法 + 路径**：`GET /api/v1/product-conditions`
* **功能**：读取 `product_conditions`；前端发布/编辑页的单选来源。
* **认证**：无需。
* **Response**：`[{ id, code, name, sortOrder }]`。 

---

### 4.6 浏览 / 搜索 / 推荐模块

#### 4.6.1 首页数据（推荐 + 最新）

* **方法 + 路径**：`GET /api/v1/home`
* **功能**：返回（登录态）“猜你喜欢”与（所有用户）“最新发布在售”分页，同时**去重**。
* **认证**：可选（登录则有推荐）。
* **Query**：`page/pageSize`（作用于“最新发布”）。
* **Response（示例）**

  ```json
  {
    "code": 0,
    "data": {
      "recommendations": [ { "id": 11, "title": "...", "mainImageUrl": "..." } ],
      "latest": {
        "items": [ { "id": 22, "title": "...", "mainImageUrl": "..." } ],
        "page": 1, "pageSize": 20, "total": 135
      }
    }
  }
  ```

> 推荐基于最近 20 条浏览的**标签**高频匹配，仅含 `ForSale` 且排除本人发布；不足 5 条时用“最新发布在售”补齐，首页展示跨模块去重。

#### 4.6.2（可选）显式记录浏览

* **方法 + 路径**：`POST /api/v1/views`
* **功能**：前端在特殊场景主动上报浏览（常规由 `GET /products/{id}` 自动记录可不调用）。
* **认证**：需要。
* **Request Body**：`{ "productId": 101, "viewedAt": "2025-10-23T12:00:00Z" }`

> 常规实现中，详情接口已自动记录并由 DB 触发器裁剪到最近 20 条。 

#### 4.6.3 获取最近浏览（当前用户）

* **方法 + 路径**：`GET /api/v1/users/recent-views`
* **功能**：返回当前用户最近浏览的商品列表（最多 20）。
* **认证**：需要。
* **Response**：`items` 为商品卡片（`id/title/price/mainImageUrl/createdAt`）。 

---

### 4.7 联系卖家模块

#### 4.7.1 点击“联系卖家”

* **方法 + 路径**：`GET /api/v1/products/{id}/contact`
* **功能**：根据请求者身份返回卖家微信号或提示需要登录。
* **认证**：可匿名（未登录返回 `null` 并提示登录）。
* **Response（示例）**

  ```json
  {
    "code": 0,
    "message": "ok",
    "data": {
      "canContact": true,
      "sellerWechat": "seller_wx_123", 
      "tips": null
    }
  }
  ```

  * 未登录：`{ "canContact": false, "sellerWechat": null, "tips": "请先登录后联系卖家" }`
  * 登录且为卖家本人：`{ "canContact": false, "sellerWechat": null }`（前端不显示入口）
  * 登录且非卖家：若卖家 `wechatId` 为空（历史数据）：`sellerWechat=null`，`tips="卖家联系方式暂不可用，请稍后再试"`。

> 卖家联系方式来自 **`users.wechat_id`**（用户级），非商品字段。

---

### 4.8 管理员后台模块

> 管理端接口仅在 `Authorization` 鉴权通过且 `isAdmin=true` 时放行；登录共用用户登录接口。 

#### 4.8.1 仪表盘统计

* **方法 + 路径**：`GET /api/v1/admin/dashboard`
* **功能**：返回用户总数、商品总数、在售/已售数等统计。
* **认证**：需要（管理员）。
* **Response（示例）**

  ```json
  {
    "code": 0,
    "data": {
      "userCount": 1200,
      "productCount": 3400,
      "forSaleCount": 2100,
      "soldCount": 980
    }
  }
  ```



#### 4.8.2 用户列表（后台）

* **方法 + 路径**：`GET /api/v1/admin/users`
* **功能**：分页查询用户，支持关键词。
* **认证**：需要（管理员）。
* **Query**：`keyword/page/pageSize`
* **Response**：分页 `items` 含 `id/account/nickname/avatarUrl/wechatId/isAdmin/createdAt`。 

#### 4.8.3 商品列表（后台）

* **方法 + 路径**：`GET /api/v1/admin/products`
* **功能**：后台分页查看商品，支持状态/发布者/关键词过滤。
* **认证**：需要（管理员）。
* **Query**：`status`（`ForSale/Delisted/Sold`）/`sellerId`/`keyword`/`page/pageSize`。 

#### 4.8.4 管理员纠错编辑已售商品（非状态字段）

* **方法 + 路径**：`PUT /api/v1/admin/products/{id}`
* **功能**：允许在 `status` **不改变** 的前提下编辑 `Sold` 商品的**非状态字段**（标题/描述/分类/标签/图片/新旧程度等）。
* **认证**：需要（管理员）。
* **Request Body**：同 4.2.2 中可编辑字段，但不得携带 `status`。
* **错误**：若尝试改变状态或将 `Sold` 改为其它状态，一律拒绝（`3004`）。

> 该能力用于**纠错/数据清洗**；与数据库触发器的“Sold 终态禁止状态字段变更”一致。

#### 4.8.5 分类管理（增改删）

* **新增分类**：`POST /api/v1/admin/categories`（需要管理员）

  * Body：`{ "name": "数码", "description": "..." }`
* **修改分类**：`PUT /api/v1/admin/categories/{id}`

  * Body：`{ "name": "数码电子", "description": "..." }`
* **删除分类**：`DELETE /api/v1/admin/categories/{id}`

  * 规则：若有商品引用（`products.category_id`），拒绝删除（`4001`）。 

#### 4.8.6 标签管理（增改删）

* **新增标签**：`POST /api/v1/admin/tags`
* **修改标签**：`PUT /api/v1/admin/tags/{id}`
* **删除标签**：`DELETE /api/v1/admin/tags/{id}`

  * 规则：若被 `product_tags` 引用，拒绝删除（`4002`）。 

---

## 5. 字段模型（DTO 摘要）

> 为便于前端类型定义与后端实现，列出常用数据结构（出参）。

### 5.1 商品卡片 `ProductCard`

```ts
{
  id: number;
  title: string;
  price: number;
  mainImageUrl: string | null; // 冗余字段，来自 products.main_image_url
  createdAt: string;           // ISO8601
}
```

### 5.2 商品详情 `ProductDetail`

```ts
{
  id: number;
  title: string;
  description: string;
  price: number;
  conditionId: number;
  conditionName?: string;
  categoryId: number;
  status: "ForSale" | "Delisted" | "Sold";
  mainImageUrl: string | null; // 冗余字段，来自 products.main_image_url
  images: { id: number; url: string; sortOrder: number; isPrimary: boolean }[];
  tagIds: number[];
  seller: { id: number; nickname: string; avatarUrl?: string | null };
  viewerIsSeller: boolean;
  sellerWechat?: string | null; // 仅在满足可见性规则时返回
  createdAt: string;
  updatedAt: string;
}
```

### 5.3 分类与标签

```ts
Category { id:number; name:string; description?:string; }
Tag { id:number; name:string; }
ProductCondition { id:number; code:string; name:string; sortOrder:number; }
```



---

## 6. 业务错误与边界情形（补充）

* **Sold 终态相关**

  * 普通卖家：`Sold` 商品不可编辑、不可上下架、不可撤销；返回 `3004`。
  * 管理员：可编辑 `Sold` 的**非状态字段**，但任何尝试修改 `status` 必须被拒绝。
  * 数据库已配置触发器防止 `Sold -> 其他` 的状态变更。 
* **发布前校验卖家联系方式**：`wechatId` 为空时禁止发布，返回明确提示。 
* **最近浏览**：由详情接口自动写入，DB 触发器裁剪到**最近 20 条**。 
* **主图唯一**：设置主图时应清理其它图片的 `isPrimary`，DB 层唯一索引保证至多一张主图；同时更新 `products.main_image_url`。 

---

## 7. 开发/联调须知

* **路由与分层**：接口由后端 Controller 暴露；业务在 Service 层实现；Repository 访问 PostgreSQL；公共模块提供统一响应/错误码/鉴权。前后端工程与模块分布详见项目结构文档。 
* **接口列表总览**（与详细设计一致）

  * 用户：注册、登录、`GET/PUT /users/profile`、`PUT /users/password`；
  * 商品：`POST /products`、`PUT /products/{id}`、`POST /products/{id}/status`、`POST /products/{id}/status/undo`、`GET /products/{id}`、`GET /products/my`、`GET /products/search`、`GET /products/category/{categoryId}`；
  * 枚举：`GET /product-conditions`；
  * 图片：`POST/PUT/PATCH/DELETE /products/{id}/images...`；
  * 分类/标签：`GET /categories`、`GET /tags`；（管理员：`POST/PUT/DELETE`）
  * 首页：`GET /home`；
  * 最近浏览：`GET /users/recent-views`（可选 `POST /views`）。 

---

## 8. 非对外接口（仅说明，不对前端公开）

* **健康检查**：`GET /internal/healthz`（用于进程存活探测，前端不调用）。
* **说明**：该类接口仅服务开发/运维自检，不属于前端联调范畴。
  （若项目实际未实现，可忽略；如实现请保持私有路由前缀/鉴权策略与网关限制。）

---

## 9. 附：请求/响应示例汇总（选摘）

> **示例：搜索商品**
> `GET /api/v1/products/search?q=ipad&minPrice=500&maxPrice=3000&publishedTimeRange=last_30_days&sort=latest&page=1&pageSize=20`

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "items": [
      { "id": 101, "title": "iPad 2021", "price": 1599.00, "mainImageUrl": "https://...", "createdAt": "2025-10-01T10:00:00Z" }
    ],
    "page": 1,
    "pageSize": 20,
    "total": 1
  }
}
```

（仅返回 `ForSale`，时间范围基于 `created_at`；排序为最新发布优先。） 

> **示例：管理员删除分类（被引用）**
> `DELETE /api/v1/admin/categories/10`

```json
{ "code": 4001, "message": "分类被引用，无法删除", "data": null }
```

（删除前执行引用检查。） 

---

### 版本与一致性说明

* 本 `api.md` 与 **PRD / 概要设计 / 详细设计 / 项目结构文档 / schema.sql** 对齐；当多处描述略有出入时，以本文当前约定为联调契约，并同步更新相关文档。
* 关键一致性锚点：**状态机不可逆（Sold 终态）**、**新旧程度来自 `product_conditions`**、**主图由 `product_images.is_primary` 决定**、**联系卖家按登录与身份差异返回**、**最近浏览仅保留 20 条**。



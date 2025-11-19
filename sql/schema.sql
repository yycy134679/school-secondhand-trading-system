-- ============================================================
-- 校园二手交易平台 (PostgreSQL) - 3NF 关系数据库结构
-- 说明：
--   1) 严格遵循“一物一件”与商品状态机（在售/已售/已下架）；
--   2) 设计符合第三范式（3NF）：无重复组、无部分依赖、无传递依赖；
--   3) 包含必要的 CHECK/UNIQUE/FK/INDEX 与触发器以约束核心业务规则；
--   4) 可直接在 Navicat 或 psql 中执行；如无超级权限，pg_trgm 扩展会自动忽略。
-- ============================================================

BEGIN;

-- 可选扩展：用于标题/描述的模糊搜索（LIKE/ILIKE、相似度匹配）
DO $$
BEGIN
    CREATE EXTENSION IF NOT EXISTS pg_trgm;
EXCEPTION
    WHEN OTHERS THEN
        RAISE NOTICE 'pg_trgm extension not created (permissions). Proceeding without it.';
END$$;

-- 先删除旧对象（谨慎：这会清空已有数据）
DROP TABLE IF EXISTS user_recent_views CASCADE;
DROP TABLE IF EXISTS product_tags CASCADE;
DROP TABLE IF EXISTS product_images CASCADE;
DROP TABLE IF EXISTS products CASCADE;
DROP TABLE IF EXISTS tags CASCADE;
DROP TABLE IF EXISTS categories CASCADE;
DROP TABLE IF EXISTS product_conditions CASCADE;
DROP TABLE IF EXISTS sessions CASCADE;
DROP TABLE IF EXISTS users CASCADE;

DROP FUNCTION IF EXISTS trg_set_updated_at() CASCADE;
DROP FUNCTION IF EXISTS trg_products_status_guard() CASCADE;
DROP FUNCTION IF EXISTS trg_prune_user_recent_views() CASCADE;

DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'product_status') THEN
        DROP TYPE product_status;
    END IF;
END$$;

-- 商品状态枚举：ForSale / Sold(终态) / Delisted
CREATE TYPE product_status AS ENUM ('ForSale', 'Sold', 'Delisted');

-- ------------------------------------------------------------
-- 通用：自动更新时间戳触发器函数
-- ------------------------------------------------------------
CREATE OR REPLACE FUNCTION trg_set_updated_at()
RETURNS trigger AS $$
BEGIN
    NEW.updated_at := NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- ------------------------------------------------------------
-- 1. 用户表
-- ------------------------------------------------------------
CREATE TABLE users (
    id                       BIGSERIAL PRIMARY KEY,
    account                  VARCHAR(32)  NOT NULL UNIQUE,
    nickname                 VARCHAR(32)  NOT NULL,
    password_hash            VARCHAR(128) NOT NULL,
    avatar_url               TEXT,
    wechat_id                VARCHAR(64), -- 允许 NULL，注册时不填，发布商品前强制要求
    is_admin                 BOOLEAN      NOT NULL DEFAULT FALSE,
    last_nickname_changed_at TIMESTAMPTZ,
    created_at               TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at               TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    CONSTRAINT ck_users_account_format CHECK (account ~ '^[A-Za-z0-9]+$'),
    CONSTRAINT ck_users_password_len    CHECK (char_length(password_hash) >= 8) -- 实际建议≥60（bcrypt），此处为最小保护
);

CREATE INDEX idx_users_created_at ON users (created_at DESC);

CREATE TRIGGER users_set_updated_at
BEFORE UPDATE ON users
FOR EACH ROW EXECUTE FUNCTION trg_set_updated_at();

COMMENT ON TABLE users IS '系统用户（学生/管理员）。账号唯一；昵称可重复；密码以哈希存储。';
COMMENT ON COLUMN users.account IS '登录账号（仅字母与数字），全局唯一。';
COMMENT ON COLUMN users.nickname IS '用于展示的昵称（非唯一）。';
COMMENT ON COLUMN users.password_hash IS '密码哈希（如 bcrypt/argon2），禁止明文存储。';
COMMENT ON COLUMN users.is_admin IS '是否管理员。';
COMMENT ON COLUMN users.last_nickname_changed_at IS '上次昵称修改时间，用于 30 天修改频控。';
COMMENT ON COLUMN users.wechat_id IS '用户微信号（用于联系卖家，用户级字段）。建议长度 4~64。注册时可为空，发布商品时要求填写（不允许空值）。';

-- ------------------------------------------------------------
-- 2. 新旧程度（枚举表）
-- ------------------------------------------------------------
CREATE TABLE product_conditions (
    id          SMALLSERIAL PRIMARY KEY,
    code        VARCHAR(32)  NOT NULL UNIQUE,  -- 机器可读编码，如 BRAND_NEW/NINE_TENTHS 等
    name        VARCHAR(50)  NOT NULL,         -- 展示名：如“全新”“九成新”
    sort_order  SMALLINT     NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX uq_product_conditions_name ON product_conditions (name);

CREATE TRIGGER product_conditions_set_updated_at
BEFORE UPDATE ON product_conditions
FOR EACH ROW EXECUTE FUNCTION trg_set_updated_at();

COMMENT ON TABLE product_conditions IS '商品新旧程度枚举表（发布页单选使用），此表为新旧程度的唯一事实来源。';
COMMENT ON COLUMN product_conditions.code IS '新旧程度编码（唯一）。';
COMMENT ON COLUMN product_conditions.name IS '新旧程度中文名（如 全新/九成新/七成新 等）。';

-- ------------------------------------------------------------
-- 3. 商品分类
-- ------------------------------------------------------------
CREATE TABLE categories (
    id          BIGSERIAL PRIMARY KEY,
    name        VARCHAR(50)  NOT NULL UNIQUE,
    description VARCHAR(255),
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TRIGGER categories_set_updated_at
BEFORE UPDATE ON categories
FOR EACH ROW EXECUTE FUNCTION trg_set_updated_at();

COMMENT ON TABLE categories IS '商品分类（由管理员维护）。删除被引用的分类将因外键而失败。';

-- ------------------------------------------------------------
-- 4. 标签
-- ------------------------------------------------------------
CREATE TABLE tags (
    id          BIGSERIAL PRIMARY KEY,
    name        VARCHAR(50)  NOT NULL UNIQUE,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TRIGGER tags_set_updated_at
BEFORE UPDATE ON tags
FOR EACH ROW EXECUTE FUNCTION trg_set_updated_at();

COMMENT ON TABLE tags IS '商品标签库（由管理员维护，商品可多选标签）。';

-- ------------------------------------------------------------
-- 5. 商品主表（严格遵循“一物一件”）
-- ------------------------------------------------------------
CREATE TABLE products (
    id           BIGSERIAL     PRIMARY KEY,
    seller_id    BIGINT        NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    title        VARCHAR(100)  NOT NULL,
    description  TEXT          NOT NULL,
    price        NUMERIC(10,2) NOT NULL,
    condition_id SMALLINT      NOT NULL REFERENCES product_conditions(id) ON DELETE RESTRICT,
    category_id  BIGINT        NOT NULL REFERENCES categories(id) ON DELETE RESTRICT,
    status       product_status NOT NULL DEFAULT 'ForSale',
    main_image_url VARCHAR(255), -- 冗余字段，用于列表查询优化，与 product_images.is_primary 同步
    created_at   TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    CONSTRAINT ck_products_price_positive CHECK (price > 0)
);

-- 索引：按常用检索建立
CREATE INDEX idx_products_seller_created ON products (seller_id, created_at DESC);
CREATE INDEX idx_products_status_created ON products (status, created_at DESC);
CREATE INDEX idx_products_category_status_created ON products (category_id, status, created_at DESC);
CREATE INDEX idx_products_status_price ON products (status, price);

-- 可选：模糊搜索优化（若 pg_trgm 可用）
DO $$
BEGIN
    PERFORM 1 FROM pg_extension WHERE extname = 'pg_trgm';
    IF FOUND THEN
        EXECUTE 'CREATE INDEX IF NOT EXISTS idx_products_title_trgm ON products USING gin (title gin_trgm_ops)';
        EXECUTE 'CREATE INDEX IF NOT EXISTS idx_products_desc_trgm  ON products USING gin (description gin_trgm_ops)';
    END IF;
END$$;

CREATE TRIGGER products_set_updated_at
BEFORE UPDATE ON products
FOR EACH ROW EXECUTE FUNCTION trg_set_updated_at();

COMMENT ON TABLE products IS '商品主表：每条记录代表一件实物（无库存字段）。';
COMMENT ON COLUMN products.seller_id IS '发布者用户 ID（1:N 关系：用户→商品）。';
COMMENT ON COLUMN products.status IS '状态机：ForSale(在售) / Delisted(已下架) / Sold(已售-终态)。';
COMMENT ON COLUMN products.main_image_url IS '主图 URL 冗余字段，用于列表展示优化。发布/编辑/设置主图时需同步更新此字段。';

-- ------------------------------------------------------------
-- 5.1 商品状态机约束触发器：禁止非法流转，Sold 为状态终态（仅禁止状态字段由 Sold 变更为其他）
-- ------------------------------------------------------------
CREATE OR REPLACE FUNCTION trg_products_status_guard()
RETURNS trigger AS $$
BEGIN
    -- 状态校验：仅在尝试变更 status 字段时进行严格校验。
    -- 注意：Sold 为状态终态，禁止将 Sold 变更为其他状态；
    -- 但当 status 未变化（NEW.status = OLD.status）时，允许更新其他字段（例如 title/description/category/tags）。
    IF NEW.status <> OLD.status THEN
        -- 禁止 Sold -> 非 Sold 的变更
        IF OLD.status = 'Sold' THEN
            RAISE EXCEPTION 'Product % is Sold and its status cannot be changed (terminal state).', OLD.id
                USING ERRCODE = '45000';
        END IF;

        -- 允许的状态流转：ForSale -> Delisted, Delisted -> ForSale, ForSale -> Sold
        IF NOT (
            (OLD.status = 'ForSale'  AND NEW.status = 'Delisted') OR
            (OLD.status = 'Delisted' AND NEW.status = 'ForSale')  OR
            (OLD.status = 'ForSale'  AND NEW.status = 'Sold')
        ) THEN
            RAISE EXCEPTION 'Invalid product status transition: % -> %', OLD.status, NEW.status
                USING ERRCODE = '45000';
        END IF;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER products_status_guard
BEFORE UPDATE ON products
FOR EACH ROW EXECUTE FUNCTION trg_products_status_guard();

COMMENT ON TRIGGER products_status_guard ON products IS
'约束商品状态机：ForSale↔Delisted 互转；ForSale→Sold 终态；禁止从 Sold 变更为其他状态（状态字段不可逆）。状态为 Sold 时，若不修改 status 字段，则允许更新其他非状态字段（如标题/描述/分类/标签），以便管理员纠错或数据清洗。';

-- ------------------------------------------------------------
-- 6. 商品图片（多张，且最多一张主图）
-- ------------------------------------------------------------
CREATE TABLE product_images (
    id          BIGSERIAL PRIMARY KEY,
    product_id  BIGINT     NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    url         TEXT       NOT NULL,
    sort_order  INT        NOT NULL DEFAULT 0,
    is_primary  BOOLEAN    NOT NULL DEFAULT FALSE
);

CREATE INDEX idx_product_images_product_sort ON product_images (product_id, sort_order);

-- 保证每个商品最多只有一张主图
CREATE UNIQUE INDEX uq_product_images_primary_one
ON product_images (product_id)
WHERE is_primary = TRUE;

COMMENT ON TABLE product_images IS '商品图片表：一对多。主图通过 is_primary 标记并由唯一索引保证每商品最多一张主图；其余按 sort_order 排序。';

COMMENT ON COLUMN products.condition_id IS '引用 product_conditions 表（唯一事实来源）；前端应使用 conditionId 作为入参，响应可返回 id 与名称/编码供展示。';

-- ------------------------------------------------------------
-- 7. 商品-标签 多对多
-- ------------------------------------------------------------
CREATE TABLE product_tags (
    product_id BIGINT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    tag_id     BIGINT NOT NULL REFERENCES tags(id)      ON DELETE RESTRICT,
    PRIMARY KEY (product_id, tag_id)
);

CREATE INDEX idx_product_tags_tag ON product_tags (tag_id);

COMMENT ON TABLE product_tags IS '商品与标签的多对多关联表（复合主键防重复）。';

-- ------------------------------------------------------------
-- 8. 最近浏览记录（用于推荐；每用户仅保留最新 20 条）
-- ------------------------------------------------------------
CREATE TABLE user_recent_views (
    id         BIGSERIAL PRIMARY KEY,
    user_id    BIGINT       NOT NULL REFERENCES users(id)    ON DELETE CASCADE,
    product_id BIGINT       NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    viewed_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_views_user_time ON user_recent_views (user_id, viewed_at DESC, id DESC);

COMMENT ON TABLE user_recent_views IS '用户最近浏览商品记录（用于“猜你喜欢”）；按 user_id + viewed_at 倒序查询。';

-- 触发器：插入后裁剪为最近 20 条
CREATE OR REPLACE FUNCTION trg_prune_user_recent_views()
RETURNS trigger AS $$
BEGIN
    DELETE FROM user_recent_views
    WHERE user_id = NEW.user_id
      AND id IN (
        SELECT id
        FROM user_recent_views
        WHERE user_id = NEW.user_id
        ORDER BY viewed_at DESC, id DESC
        OFFSET 20
      );
    RETURN NULL; -- AFTER INSERT, return value ignored
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER user_recent_views_prune_after_insert
AFTER INSERT ON user_recent_views
FOR EACH ROW EXECUTE FUNCTION trg_prune_user_recent_views();

-- ------------------------------------------------------------
-- 9.（可选）会话表：若采用服务端会话/黑名单，可启用；使用纯 JWT 可忽略
-- ------------------------------------------------------------
CREATE TABLE sessions (
    id         BIGSERIAL PRIMARY KEY,
    user_id    BIGINT       NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token      VARCHAR(128) NOT NULL UNIQUE,
    expired_at TIMESTAMPTZ  NOT NULL,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_sessions_user ON sessions (user_id);
CREATE INDEX idx_sessions_expired_at ON sessions (expired_at);

COMMENT ON TABLE sessions IS '可选：服务端会话/令牌黑名单存储；若使用纯 JWT + Redis，可不创建本表。';

COMMIT;

-- ======================== 关系说明（摘要） ========================
-- users (1) —— (N) products：一个用户可发布多件商品；
-- categories (1) —— (N) products：每件商品隶属一个分类；
-- product_conditions (1) —— (N) products：每件商品选择一个新旧程度；
-- products (1) —— (N) product_images：商品可有多张图片；每商品最多一张 is_primary=TRUE；
-- products (N) —— (N) tags：通过 product_tags 关联；
-- users (1) —— (N) user_recent_views：用于记录浏览历史与推荐；
-- 关键业务约束在 DB 层显式落实：
--   · “一物一件”：不设库存字段；每条 products 即一件实物；
--   · 状态机：触发器 products_status_guard 严格限制流转；Sold 为不可逆终态；
--   · 状态机：触发器 products_status_guard 严格限制状态字段流转（禁止 Sold -> 其他）；状态为 Sold 时若 status 未变则允许更新其他字段（例如标题/描述/分类/标签），以支持管理员纠错或数据清洗。
--   · 最近浏览：触发器 user_recent_views_prune_after_insert 保持每用户仅保留 20 条记录；
--   · 联系方式：用户 `wechat_id` 为用户级联系方式；系统在商品详情页点击“联系卖家”时从 `users` 表读取该字段并返回给满足权限的已登录用户；若历史记录缺失请返回降级文案或 null。 

-- ============================================================
-- 迁移建议（注意：如果已有 users 数据，请采用两步迁移以避免破坏现有记录）：
-- 1) 新增列（允许 NULL），并在前端/后端强制要求新/编辑用户填写微信号；
--    ALTER TABLE users ADD COLUMN wechat_id VARCHAR(64);
--    COMMENT ON COLUMN users.wechat_id IS '用户微信号，迁移初期允许 NULL；发布商品需填写';
-- 2) 后续通过脚本或手动操作补齐旧数据，并通过 UI 强制用户完善；
--    UPDATE users SET wechat_id = '<placeholder>' WHERE wechat_id IS NULL;
-- 3) 验证后将列改为 NOT NULL（同样可加 DEFAULT 以便批量回填）：
--    ALTER TABLE users ALTER COLUMN wechat_id SET NOT NULL;
--    -- 或者逐步实施：
--    -- 1. 先设置默认值空串（或占位），保证现有记录能通过 NOT NULL：
--    --    ALTER TABLE users ALTER COLUMN wechat_id SET DEFAULT '';
--    --    UPDATE users SET wechat_id = '' WHERE wechat_id IS NULL;
--    -- 2. 再将列改为 NOT NULL：
--    --    ALTER TABLE users ALTER COLUMN wechat_id SET NOT NULL;
--    -- 可在上一步先设置默认值再修改为 NOT NULL，或者在更新后删除默认值。
-- ============================================================

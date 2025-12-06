/*
 Navicat Premium Dump SQL

 Source Server         : postgres
 Source Server Type    : PostgreSQL
 Source Server Version : 160010 (160010)
 Source Host           : 43.136.104.67:5432
 Source Catalog        : school-secondhand-trading
 Source Schema         : public

 Target Server Type    : PostgreSQL
 Target Server Version : 160010 (160010)
 File Encoding         : 65001

 Date: 06/12/2025 12:12:15
*/


-- ----------------------------
-- Type structure for gtrgm
-- ----------------------------
DROP TYPE IF EXISTS "public"."gtrgm";
CREATE TYPE "public"."gtrgm" (
  INPUT = "public"."gtrgm_in",
  OUTPUT = "public"."gtrgm_out",
  INTERNALLENGTH = VARIABLE,
  CATEGORY = U,
  DELIMITER = ','
);
ALTER TYPE "public"."gtrgm" OWNER TO "postgres";

-- ----------------------------
-- Type structure for product_status
-- ----------------------------
DROP TYPE IF EXISTS "public"."product_status";
CREATE TYPE "public"."product_status" AS ENUM (
  'ForSale',
  'Sold',
  'Delisted'
);
ALTER TYPE "public"."product_status" OWNER TO "postgres";

-- ----------------------------
-- Sequence structure for categories_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."categories_id_seq";
CREATE SEQUENCE "public"."categories_id_seq"
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;
ALTER SEQUENCE "public"."categories_id_seq" OWNER TO "postgres";

-- ----------------------------
-- Sequence structure for product_conditions_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."product_conditions_id_seq";
CREATE SEQUENCE "public"."product_conditions_id_seq"
INCREMENT 1
MINVALUE  1
MAXVALUE 32767
START 1
CACHE 1;
ALTER SEQUENCE "public"."product_conditions_id_seq" OWNER TO "postgres";

-- ----------------------------
-- Sequence structure for product_images_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."product_images_id_seq";
CREATE SEQUENCE "public"."product_images_id_seq"
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;
ALTER SEQUENCE "public"."product_images_id_seq" OWNER TO "postgres";

-- ----------------------------
-- Sequence structure for products_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."products_id_seq";
CREATE SEQUENCE "public"."products_id_seq"
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;
ALTER SEQUENCE "public"."products_id_seq" OWNER TO "postgres";

-- ----------------------------
-- Sequence structure for sessions_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."sessions_id_seq";
CREATE SEQUENCE "public"."sessions_id_seq"
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;
ALTER SEQUENCE "public"."sessions_id_seq" OWNER TO "postgres";

-- ----------------------------
-- Sequence structure for simple_users_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."simple_users_id_seq";
CREATE SEQUENCE "public"."simple_users_id_seq"
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;
ALTER SEQUENCE "public"."simple_users_id_seq" OWNER TO "postgres";

-- ----------------------------
-- Sequence structure for tags_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."tags_id_seq";
CREATE SEQUENCE "public"."tags_id_seq"
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;
ALTER SEQUENCE "public"."tags_id_seq" OWNER TO "postgres";

-- ----------------------------
-- Sequence structure for test_products_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."test_products_id_seq";
CREATE SEQUENCE "public"."test_products_id_seq"
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;
ALTER SEQUENCE "public"."test_products_id_seq" OWNER TO "postgres";

-- ----------------------------
-- Sequence structure for test_users_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."test_users_id_seq";
CREATE SEQUENCE "public"."test_users_id_seq"
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;
ALTER SEQUENCE "public"."test_users_id_seq" OWNER TO "postgres";

-- ----------------------------
-- Sequence structure for user_recent_views_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."user_recent_views_id_seq";
CREATE SEQUENCE "public"."user_recent_views_id_seq"
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;
ALTER SEQUENCE "public"."user_recent_views_id_seq" OWNER TO "postgres";

-- ----------------------------
-- Sequence structure for users_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."users_id_seq";
CREATE SEQUENCE "public"."users_id_seq"
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;
ALTER SEQUENCE "public"."users_id_seq" OWNER TO "postgres";

-- ----------------------------
-- Table structure for categories
-- ----------------------------
DROP TABLE IF EXISTS "public"."categories";
CREATE TABLE "public"."categories" (
  "id" int8 NOT NULL DEFAULT nextval('categories_id_seq'::regclass),
  "name" varchar(50) COLLATE "pg_catalog"."default" NOT NULL,
  "description" varchar(255) COLLATE "pg_catalog"."default",
  "created_at" timestamptz(6) NOT NULL DEFAULT now(),
  "updated_at" timestamptz(6) NOT NULL DEFAULT now()
)
;
ALTER TABLE "public"."categories" OWNER TO "postgres";
COMMENT ON TABLE "public"."categories" IS '商品分类（由管理员维护）。删除被引用的分类将因外键而失败。';

-- ----------------------------
-- Records of categories
-- ----------------------------
BEGIN;
INSERT INTO "public"."categories" ("id", "name", "description", "created_at", "updated_at") VALUES (1, '数码电子', '各类手机、电脑及配件', '2025-12-06 11:28:34.396243+08', '2025-12-06 11:28:34.396243+08');
INSERT INTO "public"."categories" ("id", "name", "description", "created_at", "updated_at") VALUES (2, '图书教材', '本科教材、考研考公资料', '2025-12-06 11:28:34.396243+08', '2025-12-06 11:28:34.396243+08');
INSERT INTO "public"."categories" ("id", "name", "description", "created_at", "updated_at") VALUES (3, '生活日用', '宿舍神器、收纳、小家电', '2025-12-06 11:28:34.396243+08', '2025-12-06 11:28:34.396243+08');
INSERT INTO "public"."categories" ("id", "name", "description", "created_at", "updated_at") VALUES (4, '运动器材', '各类球拍、健身器材', '2025-12-06 11:28:34.396243+08', '2025-12-06 11:28:34.396243+08');
COMMIT;

-- ----------------------------
-- Table structure for product_conditions
-- ----------------------------
DROP TABLE IF EXISTS "public"."product_conditions";
CREATE TABLE "public"."product_conditions" (
  "id" int2 NOT NULL DEFAULT nextval('product_conditions_id_seq'::regclass),
  "code" varchar(32) COLLATE "pg_catalog"."default" NOT NULL,
  "name" varchar(50) COLLATE "pg_catalog"."default" NOT NULL,
  "sort_order" int2 NOT NULL DEFAULT 0,
  "created_at" timestamptz(6) NOT NULL DEFAULT now(),
  "updated_at" timestamptz(6) NOT NULL DEFAULT now()
)
;
ALTER TABLE "public"."product_conditions" OWNER TO "postgres";
COMMENT ON COLUMN "public"."product_conditions"."code" IS '新旧程度编码（唯一）。';
COMMENT ON COLUMN "public"."product_conditions"."name" IS '新旧程度中文名（如 全新/九成新/七成新 等）。';
COMMENT ON TABLE "public"."product_conditions" IS '商品新旧程度枚举表（发布页单选使用），此表为新旧程度的唯一事实来源。';

-- ----------------------------
-- Records of product_conditions
-- ----------------------------
BEGIN;
INSERT INTO "public"."product_conditions" ("id", "code", "name", "sort_order", "created_at", "updated_at") VALUES (1, 'BRAND_NEW', '全新', 10, '2025-12-05 22:11:21.909582+08', '2025-12-05 22:11:21.909582+08');
INSERT INTO "public"."product_conditions" ("id", "code", "name", "sort_order", "created_at", "updated_at") VALUES (2, 'NINE_TENTHS', '九成新', 20, '2025-12-05 22:11:21.909582+08', '2025-12-05 22:11:21.909582+08');
INSERT INTO "public"."product_conditions" ("id", "code", "name", "sort_order", "created_at", "updated_at") VALUES (3, 'EIGHT_TENTHS', '八成新', 30, '2025-12-05 22:11:21.909582+08', '2025-12-05 22:11:21.909582+08');
INSERT INTO "public"."product_conditions" ("id", "code", "name", "sort_order", "created_at", "updated_at") VALUES (4, 'SEVEN_TENTHS', '七成新', 40, '2025-12-05 22:11:21.909582+08', '2025-12-05 22:11:21.909582+08');
COMMIT;

-- ----------------------------
-- Table structure for product_images
-- ----------------------------
DROP TABLE IF EXISTS "public"."product_images";
CREATE TABLE "public"."product_images" (
  "id" int8 NOT NULL DEFAULT nextval('product_images_id_seq'::regclass),
  "product_id" int8 NOT NULL,
  "url" text COLLATE "pg_catalog"."default" NOT NULL,
  "sort_order" int4 NOT NULL DEFAULT 0,
  "is_primary" bool NOT NULL DEFAULT false
)
;
ALTER TABLE "public"."product_images" OWNER TO "postgres";
COMMENT ON TABLE "public"."product_images" IS '商品图片表：一对多。主图通过 is_primary 标记并由唯一索引保证每商品最多一张主图；其余按 sort_order 排序。';

-- ----------------------------
-- Records of product_images
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for product_tags
-- ----------------------------
DROP TABLE IF EXISTS "public"."product_tags";
CREATE TABLE "public"."product_tags" (
  "product_id" int8 NOT NULL,
  "tag_id" int8 NOT NULL
)
;
ALTER TABLE "public"."product_tags" OWNER TO "postgres";
COMMENT ON TABLE "public"."product_tags" IS '商品与标签的多对多关联表（复合主键防重复）。';

-- ----------------------------
-- Records of product_tags
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for products
-- ----------------------------
DROP TABLE IF EXISTS "public"."products";
CREATE TABLE "public"."products" (
  "id" int8 NOT NULL DEFAULT nextval('products_id_seq'::regclass),
  "seller_id" int8 NOT NULL,
  "title" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "description" text COLLATE "pg_catalog"."default" NOT NULL,
  "price" numeric(10,2) NOT NULL,
  "condition_id" int2 NOT NULL,
  "category_id" int8 NOT NULL,
  "status" "public"."product_status" NOT NULL DEFAULT 'ForSale'::product_status,
  "main_image_url" varchar(255) COLLATE "pg_catalog"."default",
  "created_at" timestamptz(6) NOT NULL DEFAULT now(),
  "updated_at" timestamptz(6) NOT NULL DEFAULT now()
)
;
ALTER TABLE "public"."products" OWNER TO "postgres";
COMMENT ON COLUMN "public"."products"."seller_id" IS '发布者用户 ID（1:N 关系：用户→商品）。';
COMMENT ON COLUMN "public"."products"."condition_id" IS '引用 product_conditions 表（唯一事实来源）；前端应使用 conditionId 作为入参，响应可返回 id 与名称/编码供展示。';
COMMENT ON COLUMN "public"."products"."status" IS '状态机：ForSale(在售) / Delisted(已下架) / Sold(已售-终态)。';
COMMENT ON COLUMN "public"."products"."main_image_url" IS '主图 URL 冗余字段，用于列表展示优化。发布/编辑/设置主图时需同步更新此字段。';
COMMENT ON TABLE "public"."products" IS '商品主表：每条记录代表一件实物（无库存字段）。';

-- ----------------------------
-- Records of products
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for sessions
-- ----------------------------
DROP TABLE IF EXISTS "public"."sessions";
CREATE TABLE "public"."sessions" (
  "id" int8 NOT NULL DEFAULT nextval('sessions_id_seq'::regclass),
  "user_id" int8 NOT NULL,
  "token" varchar(128) COLLATE "pg_catalog"."default" NOT NULL,
  "expired_at" timestamptz(6) NOT NULL,
  "created_at" timestamptz(6) NOT NULL DEFAULT now()
)
;
ALTER TABLE "public"."sessions" OWNER TO "postgres";
COMMENT ON TABLE "public"."sessions" IS '可选：服务端会话/令牌黑名单存储；若使用纯 JWT + Redis，可不创建本表。';

-- ----------------------------
-- Records of sessions
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for simple_users
-- ----------------------------
DROP TABLE IF EXISTS "public"."simple_users";
CREATE TABLE "public"."simple_users" (
  "id" int8 NOT NULL DEFAULT nextval('simple_users_id_seq'::regclass),
  "account" text COLLATE "pg_catalog"."default",
  "created_at" timestamptz(6)
)
;
ALTER TABLE "public"."simple_users" OWNER TO "postgres";

-- ----------------------------
-- Records of simple_users
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for tags
-- ----------------------------
DROP TABLE IF EXISTS "public"."tags";
CREATE TABLE "public"."tags" (
  "id" int8 NOT NULL DEFAULT nextval('tags_id_seq'::regclass),
  "name" varchar(50) COLLATE "pg_catalog"."default" NOT NULL,
  "created_at" timestamptz(6) NOT NULL DEFAULT now(),
  "updated_at" timestamptz(6) NOT NULL DEFAULT now(),
  "category_id" int8 NOT NULL
)
;
ALTER TABLE "public"."tags" OWNER TO "postgres";
COMMENT ON TABLE "public"."tags" IS '商品标签库（由管理员维护，商品可多选标签）。';

-- ----------------------------
-- Records of tags
-- ----------------------------
BEGIN;
INSERT INTO "public"."tags" ("id", "name", "created_at", "updated_at", "category_id") VALUES (1, '手机', '2025-12-06 11:33:15.301066+08', '2025-12-06 11:33:15.301066+08', 1);
INSERT INTO "public"."tags" ("id", "name", "created_at", "updated_at", "category_id") VALUES (2, '平板', '2025-12-06 11:33:15.301066+08', '2025-12-06 11:33:15.301066+08', 1);
INSERT INTO "public"."tags" ("id", "name", "created_at", "updated_at", "category_id") VALUES (3, '耳机', '2025-12-06 11:33:15.301066+08', '2025-12-06 11:33:15.301066+08', 1);
INSERT INTO "public"."tags" ("id", "name", "created_at", "updated_at", "category_id") VALUES (4, '教材', '2025-12-06 11:33:15.665251+08', '2025-12-06 11:33:15.665251+08', 2);
INSERT INTO "public"."tags" ("id", "name", "created_at", "updated_at", "category_id") VALUES (5, '考研', '2025-12-06 11:33:15.665251+08', '2025-12-06 11:33:15.665251+08', 2);
INSERT INTO "public"."tags" ("id", "name", "created_at", "updated_at", "category_id") VALUES (6, '小说', '2025-12-06 11:33:15.665251+08', '2025-12-06 11:33:15.665251+08', 2);
INSERT INTO "public"."tags" ("id", "name", "created_at", "updated_at", "category_id") VALUES (7, '台灯', '2025-12-06 11:33:15.805128+08', '2025-12-06 11:33:15.805128+08', 3);
INSERT INTO "public"."tags" ("id", "name", "created_at", "updated_at", "category_id") VALUES (8, '收纳', '2025-12-06 11:33:15.805128+08', '2025-12-06 11:33:15.805128+08', 3);
INSERT INTO "public"."tags" ("id", "name", "created_at", "updated_at", "category_id") VALUES (9, '雨伞', '2025-12-06 11:33:15.805128+08', '2025-12-06 11:33:15.805128+08', 3);
INSERT INTO "public"."tags" ("id", "name", "created_at", "updated_at", "category_id") VALUES (10, '球拍', '2025-12-06 11:33:15.867164+08', '2025-12-06 11:33:15.867164+08', 4);
INSERT INTO "public"."tags" ("id", "name", "created_at", "updated_at", "category_id") VALUES (11, '瑜伽垫', '2025-12-06 11:33:15.867164+08', '2025-12-06 11:33:15.867164+08', 4);
INSERT INTO "public"."tags" ("id", "name", "created_at", "updated_at", "category_id") VALUES (12, '滑板', '2025-12-06 11:33:15.867164+08', '2025-12-06 11:33:15.867164+08', 4);
COMMIT;

-- ----------------------------
-- Table structure for test_products
-- ----------------------------
DROP TABLE IF EXISTS "public"."test_products";
CREATE TABLE "public"."test_products" (
  "id" int8 NOT NULL DEFAULT nextval('test_products_id_seq'::regclass),
  "seller_id" int8 NOT NULL,
  "title" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "created_at" timestamp(6) DEFAULT CURRENT_TIMESTAMP
)
;
ALTER TABLE "public"."test_products" OWNER TO "postgres";

-- ----------------------------
-- Records of test_products
-- ----------------------------
BEGIN;
INSERT INTO "public"."test_products" ("id", "seller_id", "title", "created_at") VALUES (1, 1, 'TestProduct_1764610630', '2025-12-02 01:37:10.78393');
INSERT INTO "public"."test_products" ("id", "seller_id", "title", "created_at") VALUES (2, 3, 'TestProductNoWechat_1764610630', '2025-12-02 01:37:10.807305');
COMMIT;

-- ----------------------------
-- Table structure for test_users
-- ----------------------------
DROP TABLE IF EXISTS "public"."test_users";
CREATE TABLE "public"."test_users" (
  "id" int8 NOT NULL DEFAULT nextval('test_users_id_seq'::regclass),
  "account" varchar(50) COLLATE "pg_catalog"."default" NOT NULL,
  "wechat_id" varchar(64) COLLATE "pg_catalog"."default"
)
;
ALTER TABLE "public"."test_users" OWNER TO "postgres";

-- ----------------------------
-- Records of test_users
-- ----------------------------
BEGIN;
INSERT INTO "public"."test_users" ("id", "account", "wechat_id") VALUES (1, 'test_seller_1764610630', 'seller_wechat_123');
INSERT INTO "public"."test_users" ("id", "account", "wechat_id") VALUES (2, 'test_buyer_1764610630', 'buyer_wechat_456');
INSERT INTO "public"."test_users" ("id", "account", "wechat_id") VALUES (3, 'test_seller_no_wechat_1764610630', '');
COMMIT;

-- ----------------------------
-- Table structure for user_recent_views
-- ----------------------------
DROP TABLE IF EXISTS "public"."user_recent_views";
CREATE TABLE "public"."user_recent_views" (
  "id" int8 NOT NULL DEFAULT nextval('user_recent_views_id_seq'::regclass),
  "user_id" int8 NOT NULL,
  "product_id" int8 NOT NULL,
  "viewed_at" timestamptz(6) NOT NULL DEFAULT now()
)
;
ALTER TABLE "public"."user_recent_views" OWNER TO "postgres";
COMMENT ON TABLE "public"."user_recent_views" IS '用户最近浏览商品记录（用于“猜你喜欢”）；按 user_id + viewed_at 倒序查询。';

-- ----------------------------
-- Records of user_recent_views
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for users
-- ----------------------------
DROP TABLE IF EXISTS "public"."users";
CREATE TABLE "public"."users" (
  "id" int8 NOT NULL DEFAULT nextval('users_id_seq'::regclass),
  "account" varchar(32) COLLATE "pg_catalog"."default" NOT NULL,
  "nickname" varchar(32) COLLATE "pg_catalog"."default" NOT NULL,
  "password_hash" varchar(128) COLLATE "pg_catalog"."default" NOT NULL,
  "avatar_url" text COLLATE "pg_catalog"."default",
  "wechat_id" varchar(64) COLLATE "pg_catalog"."default",
  "is_admin" bool NOT NULL DEFAULT false,
  "last_nickname_changed_at" timestamptz(6),
  "created_at" timestamptz(6) NOT NULL DEFAULT now(),
  "updated_at" timestamptz(6) NOT NULL DEFAULT now()
)
;
ALTER TABLE "public"."users" OWNER TO "postgres";
COMMENT ON COLUMN "public"."users"."account" IS '登录账号（仅字母与数字），全局唯一。';
COMMENT ON COLUMN "public"."users"."nickname" IS '用于展示的昵称（非唯一）。';
COMMENT ON COLUMN "public"."users"."password_hash" IS '密码哈希（如 bcrypt/argon2），禁止明文存储。';
COMMENT ON COLUMN "public"."users"."wechat_id" IS '用户微信号（用于联系卖家，用户级字段）。建议长度 4~64。注册时可为空，发布商品时要求填写（不允许空值）。';
COMMENT ON COLUMN "public"."users"."is_admin" IS '是否管理员。';
COMMENT ON COLUMN "public"."users"."last_nickname_changed_at" IS '上次昵称修改时间，用于 30 天修改频控。';
COMMENT ON TABLE "public"."users" IS '系统用户（学生/管理员）。账号唯一；昵称可重复；密码以哈希存储。';

-- ----------------------------
-- Records of users
-- ----------------------------
BEGIN;
INSERT INTO "public"."users" ("id", "account", "nickname", "password_hash", "avatar_url", "wechat_id", "is_admin", "last_nickname_changed_at", "created_at", "updated_at") VALUES (2, 'zrl123456', 'zrl123456', '$2a$10$cpw5U9iPEr0DfdsrLqM8zeQf7hxjU9OssRtj5ReeJovOoV81Ek9iq', '', 'zrl123456', 'f', NULL, '2025-12-04 11:06:43.495689+08', '2025-12-04 11:06:43.495689+08');
INSERT INTO "public"."users" ("id", "account", "nickname", "password_hash", "avatar_url", "wechat_id", "is_admin", "last_nickname_changed_at", "created_at", "updated_at") VALUES (1, 'czs123456', 'czs123456', '$2a$10$mSPRkSvgNKBXaEEb/n0WSuNKeZdKD7cMM4WEQJUlD7TjzP1dMNAA2', 'http://localhost:8080/uploads/2025/12/1764933578599.jpg', 'czs123456', 'f', NULL, '2025-12-04 11:03:13.50694+08', '2025-12-05 19:19:38.727693+08');
INSERT INTO "public"."users" ("id", "account", "nickname", "password_hash", "avatar_url", "wechat_id", "is_admin", "last_nickname_changed_at", "created_at", "updated_at") VALUES (7, 'yycy134679', '云烟成雨', '$2a$10$322Qy97u9unXq2uU7cTNBO1LVSnEL3l2Ulk5XeB0E093BWot946jO', 'http://localhost:8080/uploads/2025/12/1764935884835.png', 'yycy134679', 'f', NULL, '2025-12-05 19:39:04.513154+08', '2025-12-05 19:58:08.990622+08');
INSERT INTO "public"."users" ("id", "account", "nickname", "password_hash", "avatar_url", "wechat_id", "is_admin", "last_nickname_changed_at", "created_at", "updated_at") VALUES (8, '123456', '12345678', '$2a$10$qllt0fUhRDJ1KjxiaXpD5Oxt0.iGu.W3Xc6on8LwHFEfVaXpEy73C', '', '123456', 'f', '2025-12-05 19:59:19.171266+08', '2025-12-05 19:59:04.318609+08', '2025-12-05 19:59:19.326148+08');
INSERT INTO "public"."users" ("id", "account", "nickname", "password_hash", "avatar_url", "wechat_id", "is_admin", "last_nickname_changed_at", "created_at", "updated_at") VALUES (9, 'testuser1', 'Test User1', '$2a$10$y5NAueU416BumanQaKySMu/u4KLGluSxHCj2LSsfY0ESTftdNxowy', '', 'wx_test', 'f', NULL, '2025-12-06 10:25:35.434736+08', '2025-12-06 10:25:35.434736+08');
COMMIT;

-- ----------------------------
-- Function structure for gin_extract_query_trgm
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."gin_extract_query_trgm"(text, internal, int2, internal, internal, internal, internal);
CREATE FUNCTION "public"."gin_extract_query_trgm"(text, internal, int2, internal, internal, internal, internal)
  RETURNS "pg_catalog"."internal" AS '$libdir/pg_trgm', 'gin_extract_query_trgm'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."gin_extract_query_trgm"(text, internal, int2, internal, internal, internal, internal) OWNER TO "postgres";

-- ----------------------------
-- Function structure for gin_extract_value_trgm
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."gin_extract_value_trgm"(text, internal);
CREATE FUNCTION "public"."gin_extract_value_trgm"(text, internal)
  RETURNS "pg_catalog"."internal" AS '$libdir/pg_trgm', 'gin_extract_value_trgm'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."gin_extract_value_trgm"(text, internal) OWNER TO "postgres";

-- ----------------------------
-- Function structure for gin_trgm_consistent
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."gin_trgm_consistent"(internal, int2, text, int4, internal, internal, internal, internal);
CREATE FUNCTION "public"."gin_trgm_consistent"(internal, int2, text, int4, internal, internal, internal, internal)
  RETURNS "pg_catalog"."bool" AS '$libdir/pg_trgm', 'gin_trgm_consistent'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."gin_trgm_consistent"(internal, int2, text, int4, internal, internal, internal, internal) OWNER TO "postgres";

-- ----------------------------
-- Function structure for gin_trgm_triconsistent
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."gin_trgm_triconsistent"(internal, int2, text, int4, internal, internal, internal);
CREATE FUNCTION "public"."gin_trgm_triconsistent"(internal, int2, text, int4, internal, internal, internal)
  RETURNS "pg_catalog"."char" AS '$libdir/pg_trgm', 'gin_trgm_triconsistent'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."gin_trgm_triconsistent"(internal, int2, text, int4, internal, internal, internal) OWNER TO "postgres";

-- ----------------------------
-- Function structure for gtrgm_compress
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."gtrgm_compress"(internal);
CREATE FUNCTION "public"."gtrgm_compress"(internal)
  RETURNS "pg_catalog"."internal" AS '$libdir/pg_trgm', 'gtrgm_compress'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."gtrgm_compress"(internal) OWNER TO "postgres";

-- ----------------------------
-- Function structure for gtrgm_consistent
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."gtrgm_consistent"(internal, text, int2, oid, internal);
CREATE FUNCTION "public"."gtrgm_consistent"(internal, text, int2, oid, internal)
  RETURNS "pg_catalog"."bool" AS '$libdir/pg_trgm', 'gtrgm_consistent'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."gtrgm_consistent"(internal, text, int2, oid, internal) OWNER TO "postgres";

-- ----------------------------
-- Function structure for gtrgm_decompress
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."gtrgm_decompress"(internal);
CREATE FUNCTION "public"."gtrgm_decompress"(internal)
  RETURNS "pg_catalog"."internal" AS '$libdir/pg_trgm', 'gtrgm_decompress'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."gtrgm_decompress"(internal) OWNER TO "postgres";

-- ----------------------------
-- Function structure for gtrgm_distance
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."gtrgm_distance"(internal, text, int2, oid, internal);
CREATE FUNCTION "public"."gtrgm_distance"(internal, text, int2, oid, internal)
  RETURNS "pg_catalog"."float8" AS '$libdir/pg_trgm', 'gtrgm_distance'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."gtrgm_distance"(internal, text, int2, oid, internal) OWNER TO "postgres";

-- ----------------------------
-- Function structure for gtrgm_in
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."gtrgm_in"(cstring);
CREATE FUNCTION "public"."gtrgm_in"(cstring)
  RETURNS "public"."gtrgm" AS '$libdir/pg_trgm', 'gtrgm_in'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."gtrgm_in"(cstring) OWNER TO "postgres";

-- ----------------------------
-- Function structure for gtrgm_options
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."gtrgm_options"(internal);
CREATE FUNCTION "public"."gtrgm_options"(internal)
  RETURNS "pg_catalog"."void" AS '$libdir/pg_trgm', 'gtrgm_options'
  LANGUAGE c IMMUTABLE
  COST 1;
ALTER FUNCTION "public"."gtrgm_options"(internal) OWNER TO "postgres";

-- ----------------------------
-- Function structure for gtrgm_out
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."gtrgm_out"("public"."gtrgm");
CREATE FUNCTION "public"."gtrgm_out"("public"."gtrgm")
  RETURNS "pg_catalog"."cstring" AS '$libdir/pg_trgm', 'gtrgm_out'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."gtrgm_out"("public"."gtrgm") OWNER TO "postgres";

-- ----------------------------
-- Function structure for gtrgm_penalty
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."gtrgm_penalty"(internal, internal, internal);
CREATE FUNCTION "public"."gtrgm_penalty"(internal, internal, internal)
  RETURNS "pg_catalog"."internal" AS '$libdir/pg_trgm', 'gtrgm_penalty'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."gtrgm_penalty"(internal, internal, internal) OWNER TO "postgres";

-- ----------------------------
-- Function structure for gtrgm_picksplit
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."gtrgm_picksplit"(internal, internal);
CREATE FUNCTION "public"."gtrgm_picksplit"(internal, internal)
  RETURNS "pg_catalog"."internal" AS '$libdir/pg_trgm', 'gtrgm_picksplit'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."gtrgm_picksplit"(internal, internal) OWNER TO "postgres";

-- ----------------------------
-- Function structure for gtrgm_same
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."gtrgm_same"("public"."gtrgm", "public"."gtrgm", internal);
CREATE FUNCTION "public"."gtrgm_same"("public"."gtrgm", "public"."gtrgm", internal)
  RETURNS "pg_catalog"."internal" AS '$libdir/pg_trgm', 'gtrgm_same'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."gtrgm_same"("public"."gtrgm", "public"."gtrgm", internal) OWNER TO "postgres";

-- ----------------------------
-- Function structure for gtrgm_union
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."gtrgm_union"(internal, internal);
CREATE FUNCTION "public"."gtrgm_union"(internal, internal)
  RETURNS "public"."gtrgm" AS '$libdir/pg_trgm', 'gtrgm_union'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."gtrgm_union"(internal, internal) OWNER TO "postgres";

-- ----------------------------
-- Function structure for set_limit
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."set_limit"(float4);
CREATE FUNCTION "public"."set_limit"(float4)
  RETURNS "pg_catalog"."float4" AS '$libdir/pg_trgm', 'set_limit'
  LANGUAGE c VOLATILE STRICT
  COST 1;
ALTER FUNCTION "public"."set_limit"(float4) OWNER TO "postgres";

-- ----------------------------
-- Function structure for show_limit
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."show_limit"();
CREATE FUNCTION "public"."show_limit"()
  RETURNS "pg_catalog"."float4" AS '$libdir/pg_trgm', 'show_limit'
  LANGUAGE c STABLE STRICT
  COST 1;
ALTER FUNCTION "public"."show_limit"() OWNER TO "postgres";

-- ----------------------------
-- Function structure for show_trgm
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."show_trgm"(text);
CREATE FUNCTION "public"."show_trgm"(text)
  RETURNS "pg_catalog"."_text" AS '$libdir/pg_trgm', 'show_trgm'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."show_trgm"(text) OWNER TO "postgres";

-- ----------------------------
-- Function structure for similarity
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."similarity"(text, text);
CREATE FUNCTION "public"."similarity"(text, text)
  RETURNS "pg_catalog"."float4" AS '$libdir/pg_trgm', 'similarity'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."similarity"(text, text) OWNER TO "postgres";

-- ----------------------------
-- Function structure for similarity_dist
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."similarity_dist"(text, text);
CREATE FUNCTION "public"."similarity_dist"(text, text)
  RETURNS "pg_catalog"."float4" AS '$libdir/pg_trgm', 'similarity_dist'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."similarity_dist"(text, text) OWNER TO "postgres";

-- ----------------------------
-- Function structure for similarity_op
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."similarity_op"(text, text);
CREATE FUNCTION "public"."similarity_op"(text, text)
  RETURNS "pg_catalog"."bool" AS '$libdir/pg_trgm', 'similarity_op'
  LANGUAGE c STABLE STRICT
  COST 1;
ALTER FUNCTION "public"."similarity_op"(text, text) OWNER TO "postgres";

-- ----------------------------
-- Function structure for strict_word_similarity
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."strict_word_similarity"(text, text);
CREATE FUNCTION "public"."strict_word_similarity"(text, text)
  RETURNS "pg_catalog"."float4" AS '$libdir/pg_trgm', 'strict_word_similarity'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."strict_word_similarity"(text, text) OWNER TO "postgres";

-- ----------------------------
-- Function structure for strict_word_similarity_commutator_op
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."strict_word_similarity_commutator_op"(text, text);
CREATE FUNCTION "public"."strict_word_similarity_commutator_op"(text, text)
  RETURNS "pg_catalog"."bool" AS '$libdir/pg_trgm', 'strict_word_similarity_commutator_op'
  LANGUAGE c STABLE STRICT
  COST 1;
ALTER FUNCTION "public"."strict_word_similarity_commutator_op"(text, text) OWNER TO "postgres";

-- ----------------------------
-- Function structure for strict_word_similarity_dist_commutator_op
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."strict_word_similarity_dist_commutator_op"(text, text);
CREATE FUNCTION "public"."strict_word_similarity_dist_commutator_op"(text, text)
  RETURNS "pg_catalog"."float4" AS '$libdir/pg_trgm', 'strict_word_similarity_dist_commutator_op'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."strict_word_similarity_dist_commutator_op"(text, text) OWNER TO "postgres";

-- ----------------------------
-- Function structure for strict_word_similarity_dist_op
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."strict_word_similarity_dist_op"(text, text);
CREATE FUNCTION "public"."strict_word_similarity_dist_op"(text, text)
  RETURNS "pg_catalog"."float4" AS '$libdir/pg_trgm', 'strict_word_similarity_dist_op'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."strict_word_similarity_dist_op"(text, text) OWNER TO "postgres";

-- ----------------------------
-- Function structure for strict_word_similarity_op
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."strict_word_similarity_op"(text, text);
CREATE FUNCTION "public"."strict_word_similarity_op"(text, text)
  RETURNS "pg_catalog"."bool" AS '$libdir/pg_trgm', 'strict_word_similarity_op'
  LANGUAGE c STABLE STRICT
  COST 1;
ALTER FUNCTION "public"."strict_word_similarity_op"(text, text) OWNER TO "postgres";

-- ----------------------------
-- Function structure for trg_products_status_guard
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."trg_products_status_guard"();
CREATE FUNCTION "public"."trg_products_status_guard"()
  RETURNS "pg_catalog"."trigger" AS $BODY$
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
$BODY$
  LANGUAGE plpgsql VOLATILE
  COST 100;
ALTER FUNCTION "public"."trg_products_status_guard"() OWNER TO "postgres";

-- ----------------------------
-- Function structure for trg_prune_user_recent_views
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."trg_prune_user_recent_views"();
CREATE FUNCTION "public"."trg_prune_user_recent_views"()
  RETURNS "pg_catalog"."trigger" AS $BODY$
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
$BODY$
  LANGUAGE plpgsql VOLATILE
  COST 100;
ALTER FUNCTION "public"."trg_prune_user_recent_views"() OWNER TO "postgres";

-- ----------------------------
-- Function structure for trg_set_updated_at
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."trg_set_updated_at"();
CREATE FUNCTION "public"."trg_set_updated_at"()
  RETURNS "pg_catalog"."trigger" AS $BODY$
BEGIN
    NEW.updated_at := NOW();
    RETURN NEW;
END;
$BODY$
  LANGUAGE plpgsql VOLATILE
  COST 100;
ALTER FUNCTION "public"."trg_set_updated_at"() OWNER TO "postgres";

-- ----------------------------
-- Function structure for word_similarity
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."word_similarity"(text, text);
CREATE FUNCTION "public"."word_similarity"(text, text)
  RETURNS "pg_catalog"."float4" AS '$libdir/pg_trgm', 'word_similarity'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."word_similarity"(text, text) OWNER TO "postgres";

-- ----------------------------
-- Function structure for word_similarity_commutator_op
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."word_similarity_commutator_op"(text, text);
CREATE FUNCTION "public"."word_similarity_commutator_op"(text, text)
  RETURNS "pg_catalog"."bool" AS '$libdir/pg_trgm', 'word_similarity_commutator_op'
  LANGUAGE c STABLE STRICT
  COST 1;
ALTER FUNCTION "public"."word_similarity_commutator_op"(text, text) OWNER TO "postgres";

-- ----------------------------
-- Function structure for word_similarity_dist_commutator_op
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."word_similarity_dist_commutator_op"(text, text);
CREATE FUNCTION "public"."word_similarity_dist_commutator_op"(text, text)
  RETURNS "pg_catalog"."float4" AS '$libdir/pg_trgm', 'word_similarity_dist_commutator_op'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."word_similarity_dist_commutator_op"(text, text) OWNER TO "postgres";

-- ----------------------------
-- Function structure for word_similarity_dist_op
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."word_similarity_dist_op"(text, text);
CREATE FUNCTION "public"."word_similarity_dist_op"(text, text)
  RETURNS "pg_catalog"."float4" AS '$libdir/pg_trgm', 'word_similarity_dist_op'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;
ALTER FUNCTION "public"."word_similarity_dist_op"(text, text) OWNER TO "postgres";

-- ----------------------------
-- Function structure for word_similarity_op
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."word_similarity_op"(text, text);
CREATE FUNCTION "public"."word_similarity_op"(text, text)
  RETURNS "pg_catalog"."bool" AS '$libdir/pg_trgm', 'word_similarity_op'
  LANGUAGE c STABLE STRICT
  COST 1;
ALTER FUNCTION "public"."word_similarity_op"(text, text) OWNER TO "postgres";

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."categories_id_seq"
OWNED BY "public"."categories"."id";
SELECT setval('"public"."categories_id_seq"', 4, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."product_conditions_id_seq"
OWNED BY "public"."product_conditions"."id";
SELECT setval('"public"."product_conditions_id_seq"', 4, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."product_images_id_seq"
OWNED BY "public"."product_images"."id";
SELECT setval('"public"."product_images_id_seq"', 25, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."products_id_seq"
OWNED BY "public"."products"."id";
SELECT setval('"public"."products_id_seq"', 28, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."sessions_id_seq"
OWNED BY "public"."sessions"."id";
SELECT setval('"public"."sessions_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."simple_users_id_seq"
OWNED BY "public"."simple_users"."id";
SELECT setval('"public"."simple_users_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."tags_id_seq"
OWNED BY "public"."tags"."id";
SELECT setval('"public"."tags_id_seq"', 12, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."test_products_id_seq"
OWNED BY "public"."test_products"."id";
SELECT setval('"public"."test_products_id_seq"', 2, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."test_users_id_seq"
OWNED BY "public"."test_users"."id";
SELECT setval('"public"."test_users_id_seq"', 3, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."user_recent_views_id_seq"
OWNED BY "public"."user_recent_views"."id";
SELECT setval('"public"."user_recent_views_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."users_id_seq"
OWNED BY "public"."users"."id";
SELECT setval('"public"."users_id_seq"', 9, true);

-- ----------------------------
-- Triggers structure for table categories
-- ----------------------------
CREATE TRIGGER "categories_set_updated_at" BEFORE UPDATE ON "public"."categories"
FOR EACH ROW
EXECUTE PROCEDURE "public"."trg_set_updated_at"();

-- ----------------------------
-- Uniques structure for table categories
-- ----------------------------
ALTER TABLE "public"."categories" ADD CONSTRAINT "categories_name_key" UNIQUE ("name");

-- ----------------------------
-- Primary Key structure for table categories
-- ----------------------------
ALTER TABLE "public"."categories" ADD CONSTRAINT "categories_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table product_conditions
-- ----------------------------
CREATE UNIQUE INDEX "uq_product_conditions_name" ON "public"."product_conditions" USING btree (
  "name" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

-- ----------------------------
-- Triggers structure for table product_conditions
-- ----------------------------
CREATE TRIGGER "product_conditions_set_updated_at" BEFORE UPDATE ON "public"."product_conditions"
FOR EACH ROW
EXECUTE PROCEDURE "public"."trg_set_updated_at"();

-- ----------------------------
-- Uniques structure for table product_conditions
-- ----------------------------
ALTER TABLE "public"."product_conditions" ADD CONSTRAINT "product_conditions_code_key" UNIQUE ("code");

-- ----------------------------
-- Primary Key structure for table product_conditions
-- ----------------------------
ALTER TABLE "public"."product_conditions" ADD CONSTRAINT "product_conditions_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table product_images
-- ----------------------------
CREATE INDEX "idx_product_images_product_sort" ON "public"."product_images" USING btree (
  "product_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "sort_order" "pg_catalog"."int4_ops" ASC NULLS LAST
);
CREATE UNIQUE INDEX "uq_product_images_primary_one" ON "public"."product_images" USING btree (
  "product_id" "pg_catalog"."int8_ops" ASC NULLS LAST
) WHERE is_primary = true;

-- ----------------------------
-- Primary Key structure for table product_images
-- ----------------------------
ALTER TABLE "public"."product_images" ADD CONSTRAINT "product_images_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table product_tags
-- ----------------------------
CREATE INDEX "idx_product_tags_tag" ON "public"."product_tags" USING btree (
  "tag_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table product_tags
-- ----------------------------
ALTER TABLE "public"."product_tags" ADD CONSTRAINT "product_tags_pkey" PRIMARY KEY ("product_id", "tag_id");

-- ----------------------------
-- Indexes structure for table products
-- ----------------------------
CREATE INDEX "idx_products_category_status_created" ON "public"."products" USING btree (
  "category_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "status" "pg_catalog"."enum_ops" ASC NULLS LAST,
  "created_at" "pg_catalog"."timestamptz_ops" DESC NULLS FIRST
);
CREATE INDEX "idx_products_desc_trgm" ON "public"."products" USING gin (
  "description" COLLATE "pg_catalog"."default" "public"."gin_trgm_ops"
);
CREATE INDEX "idx_products_seller_created" ON "public"."products" USING btree (
  "seller_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "created_at" "pg_catalog"."timestamptz_ops" DESC NULLS FIRST
);
CREATE INDEX "idx_products_status_created" ON "public"."products" USING btree (
  "status" "pg_catalog"."enum_ops" ASC NULLS LAST,
  "created_at" "pg_catalog"."timestamptz_ops" DESC NULLS FIRST
);
CREATE INDEX "idx_products_status_price" ON "public"."products" USING btree (
  "status" "pg_catalog"."enum_ops" ASC NULLS LAST,
  "price" "pg_catalog"."numeric_ops" ASC NULLS LAST
);
CREATE INDEX "idx_products_title_trgm" ON "public"."products" USING gin (
  "title" COLLATE "pg_catalog"."default" "public"."gin_trgm_ops"
);

-- ----------------------------
-- Triggers structure for table products
-- ----------------------------
CREATE TRIGGER "products_set_updated_at" BEFORE UPDATE ON "public"."products"
FOR EACH ROW
EXECUTE PROCEDURE "public"."trg_set_updated_at"();
CREATE TRIGGER "products_status_guard" BEFORE UPDATE ON "public"."products"
FOR EACH ROW
EXECUTE PROCEDURE "public"."trg_products_status_guard"();
COMMENT ON TRIGGER "products_status_guard" ON "public"."products" IS '约束商品状态机：ForSale↔Delisted 互转；ForSale→Sold 终态；禁止从 Sold 变更为其他状态（状态字段不可逆）。状态为 Sold 时，若不修改 status 字段，则允许更新其他非状态字段（如标题/描述/分类/标签），以便管理员纠错或数据清洗。';

-- ----------------------------
-- Checks structure for table products
-- ----------------------------
ALTER TABLE "public"."products" ADD CONSTRAINT "ck_products_price_positive" CHECK (price > 0::numeric);

-- ----------------------------
-- Primary Key structure for table products
-- ----------------------------
ALTER TABLE "public"."products" ADD CONSTRAINT "products_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table sessions
-- ----------------------------
CREATE INDEX "idx_sessions_expired_at" ON "public"."sessions" USING btree (
  "expired_at" "pg_catalog"."timestamptz_ops" ASC NULLS LAST
);
CREATE INDEX "idx_sessions_user" ON "public"."sessions" USING btree (
  "user_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Uniques structure for table sessions
-- ----------------------------
ALTER TABLE "public"."sessions" ADD CONSTRAINT "sessions_token_key" UNIQUE ("token");

-- ----------------------------
-- Primary Key structure for table sessions
-- ----------------------------
ALTER TABLE "public"."sessions" ADD CONSTRAINT "sessions_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table simple_users
-- ----------------------------
ALTER TABLE "public"."simple_users" ADD CONSTRAINT "simple_users_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table tags
-- ----------------------------
CREATE INDEX "idx_tags_category_id" ON "public"."tags" USING btree (
  "category_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Triggers structure for table tags
-- ----------------------------
CREATE TRIGGER "tags_set_updated_at" BEFORE UPDATE ON "public"."tags"
FOR EACH ROW
EXECUTE PROCEDURE "public"."trg_set_updated_at"();

-- ----------------------------
-- Uniques structure for table tags
-- ----------------------------
ALTER TABLE "public"."tags" ADD CONSTRAINT "tags_name_key" UNIQUE ("name");

-- ----------------------------
-- Primary Key structure for table tags
-- ----------------------------
ALTER TABLE "public"."tags" ADD CONSTRAINT "tags_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table test_products
-- ----------------------------
ALTER TABLE "public"."test_products" ADD CONSTRAINT "test_products_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table test_users
-- ----------------------------
ALTER TABLE "public"."test_users" ADD CONSTRAINT "test_users_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table user_recent_views
-- ----------------------------
CREATE INDEX "idx_views_user_time" ON "public"."user_recent_views" USING btree (
  "user_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "viewed_at" "pg_catalog"."timestamptz_ops" DESC NULLS FIRST,
  "id" "pg_catalog"."int8_ops" DESC NULLS FIRST
);

-- ----------------------------
-- Triggers structure for table user_recent_views
-- ----------------------------
CREATE TRIGGER "user_recent_views_prune_after_insert" AFTER INSERT ON "public"."user_recent_views"
FOR EACH ROW
EXECUTE PROCEDURE "public"."trg_prune_user_recent_views"();

-- ----------------------------
-- Primary Key structure for table user_recent_views
-- ----------------------------
ALTER TABLE "public"."user_recent_views" ADD CONSTRAINT "user_recent_views_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table users
-- ----------------------------
CREATE INDEX "idx_users_created_at" ON "public"."users" USING btree (
  "created_at" "pg_catalog"."timestamptz_ops" DESC NULLS FIRST
);

-- ----------------------------
-- Triggers structure for table users
-- ----------------------------
CREATE TRIGGER "users_set_updated_at" BEFORE UPDATE ON "public"."users"
FOR EACH ROW
EXECUTE PROCEDURE "public"."trg_set_updated_at"();

-- ----------------------------
-- Uniques structure for table users
-- ----------------------------
ALTER TABLE "public"."users" ADD CONSTRAINT "users_account_key" UNIQUE ("account");

-- ----------------------------
-- Checks structure for table users
-- ----------------------------
ALTER TABLE "public"."users" ADD CONSTRAINT "ck_users_account_format" CHECK (account::text ~ '^[A-Za-z0-9]+$'::text);
ALTER TABLE "public"."users" ADD CONSTRAINT "ck_users_password_len" CHECK (char_length(password_hash::text) >= 8);

-- ----------------------------
-- Primary Key structure for table users
-- ----------------------------
ALTER TABLE "public"."users" ADD CONSTRAINT "users_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Foreign Keys structure for table product_images
-- ----------------------------
ALTER TABLE "public"."product_images" ADD CONSTRAINT "product_images_product_id_fkey" FOREIGN KEY ("product_id") REFERENCES "public"."products" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table product_tags
-- ----------------------------
ALTER TABLE "public"."product_tags" ADD CONSTRAINT "product_tags_product_id_fkey" FOREIGN KEY ("product_id") REFERENCES "public"."products" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;
ALTER TABLE "public"."product_tags" ADD CONSTRAINT "product_tags_tag_id_fkey" FOREIGN KEY ("tag_id") REFERENCES "public"."tags" ("id") ON DELETE RESTRICT ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table products
-- ----------------------------
ALTER TABLE "public"."products" ADD CONSTRAINT "products_category_id_fkey" FOREIGN KEY ("category_id") REFERENCES "public"."categories" ("id") ON DELETE RESTRICT ON UPDATE NO ACTION;
ALTER TABLE "public"."products" ADD CONSTRAINT "products_condition_id_fkey" FOREIGN KEY ("condition_id") REFERENCES "public"."product_conditions" ("id") ON DELETE RESTRICT ON UPDATE NO ACTION;
ALTER TABLE "public"."products" ADD CONSTRAINT "products_seller_id_fkey" FOREIGN KEY ("seller_id") REFERENCES "public"."users" ("id") ON DELETE RESTRICT ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table sessions
-- ----------------------------
ALTER TABLE "public"."sessions" ADD CONSTRAINT "sessions_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table tags
-- ----------------------------
ALTER TABLE "public"."tags" ADD CONSTRAINT "fk_tags_category" FOREIGN KEY ("category_id") REFERENCES "public"."categories" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table user_recent_views
-- ----------------------------
ALTER TABLE "public"."user_recent_views" ADD CONSTRAINT "user_recent_views_product_id_fkey" FOREIGN KEY ("product_id") REFERENCES "public"."products" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;
ALTER TABLE "public"."user_recent_views" ADD CONSTRAINT "user_recent_views_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;

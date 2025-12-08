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

 Date: 08/12/2025 18:32:16
*/


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
COMMENT ON TABLE "public"."categories" IS '商品分类（由管理员维护）。删除被引用的分类将因外键而失败。';

-- ----------------------------
-- Records of categories
-- ----------------------------
INSERT INTO "public"."categories" VALUES (1, '数码电子', '各类手机、电脑及配件', '2025-12-06 11:28:34.396243+08', '2025-12-06 11:28:34.396243+08');
INSERT INTO "public"."categories" VALUES (2, '图书教材', '本科教材、考研考公资料', '2025-12-06 11:28:34.396243+08', '2025-12-06 11:28:34.396243+08');
INSERT INTO "public"."categories" VALUES (3, '生活日用', '宿舍神器、收纳、小家电', '2025-12-06 11:28:34.396243+08', '2025-12-06 11:28:34.396243+08');
INSERT INTO "public"."categories" VALUES (4, '运动器材', '各类球拍、健身器材', '2025-12-06 11:28:34.396243+08', '2025-12-06 11:28:34.396243+08');
INSERT INTO "public"."categories" VALUES (5, '珠宝首饰', '各类金银首饰，手链手环', '2025-12-06 16:08:57.253348+08', '2025-12-06 16:30:14.846178+08');
INSERT INTO "public"."categories" VALUES (6, '运动服饰', 'T恤速干衣瑜伽裤运动裤运动套装夹克/风衣卫衣/套头衫运动背心健身服运动内衣', '2025-12-06 16:31:26.53596+08', '2025-12-06 16:31:26.53596+08');
INSERT INTO "public"."categories" VALUES (7, '宠物食品', '狗粮猫粮宠物零食猫咪罐头狗狗罐头小宠食品水族食品', '2025-12-06 16:43:57.441053+08', '2025-12-06 16:52:23.063292+08');

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
COMMENT ON COLUMN "public"."product_conditions"."code" IS '新旧程度编码（唯一）。';
COMMENT ON COLUMN "public"."product_conditions"."name" IS '新旧程度中文名（如 全新/九成新/七成新 等）。';
COMMENT ON TABLE "public"."product_conditions" IS '商品新旧程度枚举表（发布页单选使用），此表为新旧程度的唯一事实来源。';

-- ----------------------------
-- Records of product_conditions
-- ----------------------------
INSERT INTO "public"."product_conditions" VALUES (1, 'BRAND_NEW', '全新', 10, '2025-12-05 22:11:21.909582+08', '2025-12-05 22:11:21.909582+08');
INSERT INTO "public"."product_conditions" VALUES (2, 'NINE_TENTHS', '九成新', 20, '2025-12-05 22:11:21.909582+08', '2025-12-05 22:11:21.909582+08');
INSERT INTO "public"."product_conditions" VALUES (3, 'EIGHT_TENTHS', '八成新', 30, '2025-12-05 22:11:21.909582+08', '2025-12-05 22:11:21.909582+08');
INSERT INTO "public"."product_conditions" VALUES (4, 'SEVEN_TENTHS', '七成新', 40, '2025-12-05 22:11:21.909582+08', '2025-12-05 22:11:21.909582+08');

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
COMMENT ON TABLE "public"."product_images" IS '商品图片表：一对多。主图通过 is_primary 标记并由唯一索引保证每商品最多一张主图；其余按 sort_order 排序。';

-- ----------------------------
-- Records of product_images
-- ----------------------------
INSERT INTO "public"."product_images" VALUES (26, 29, 'http://43.136.104.67:8080/uploads/2025/12/1764995182859.jpg', 1, 't');
INSERT INTO "public"."product_images" VALUES (27, 30, 'http://43.136.104.67:8080/uploads/2025/12/1765006043200.jpg', 1, 't');
INSERT INTO "public"."product_images" VALUES (28, 31, 'http://43.136.104.67:8080/uploads/2025/12/1765006266153.jpg', 1, 't');
INSERT INTO "public"."product_images" VALUES (29, 32, 'http://43.136.104.67:8080/uploads/2025/12/1765006433012.jpg', 1, 't');
INSERT INTO "public"."product_images" VALUES (30, 33, 'http://43.136.104.67:8080/uploads/2025/12/1765006514143.jpg', 1, 't');
INSERT INTO "public"."product_images" VALUES (31, 34, 'http://43.136.104.67:8080/uploads/2025/12/1765006612031.jpg', 1, 't');
INSERT INTO "public"."product_images" VALUES (32, 35, 'http://43.136.104.67:8080/uploads/2025/12/1765006665065.jpg', 1, 't');
INSERT INTO "public"."product_images" VALUES (33, 36, 'http://43.136.104.67:8080/uploads/2025/12/1765006760273.jpg', 1, 't');
INSERT INTO "public"."product_images" VALUES (34, 37, 'http://43.136.104.67:8080/uploads/2025/12/1765006829300.jpg', 1, 't');
INSERT INTO "public"."product_images" VALUES (35, 38, 'http://43.136.104.67:8080/uploads/2025/12/1765007022350.jpg', 1, 't');
INSERT INTO "public"."product_images" VALUES (37, 40, 'http://43.136.104.67:8080/uploads/2025/12/1765008346711.jpg', 1, 't');
INSERT INTO "public"."product_images" VALUES (38, 41, 'http://43.136.104.67:8080/uploads/2025/12/1765009175722.jpg', 1, 't');
INSERT INTO "public"."product_images" VALUES (40, 42, 'http://43.136.104.67:8080/uploads/2025/12/1765009610632.jpg', 1, 't');
INSERT INTO "public"."product_images" VALUES (41, 43, 'http://43.136.104.67:8080/uploads/2025/12/1765009787006.jpg', 1, 't');
INSERT INTO "public"."product_images" VALUES (42, 44, 'http://43.136.104.67:8080/uploads/2025/12/1765010089409.jpg', 1, 't');
INSERT INTO "public"."product_images" VALUES (43, 45, 'http://43.136.104.67:8080/uploads/2025/12/1765010147134.jpg', 1, 't');
INSERT INTO "public"."product_images" VALUES (44, 46, 'http://43.136.104.67:8080/uploads/2025/12/1765010195327.jpg', 1, 't');
INSERT INTO "public"."product_images" VALUES (45, 47, 'http://43.136.104.67:8080/uploads/2025/12/1765010254078.jpg', 1, 't');
INSERT INTO "public"."product_images" VALUES (46, 48, 'http://43.136.104.67:8080/uploads/2025/12/1765010303049.jpg', 1, 't');
INSERT INTO "public"."product_images" VALUES (47, 49, 'http://43.136.104.67:8080/uploads/2025/12/1765011192249.jpg', 1, 't');
INSERT INTO "public"."product_images" VALUES (48, 50, 'http://43.136.104.67:8080/uploads/2025/12/1765011252955.jpg', 1, 't');
INSERT INTO "public"."product_images" VALUES (49, 51, 'http://43.136.104.67:8080/uploads/2025/12/1765011383237.jpg', 1, 't');
INSERT INTO "public"."product_images" VALUES (50, 52, 'http://43.136.104.67:8080/uploads/2025/12/1765011430474.jpg', 1, 't');
INSERT INTO "public"."product_images" VALUES (51, 53, 'http://43.136.104.67:8080/uploads/2025/12/1765011472039.jpg', 1, 't');
INSERT INTO "public"."product_images" VALUES (64, 54, 'http://43.136.104.67:8080/uploads/2025/12/1765161366945.jpg', 1, 't');
INSERT INTO "public"."product_images" VALUES (65, 54, 'http://43.136.104.67:8080/uploads/2025/12/1765161366945.jpg', 2, 'f');
INSERT INTO "public"."product_images" VALUES (71, 55, 'http://43.136.104.67:8080/uploads/2025/12/1765162745112.jpg', 1, 't');
INSERT INTO "public"."product_images" VALUES (73, 39, 'http://43.136.104.67:8080/uploads/2025/12/1765007270015.jpg', 1, 't');

-- ----------------------------
-- Table structure for product_tags
-- ----------------------------
DROP TABLE IF EXISTS "public"."product_tags";
CREATE TABLE "public"."product_tags" (
  "product_id" int8 NOT NULL,
  "tag_id" int8 NOT NULL
)
;
COMMENT ON TABLE "public"."product_tags" IS '商品与标签的多对多关联表（复合主键防重复）。';

-- ----------------------------
-- Records of product_tags
-- ----------------------------
INSERT INTO "public"."product_tags" VALUES (29, 1);
INSERT INTO "public"."product_tags" VALUES (30, 1);
INSERT INTO "public"."product_tags" VALUES (31, 2);
INSERT INTO "public"."product_tags" VALUES (32, 2);
INSERT INTO "public"."product_tags" VALUES (33, 2);
INSERT INTO "public"."product_tags" VALUES (34, 6);
INSERT INTO "public"."product_tags" VALUES (35, 6);
INSERT INTO "public"."product_tags" VALUES (36, 6);
INSERT INTO "public"."product_tags" VALUES (37, 6);
INSERT INTO "public"."product_tags" VALUES (38, 7);
INSERT INTO "public"."product_tags" VALUES (40, 9);
INSERT INTO "public"."product_tags" VALUES (41, 14);
INSERT INTO "public"."product_tags" VALUES (42, 15);
INSERT INTO "public"."product_tags" VALUES (43, 13);
INSERT INTO "public"."product_tags" VALUES (44, 16);
INSERT INTO "public"."product_tags" VALUES (45, 17);
INSERT INTO "public"."product_tags" VALUES (46, 18);
INSERT INTO "public"."product_tags" VALUES (47, 19);
INSERT INTO "public"."product_tags" VALUES (48, 20);
INSERT INTO "public"."product_tags" VALUES (49, 22);
INSERT INTO "public"."product_tags" VALUES (50, 21);
INSERT INTO "public"."product_tags" VALUES (51, 23);
INSERT INTO "public"."product_tags" VALUES (52, 24);
INSERT INTO "public"."product_tags" VALUES (53, 25);
INSERT INTO "public"."product_tags" VALUES (54, 1);
INSERT INTO "public"."product_tags" VALUES (55, 1);
INSERT INTO "public"."product_tags" VALUES (39, 9);

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
COMMENT ON COLUMN "public"."products"."seller_id" IS '发布者用户 ID（1:N 关系：用户→商品）。';
COMMENT ON COLUMN "public"."products"."condition_id" IS '引用 product_conditions 表（唯一事实来源）；前端应使用 conditionId 作为入参，响应可返回 id 与名称/编码供展示。';
COMMENT ON COLUMN "public"."products"."status" IS '状态机：ForSale(在售) / Delisted(已下架) / Sold(已售-终态)。';
COMMENT ON COLUMN "public"."products"."main_image_url" IS '主图 URL 冗余字段，用于列表展示优化。发布/编辑/设置主图时需同步更新此字段。';
COMMENT ON TABLE "public"."products" IS '商品主表：每条记录代表一件实物（无库存字段）。';

-- ----------------------------
-- Records of products
-- ----------------------------
INSERT INTO "public"."products" VALUES (36, 1, '【当当专享印签版】明朝那些事儿1-9册单本&全套9册增补版 当年明月著作 二十四史中国明清史国民史学读本全本白话明史历史畅销书', '【当当专享印签版】明朝那些事儿1-9册单本&全套9册增补版 当年明月著作 二十四史中国明清史国民史学读本全本白话明史历史畅销书', 243.00, 1, 2, 'ForSale', 'http://43.136.104.67:8080/uploads/2025/12/1765006760273.jpg', '2025-12-06 15:39:20.313046+08', '2025-12-08 18:15:30.897917+08');
INSERT INTO "public"."products" VALUES (29, 7, '全新苹果 17', '全新苹果 17 标准版', 5999.00, 1, 1, 'Sold', 'http://43.136.104.67:8080/uploads/2025/12/1764995182859.jpg', '2025-12-06 12:26:22.867743+08', '2025-12-08 18:15:30.897917+08');
INSERT INTO "public"."products" VALUES (51, 2, '顽皮狗狗零食宠物鸭肉干鸡肉干磨牙棒幼犬训练奖励小型犬小狗零食', '顽皮狗狗零食宠物鸭肉干鸡肉干磨牙棒幼犬训练奖励小型犬小狗零食', 9.90, 1, 7, 'ForSale', 'http://43.136.104.67:8080/uploads/2025/12/1765011383237.jpg', '2025-12-06 16:56:23.276782+08', '2025-12-08 18:15:30.897917+08');
INSERT INTO "public"."products" VALUES (52, 2, '鲜朗猫罐头主食罐幼猫成猫湿粮猫咪零食无谷全价布偶美短加菲通用', '鲜朗猫罐头主食罐幼猫成猫湿粮猫咪零食无谷全价布偶美短加菲通用', 75.90, 1, 7, 'ForSale', 'http://43.136.104.67:8080/uploads/2025/12/1765011430474.jpg', '2025-12-06 16:57:10.487274+08', '2025-12-08 18:15:30.897917+08');
INSERT INTO "public"."products" VALUES (53, 2, '顽皮狗狗零食狗罐头营养湿粮拌饭增肥非主食牛肉宠物罐头24罐整箱', '顽皮狗狗零食狗罐头营养湿粮拌饭增肥非主食牛肉宠物罐头24罐整箱', 34.91, 1, 7, 'ForSale', 'http://43.136.104.67:8080/uploads/2025/12/1765011472039.jpg', '2025-12-06 16:57:52.069685+08', '2025-12-08 18:15:30.897917+08');
INSERT INTO "public"."products" VALUES (54, 7, '苹果 17', '苹果 17,16+256g，九成新', 5299.00, 2, 1, 'Delisted', 'http://43.136.104.67:8080/uploads/2025/12/1765161366945.jpg', '2025-12-08 10:36:06.959504+08', '2025-12-08 18:15:30.897917+08');
INSERT INTO "public"."products" VALUES (55, 7, '苹果 17，16+256g', '苹果 17 标准版，16+256g，九成新', 5099.00, 2, 1, 'ForSale', 'http://43.136.104.67:8080/uploads/2025/12/1765162745112.jpg', '2025-12-08 10:59:05.125914+08', '2025-12-08 18:15:30.897917+08');
INSERT INTO "public"."products" VALUES (39, 1, '天堂伞超大加大号抗风暴雨伞双人防风加固防晒折叠晴雨两用伞', '天堂伞超大加大号抗风暴雨伞双人防风加固防晒折叠晴雨两用伞男女', 40.10, 4, 3, 'ForSale', 'http://43.136.104.67:8080/uploads/2025/12/1765007270015.jpg', '2025-12-06 15:47:50.033139+08', '2025-12-08 18:15:30.897917+08');
INSERT INTO "public"."products" VALUES (37, 1, '【现货速发】活着 余华正版原著 当代文学小说兄弟第七天许三观卖血记卢克明的偷偷一笑平凡的世界在细雨中呼喊文学畅销书籍排行榜', '【现货速发】活着 余华正版原著 当代文学小说兄弟第七天许三观卖血记卢克明的偷偷一笑平凡的世界在细雨中呼喊文学畅销书籍排行榜', 19.30, 3, 2, 'ForSale', 'http://43.136.104.67:8080/uploads/2025/12/1765006829300.jpg', '2025-12-06 15:40:29.332216+08', '2025-12-08 18:15:30.897917+08');
INSERT INTO "public"."products" VALUES (38, 1, '麦朵尔MindDuo 2全光谱儿童书桌学生宿舍阅读学习专用LED护眼台灯', '麦朵尔MindDuo 2全光谱儿童书桌学生宿舍阅读学习专用LED护眼台灯', 1999.00, 1, 3, 'ForSale', 'http://43.136.104.67:8080/uploads/2025/12/1765007022350.jpg', '2025-12-06 15:43:42.544304+08', '2025-12-08 18:15:30.897917+08');
INSERT INTO "public"."products" VALUES (30, 1, 'Redmi Note 14 5G手机红米note手机小米手机小米官方旗舰店官网新品小米note14', '速度快，使用流畅，电池耐用', 801.27, 2, 1, 'ForSale', 'http://43.136.104.67:8080/uploads/2025/12/1765006043200.jpg', '2025-12-06 15:27:23.217214+08', '2025-12-08 18:15:30.897917+08');
INSERT INTO "public"."products" VALUES (31, 1, '【全国补立减500元】科大讯飞AI学习机S30 Turbo疯狂动物城特别定制版智能学生平板小中高英语平板官方旗舰店', '专业的学习平板，里面内容强大，初中生用，可以自行下载一些非游戏app，有家长监管', 4949.00, 3, 1, 'ForSale', 'http://43.136.104.67:8080/uploads/2025/12/1765006266153.jpg', '2025-12-06 15:31:06.486919+08', '2025-12-08 18:15:30.897917+08');
INSERT INTO "public"."products" VALUES (32, 1, '联想小新平板11 2025/小新Pad 11学习平板护眼娱乐办公学生平板大屏学习娱乐平板新品联想 333', '小新平板11 2025 天玑6300 8GB+128GB WIFI 深空灰 小新平板12.1 2025款 2.5k屏天玑6400 8GB+128', 674.00, 2, 1, 'ForSale', 'http://43.136.104.67:8080/uploads/2025/12/1765006433012.jpg', '2025-12-06 15:33:53.052456+08', '2025-12-08 18:15:30.897917+08');
INSERT INTO "public"."products" VALUES (33, 1, '【88券+政府补贴15%】联想小新平板Pro GT 11.1英寸3.2K超清屏骁龙8Gen3高性能AI平板电脑大学生平板小新平板', '浅海贝 信风灰 流岚粉 浅海贝（键盘笔套装） 信风灰（键盘笔套装）', 1614.00, 4, 1, 'ForSale', 'http://43.136.104.67:8080/uploads/2025/12/1765006514143.jpg', '2025-12-06 15:35:14.169425+08', '2025-12-08 18:15:30.897917+08');
INSERT INTO "public"."products" VALUES (34, 1, '从内耗到心流 复杂时代下的熵减行动指南 杨鸣 著 成功社科 新华书店正版图书籍 电子工业出版社', '从内耗到心流:复杂时代下的熵减行动指南', 27.00, 4, 2, 'ForSale', 'http://43.136.104.67:8080/uploads/2025/12/1765006612031.jpg', '2025-12-06 15:36:52.065472+08', '2025-12-08 18:15:30.897917+08');
INSERT INTO "public"."products" VALUES (35, 1, '正版 追风筝的人 中文版 作者胡塞尼著 学生课外读物书现当代文学中文小说文艺书籍 凤凰新华书店旗舰店', '正版 追风筝的人 中文版 作者胡塞尼著 学生课外读物书现当代文学中文小说文艺书籍 凤凰新华书店旗舰店', 20.25, 1, 2, 'ForSale', 'http://43.136.104.67:8080/uploads/2025/12/1765006665065.jpg', '2025-12-06 15:37:45.081802+08', '2025-12-08 18:15:30.897917+08');
INSERT INTO "public"."products" VALUES (40, 2, '天堂伞防晒防紫外线太阳伞学生便携折叠双人晴雨两用伞男女遮阳伞', '天堂伞防晒防紫外线太阳伞学生便携折叠双人晴雨两用伞男女遮阳伞', 24.80, 2, 3, 'ForSale', 'http://43.136.104.67:8080/uploads/2025/12/1765008346711.jpg', '2025-12-06 16:05:46.729937+08', '2025-12-08 18:15:30.897917+08');
INSERT INTO "public"."products" VALUES (41, 2, 's999银耳钉男高级感莫比乌斯养耳洞男生耳环女痞帅单只男士耳圈', 's999银耳钉男高级感莫比乌斯养耳洞男生耳环女痞帅单只男士耳圈', 29.90, 1, 5, 'ForSale', 'http://43.136.104.67:8080/uploads/2025/12/1765009175722.jpg', '2025-12-06 16:19:35.835813+08', '2025-12-08 18:15:30.897917+08');
INSERT INTO "public"."products" VALUES (45, 2, '迪卡侬跑步外套男秋季防风外套男卫衣跑步健身夹克运动外套SAX1', '迪卡侬跑步外套男秋季防风外套男卫衣跑步健身夹克运动外套SAX1', 199.90, 1, 6, 'ForSale', 'http://43.136.104.67:8080/uploads/2025/12/1765010147134.jpg', '2025-12-06 16:35:47.251204+08', '2025-12-08 18:15:30.897917+08');
INSERT INTO "public"."products" VALUES (42, 2, '【可刻字】卡蒂罗无限依恋情侣对戒一对设计戒指生日礼物送男女友', '【可刻字】卡蒂罗无限依恋情侣对戒一对设计戒指生日礼物送男女友', 199.90, 2, 5, 'ForSale', 'http://43.136.104.67:8080/uploads/2025/12/1765009610632.jpg', '2025-12-06 16:26:50.643952+08', '2025-12-08 18:15:30.897917+08');
INSERT INTO "public"."products" VALUES (43, 2, '情侣手链S999纯银一对高级感小众轻奢定制名字送男女友情人节礼物', '情侣手链S999纯银一对高级感小众轻奢定制名字送男女友情人节礼物', 96.90, 3, 5, 'ForSale', 'http://43.136.104.67:8080/uploads/2025/12/1765009787006.jpg', '2025-12-06 16:29:47.033709+08', '2025-12-08 18:15:30.897917+08');
INSERT INTO "public"."products" VALUES (44, 2, '特步山海丨休闲鞋男鞋冬季新款百搭厚底复古老爹鞋轻便户外运动鞋', '特步山海丨休闲鞋男鞋冬季新款百搭厚底复古老爹鞋轻便户外运动鞋', 210.33, 1, 6, 'ForSale', 'http://43.136.104.67:8080/uploads/2025/12/1765010089409.jpg', '2025-12-06 16:34:49.424892+08', '2025-12-08 18:15:30.897917+08');
INSERT INTO "public"."products" VALUES (46, 2, 'Basic Unit 基本单元 网眼弹力紧身运动T恤男款长袖健身衣服显大', 'Basic Unit 基本单元 网眼弹力紧身运动T恤男款长袖健身衣服显大', 202.10, 1, 6, 'ForSale', 'http://43.136.104.67:8080/uploads/2025/12/1765010195327.jpg', '2025-12-06 16:36:35.343528+08', '2025-12-08 18:15:30.897917+08');
INSERT INTO "public"."products" VALUES (47, 2, '准者速干运动长裤男士秋冬加绒保暖训练跑步健身宽松篮球休闲长裤', '准者速干运动长裤男士秋冬加绒保暖训练跑步健身宽松篮球休闲长裤', 68.11, 1, 6, 'ForSale', 'http://43.136.104.67:8080/uploads/2025/12/1765010254078.jpg', '2025-12-06 16:37:34.11688+08', '2025-12-08 18:15:30.897917+08');
INSERT INTO "public"."products" VALUES (48, 2, '小野和子鲨鱼裤女外穿2025新款秋冬保暖加厚瑜伽裤收腰提臀打底裤', '小野和子鲨鱼裤女外穿2025新款秋冬保暖加厚瑜伽裤收腰提臀打底裤', 79.90, 1, 6, 'ForSale', 'http://43.136.104.67:8080/uploads/2025/12/1765010303049.jpg', '2025-12-06 16:38:23.063207+08', '2025-12-08 18:15:30.897917+08');
INSERT INTO "public"."products" VALUES (49, 2, '网易严选三拼狗粮鸭肉梨冻干小型幼犬中大金毛边牧泰迪全价成犬粮', '网易严选三拼狗粮鸭肉梨冻干小型幼犬中大金毛边牧泰迪全价成犬粮', 220.20, 1, 7, 'ForSale', 'http://43.136.104.67:8080/uploads/2025/12/1765011192249.jpg', '2025-12-06 16:53:12.263366+08', '2025-12-08 18:15:30.897917+08');
INSERT INTO "public"."products" VALUES (50, 2, '弗列加特80%高鲜肉天然猫粮高蛋白营养猫咪主食幼成老专用粮官方', '弗列加特80%高鲜肉天然猫粮高蛋白营养猫咪主食幼成老专用粮官方', 139.00, 1, 7, 'ForSale', 'http://43.136.104.67:8080/uploads/2025/12/1765011252955.jpg', '2025-12-06 16:54:12.981617+08', '2025-12-08 18:15:30.897917+08');

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
COMMENT ON TABLE "public"."sessions" IS '可选：服务端会话/令牌黑名单存储；若使用纯 JWT + Redis，可不创建本表。';

-- ----------------------------
-- Records of sessions
-- ----------------------------

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

-- ----------------------------
-- Records of simple_users
-- ----------------------------

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
COMMENT ON TABLE "public"."tags" IS '商品标签库（由管理员维护，商品可多选标签）。';

-- ----------------------------
-- Records of tags
-- ----------------------------
INSERT INTO "public"."tags" VALUES (1, '手机', '2025-12-06 11:33:15.301066+08', '2025-12-06 11:33:15.301066+08', 1);
INSERT INTO "public"."tags" VALUES (2, '平板', '2025-12-06 11:33:15.301066+08', '2025-12-06 11:33:15.301066+08', 1);
INSERT INTO "public"."tags" VALUES (3, '耳机', '2025-12-06 11:33:15.301066+08', '2025-12-06 11:33:15.301066+08', 1);
INSERT INTO "public"."tags" VALUES (4, '教材', '2025-12-06 11:33:15.665251+08', '2025-12-06 11:33:15.665251+08', 2);
INSERT INTO "public"."tags" VALUES (5, '考研', '2025-12-06 11:33:15.665251+08', '2025-12-06 11:33:15.665251+08', 2);
INSERT INTO "public"."tags" VALUES (6, '小说', '2025-12-06 11:33:15.665251+08', '2025-12-06 11:33:15.665251+08', 2);
INSERT INTO "public"."tags" VALUES (7, '台灯', '2025-12-06 11:33:15.805128+08', '2025-12-06 11:33:15.805128+08', 3);
INSERT INTO "public"."tags" VALUES (8, '收纳', '2025-12-06 11:33:15.805128+08', '2025-12-06 11:33:15.805128+08', 3);
INSERT INTO "public"."tags" VALUES (9, '雨伞', '2025-12-06 11:33:15.805128+08', '2025-12-06 11:33:15.805128+08', 3);
INSERT INTO "public"."tags" VALUES (10, '球拍', '2025-12-06 11:33:15.867164+08', '2025-12-06 11:33:15.867164+08', 4);
INSERT INTO "public"."tags" VALUES (11, '瑜伽垫', '2025-12-06 11:33:15.867164+08', '2025-12-06 11:33:15.867164+08', 4);
INSERT INTO "public"."tags" VALUES (12, '滑板', '2025-12-06 11:33:15.867164+08', '2025-12-06 11:33:15.867164+08', 4);
INSERT INTO "public"."tags" VALUES (13, '银手链', '2025-12-06 16:11:25.403965+08', '2025-12-06 16:11:25.403965+08', 5);
INSERT INTO "public"."tags" VALUES (14, '银耳环', '2025-12-06 16:11:48.707688+08', '2025-12-06 16:11:48.707688+08', 5);
INSERT INTO "public"."tags" VALUES (15, '对戒', '2025-12-06 16:26:27.131823+08', '2025-12-06 16:26:27.131823+08', 5);
INSERT INTO "public"."tags" VALUES (16, '运动鞋', '2025-12-06 16:31:47.257168+08', '2025-12-06 16:31:47.257168+08', 6);
INSERT INTO "public"."tags" VALUES (17, '运动外套', '2025-12-06 16:31:59.036761+08', '2025-12-06 16:31:59.036761+08', 6);
INSERT INTO "public"."tags" VALUES (18, '运动T恤', '2025-12-06 16:32:18.577797+08', '2025-12-06 16:32:18.577797+08', 6);
INSERT INTO "public"."tags" VALUES (19, '速干运动裤', '2025-12-06 16:32:39.173889+08', '2025-12-06 16:32:39.173889+08', 6);
INSERT INTO "public"."tags" VALUES (20, '瑜伽裤', '2025-12-06 16:33:05.618496+08', '2025-12-06 16:33:05.618496+08', 6);
INSERT INTO "public"."tags" VALUES (21, '猫粮', '2025-12-06 16:44:08.519673+08', '2025-12-06 16:44:08.519673+08', 7);
INSERT INTO "public"."tags" VALUES (22, '狗粮', '2025-12-06 16:44:16.608+08', '2025-12-06 16:44:16.608+08', 7);
INSERT INTO "public"."tags" VALUES (23, '宠物零食', '2025-12-06 16:44:34.888934+08', '2025-12-06 16:44:34.888934+08', 7);
INSERT INTO "public"."tags" VALUES (24, '猫咪罐头', '2025-12-06 16:44:50.844142+08', '2025-12-06 16:44:50.844142+08', 7);
INSERT INTO "public"."tags" VALUES (25, '狗狗罐头', '2025-12-06 16:45:06.616594+08', '2025-12-06 16:45:06.616594+08', 7);

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

-- ----------------------------
-- Records of test_products
-- ----------------------------
INSERT INTO "public"."test_products" VALUES (1, 1, 'TestProduct_1764610630', '2025-12-02 01:37:10.78393');
INSERT INTO "public"."test_products" VALUES (2, 3, 'TestProductNoWechat_1764610630', '2025-12-02 01:37:10.807305');

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

-- ----------------------------
-- Records of test_users
-- ----------------------------
INSERT INTO "public"."test_users" VALUES (1, 'test_seller_1764610630', 'seller_wechat_123');
INSERT INTO "public"."test_users" VALUES (2, 'test_buyer_1764610630', 'buyer_wechat_456');
INSERT INTO "public"."test_users" VALUES (3, 'test_seller_no_wechat_1764610630', '');

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
COMMENT ON TABLE "public"."user_recent_views" IS '用户最近浏览商品记录（用于“猜你喜欢”）；按 user_id + viewed_at 倒序查询。';

-- ----------------------------
-- Records of user_recent_views
-- ----------------------------
INSERT INTO "public"."user_recent_views" VALUES (83, 7, 55, '2025-12-08 10:59:05.289747+08');
INSERT INTO "public"."user_recent_views" VALUES (84, 7, 55, '2025-12-08 10:59:16.882023+08');
INSERT INTO "public"."user_recent_views" VALUES (85, 7, 55, '2025-12-08 10:59:20.56783+08');
INSERT INTO "public"."user_recent_views" VALUES (86, 7, 55, '2025-12-08 11:22:45.877524+08');
INSERT INTO "public"."user_recent_views" VALUES (87, 7, 55, '2025-12-08 11:23:01.330802+08');
INSERT INTO "public"."user_recent_views" VALUES (88, 7, 55, '2025-12-08 11:23:12.220973+08');
INSERT INTO "public"."user_recent_views" VALUES (89, 7, 55, '2025-12-08 11:50:33.463081+08');
INSERT INTO "public"."user_recent_views" VALUES (90, 7, 51, '2025-12-08 11:50:51.900497+08');
INSERT INTO "public"."user_recent_views" VALUES (91, 7, 52, '2025-12-08 11:53:19.857512+08');
INSERT INTO "public"."user_recent_views" VALUES (92, 7, 55, '2025-12-08 11:53:33.17318+08');
INSERT INTO "public"."user_recent_views" VALUES (93, 1, 39, '2025-12-08 11:57:45.382678+08');
INSERT INTO "public"."user_recent_views" VALUES (94, 1, 39, '2025-12-08 11:58:06.845768+08');
INSERT INTO "public"."user_recent_views" VALUES (95, 1, 39, '2025-12-08 11:58:10.546213+08');
INSERT INTO "public"."user_recent_views" VALUES (96, 1, 39, '2025-12-08 11:58:11.348873+08');
INSERT INTO "public"."user_recent_views" VALUES (60, 1, 39, '2025-12-07 18:28:03.936014+08');
INSERT INTO "public"."user_recent_views" VALUES (61, 1, 39, '2025-12-08 08:41:38.85142+08');
INSERT INTO "public"."user_recent_views" VALUES (62, 1, 39, '2025-12-08 08:42:03.761785+08');
INSERT INTO "public"."user_recent_views" VALUES (63, 1, 51, '2025-12-08 09:13:34.797197+08');
INSERT INTO "public"."user_recent_views" VALUES (64, 1, 51, '2025-12-08 09:15:22.119817+08');
INSERT INTO "public"."user_recent_views" VALUES (65, 1, 53, '2025-12-08 09:22:07.789478+08');
INSERT INTO "public"."user_recent_views" VALUES (73, 7, 29, '2025-12-08 10:24:59.70085+08');
INSERT INTO "public"."user_recent_views" VALUES (74, 7, 29, '2025-12-08 10:26:27.395325+08');
INSERT INTO "public"."user_recent_views" VALUES (75, 7, 42, '2025-12-08 10:28:16.288338+08');
INSERT INTO "public"."user_recent_views" VALUES (76, 7, 54, '2025-12-08 10:36:07.174908+08');
INSERT INTO "public"."user_recent_views" VALUES (77, 7, 54, '2025-12-08 10:39:09.157341+08');
INSERT INTO "public"."user_recent_views" VALUES (78, 7, 54, '2025-12-08 10:39:19.278323+08');
INSERT INTO "public"."user_recent_views" VALUES (79, 7, 54, '2025-12-08 10:39:40.411122+08');
INSERT INTO "public"."user_recent_views" VALUES (80, 7, 54, '2025-12-08 10:57:33.911869+08');
INSERT INTO "public"."user_recent_views" VALUES (81, 7, 31, '2025-12-08 10:58:14.585843+08');
INSERT INTO "public"."user_recent_views" VALUES (82, 7, 50, '2025-12-08 10:58:19.40128+08');

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
INSERT INTO "public"."users" VALUES (1, 'czs123456', 'czs123456', '$2a$10$mSPRkSvgNKBXaEEb/n0WSuNKeZdKD7cMM4WEQJUlD7TjzP1dMNAA2', 'http://localhost:8080/uploads/2025/12/1764933578599.jpg', 'czs123456', 'f', NULL, '2025-12-04 11:03:13.50694+08', '2025-12-05 19:19:38.727693+08');
INSERT INTO "public"."users" VALUES (8, '123456', '12345678', '$2a$10$qllt0fUhRDJ1KjxiaXpD5Oxt0.iGu.W3Xc6on8LwHFEfVaXpEy73C', '', '123456', 'f', '2025-12-05 19:59:19.171266+08', '2025-12-05 19:59:04.318609+08', '2025-12-05 19:59:19.326148+08');
INSERT INTO "public"."users" VALUES (2, 'zrl123456', 'zrl123456', '$2a$10$cpw5U9iPEr0DfdsrLqM8zeQf7hxjU9OssRtj5ReeJovOoV81Ek9iq', 'http://localhost:8080/uploads/2025/12/1765009710371.jpg', 'zrl123456', 'f', NULL, '2025-12-04 11:06:43.495689+08', '2025-12-06 16:28:35.824637+08');
INSERT INTO "public"."users" VALUES (7, 'yycy134679', '云烟成雨 yycy', '$2a$10$322Qy97u9unXq2uU7cTNBO1LVSnEL3l2Ulk5XeB0E093BWot946jO', 'http://localhost:8080/uploads/2025/12/1765012571079.jpg', 'yycy134679', 'f', '2025-12-06 23:20:37.407821+08', '2025-12-05 19:39:04.513154+08', '2025-12-06 23:20:37.923034+08');

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

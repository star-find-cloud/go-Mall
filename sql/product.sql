use shop;

-- 删除并重建product表
DROP TABLE IF EXISTS product;

-- 重新启用外键检查
# SET FOREIGN_KEY_CHECKS=1;
CREATE TABLE IF NOT EXISTS product
(
    ID              BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    merchant_id     BIGINT COMMENT '商家ID',
    title           VARCHAR(255)        NOT NULL COMMENT '商品标题',
    sub_title       VARCHAR(255) COMMENT '副标题',
    brand           VARCHAR(255) COMMENT '品牌',
    product_sn      VARCHAR(255) UNIQUE NOT NULL COMMENT '唯一商品编号',
    cate_id         INT                 NOT NULL COMMENT '分类ID',
    click_count     INT        DEFAULT 0 COMMENT '点击量',
    purchase_count  BIGINT     DEFAULT 0 COMMENT '购买量',
    product_num     INT        DEFAULT 0 COMMENT '商品数量',
    price           DECIMAL(10, 2)      NOT NULL COMMENT '价格',
    market_price    DECIMAL(10, 2) COMMENT '市场价',
    attr            TEXT COMMENT '商品属性（建议JSON格式）',
    version         VARCHAR(50) COMMENT '商品版本',
    images          JSON COMMENT '图片ID数组',
    keywords        VARCHAR(255) COMMENT '搜索关键词',
    `desc`          TEXT COMMENT '商品描述',
    content         TEXT COMMENT '详细内容',
    specs           JSON COMMENT '商品规格JSON',
    is_deleted      TINYINT(1) DEFAULT 0 COMMENT '软删除标记',
    created_at      bigint(20) COMMENT '创建时间戳',
    updated_at      bigint(20) COMMENT '更新时间戳',
    deleted_at      bigint(20) COMMENT '删除时间戳',
    is_hot          TINYINT(1) DEFAULT 0 COMMENT '热门商品标记',
    is_best         TINYINT(1) DEFAULT 0 COMMENT '精品商品标记',
    is_new          TINYINT(1) DEFAULT 0 COMMENT '新品标记',
    is_booking      TINYINT(1) DEFAULT 0 COMMENT '预售标记',
    product_type_id INT COMMENT '商品类型ID',
    sort            INT        DEFAULT 0 COMMENT '排序权重',
    status          INT COMMENT '商品状态',

    -- 基础索引
    INDEX idx_shop (merchant_id) COMMENT '商家ID索引',
    INDEX idx_cate (cate_id) COMMENT '分类ID索引',
    INDEX idx_status (status) COMMENT '商品状态索引',
    INDEX idx_hot (is_hot) COMMENT '热门商品标记索引',
    INDEX idx_best (is_best) COMMENT '精品商品标记索引',
    INDEX idx_new (is_new) COMMENT '新品标记索引',
    INDEX idx_type (product_type_id) COMMENT '商品类型ID索引',
    INDEX idx_price (price) COMMENT '价格索引',
    INDEX idx_created (created_at) COMMENT '创建时间索引',
    INDEX idx_sort (sort) COMMENT '排序权重索引',
    INDEX idx_purchase (purchase_count) COMMENT '购买量索引',
    INDEX idx_click (click_count) COMMENT '点击量索引',

    -- 文本搜索索引
    INDEX idx_title (title(32)) COMMENT '商品标题前缀索引，用于标题搜索',
    INDEX idx_keywords (keywords(32)) COMMENT '关键词前缀索引，用于关键词搜索',

    -- 复合索引，用于常见的组合查询场景
    INDEX idx_cate_price (cate_id, price) COMMENT '分类+价格复合索引，用于按分类筛选并按价格排序',
    INDEX idx_cate_created (cate_id, created_at) COMMENT '分类+创建时间复合索引，用于按分类筛选并按时间排序',
    INDEX idx_shop_cate (merchant_id, cate_id) COMMENT '商家+分类复合索引，用于查询特定商家的特定分类商品',
    INDEX idx_status_sort (status, sort) COMMENT '状态+排序权重复合索引，用于按状态筛选并按权重排序',
    INDEX idx_recommend (status, is_hot, is_best, is_new) COMMENT '推荐商品复合索引，用于首页推荐商品查询',
    INDEX idx_shop_status (merchant_id, status) COMMENT '商家+状态复合索引，用于查询特定商家的特定状态商品',
    INDEX idx_deleted (is_deleted) COMMENT '软删除标记索引，用于过滤已删除商品'

) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COMMENT ='商品主表';

-- 创建全文索引，用于更高效的商品标题和关键词搜索（如果MySQL版本支持）
ALTER TABLE product ADD FULLTEXT INDEX ft_title_keywords (title, keywords) WITH PARSER ngram;

-- 创建商品销量统计视图（可选）
CREATE OR REPLACE VIEW v_product_sales AS
SELECT id,
       title,
       merchant_id,
       cate_id,
       price,
       purchase_count,
       click_count,
       (purchase_count * 0.7 + click_count * 0.3) AS popularity_score
FROM product
WHERE is_deleted = 0
  AND status = 1
ORDER BY popularity_score DESC;

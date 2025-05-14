create database if not exists shop;
use shop;

drop table if exists product;
CREATE TABLE product (
                         id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
                         title VARCHAR(255) NOT NULL DEFAULT '' COMMENT '商品标题',
                         sub_title VARCHAR(255) NOT NULL DEFAULT '' COMMENT '商品副标题',
                         product_sn VARCHAR(50) NOT NULL DEFAULT '' COMMENT '商品编号（唯一标识）',
                         cate_id INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '分类ID',
                         click_count INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '点击量',
                         product_num INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '商品库存数量',
                         price DECIMAL(10,2) UNSIGNED NOT NULL DEFAULT 0 COMMENT '销售价格',
                         market_price DECIMAL(10,2) UNSIGNED NOT NULL DEFAULT 0 COMMENT '市场参考价',
                         relation_product TEXT COMMENT '关联商品（JSON格式存储商品ID列表）',
                         product_attr TEXT COMMENT '商品属性（JSON格式存储规格参数）',
                         product_version VARCHAR(100) NOT NULL DEFAULT '' COMMENT '商品版本',
                         product_images JSON COMMENT '商品图片URL列表（JSON数组）',
                         product_gift TEXT COMMENT '商品赠品描述',
                         product_fitting TEXT COMMENT '商品配件清单',
                         product_color VARCHAR(50) NOT NULL DEFAULT '' COMMENT '商品颜色',
                         product_keywords VARCHAR(255) NOT NULL DEFAULT '' COMMENT 'SEO关键词',
                         product_desc TEXT COMMENT '商品简短描述',
                         product_content LONGTEXT COMMENT '商品详情（富文本HTML）',
                         is_deleted TINYINT(1) UNSIGNED NOT NULL DEFAULT 0 COMMENT '是否删除（0=未删除，1=已删除）',
                         created_at INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建时间（Unix时间戳）',
                         updated_at INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '更新时间（Unix时间戳）',
                         deleted_at INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '删除时间（Unix时间戳）',
                         is_hot TINYINT(1) UNSIGNED NOT NULL DEFAULT 0 COMMENT '是否热门（0=否，1=是）',
                         is_best TINYINT(1) UNSIGNED NOT NULL DEFAULT 0 COMMENT '是否精品（0=否，1=是）',
                         is_new TINYINT(1) UNSIGNED NOT NULL DEFAULT 0 COMMENT '是否新品（0=否，1=是）',
                         product_type_id INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '商品类型ID',
                         sort SMALLINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '排序权重（值越大越靠前）',
                         status TINYINT(1) UNSIGNED NOT NULL DEFAULT 1 COMMENT '商品状态（0=禁用，1=启用）',
                         PRIMARY KEY (id),
                         UNIQUE KEY uniq_product_sn (product_sn),
                         KEY idx_cate_id (cate_id),
                         KEY idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='商品信息表';

create table if not exists cart (
    id int(11) not null primary key auto_increment comment '主键id',
    title varchar(250) default '' comment '标题',
    price decimal(10,2) default '0.00' comment '价格',
    num int(11) default '0' comment '数量',
    goods_version varchar(50) default '' comment '商品版本',
    product_gift varchar(100) default '' comment '商品赠品',
    product_fitting varchar(100) default '' comment '商品搭配',
    product_color varchar(100) default '' comment '商品颜色',
    product_img varchar(100) default '' comment '商品图片',
    product_attr varchar(100) default '' comment '商品属性'
)engine = innodb default charset = utf8mb4 comment = '购物车表';


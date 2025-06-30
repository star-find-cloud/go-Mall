create database if not exists shop;
use shop;

-- 删除可能存在的外键约束
SET FOREIGN_KEY_CHECKS = 0;

# -- 尝试删除可能存在的外键约束（如果存在）
# ALTER TABLE `merchant_product` DROP FOREIGN KEY  `fk_merchant_product_merchant`;
# ALTER TABLE `merchant_product` DROP FOREIGN KEY  `fk_merchant_product_product`;
# ALTER TABLE `product` DROP FOREIGN KEY `fk_product_merchant`;

-- 删除商家相关表（如果存在）
DROP TABLE IF EXISTS `merchant_product`;
DROP TABLE IF EXISTS `merchant`;

-- 重新启用外键检查
SET FOREIGN_KEY_CHECKS = 1;

-- 创建商家表
CREATE TABLE IF NOT EXISTS `merchant`
(
    `id`               bigint(20)   NOT NULL AUTO_INCREMENT COMMENT '商家ID',
    `userID`           bigint(20)   NOT NULL COMMENT '用户ID',
    `name`             varchar(255) NOT NULL COMMENT '商家名称',
    `phone`            varchar(20)           DEFAULT NULL COMMENT '商家电话',
    `email`            varchar(100)          DEFAULT NULL COMMENT '商家邮箱',
    `password`         varchar(255) NOT NULL COMMENT '商家密码',
    `real_name`        varchar(50)           DEFAULT NULL COMMENT '商家真实姓名',
    `real_id`          varchar(20)           DEFAULT NULL COMMENT '商家身份证号',
    `license_image_id` int(11)               DEFAULT NULL COMMENT '营业执照图片ID',
    `old_name`         varchar(255)          DEFAULT NULL COMMENT '曾用名',
    `tag`              int(11)               DEFAULT '0' COMMENT '商家荣誉标签',
    `cate_id`          bigint(20)   NOT NULL DEFAULT '0' COMMENT '分类ID',
    `business_type`    json COMMENT '商家经营类型',
    `score`            float        NOT NULL DEFAULT '3.0' COMMENT '商家评分',
    `image_id`         bigint(20)            DEFAULT NULL COMMENT '商家图片ID',
    `create_at`        bigint(20)   NOT NULL COMMENT '创建时间',
    `update_at`        bigint(20)   COMMENT '更新时间',
    `delete_at`        bigint(20)            DEFAULT NULL COMMENT '删除时间',
    `status`           int(11)      NOT NULL DEFAULT '0' COMMENT '商家状态：0-待审核，1-正常，2-禁用',
    PRIMARY KEY (`id`),
    KEY `idx_status` (`status`),
    KEY `idx_score` (`score`),
    KEY `idx_cate_id` (`cate_id`),
    KEY `idx_name` (`name`(32)) COMMENT '商家名称前缀索引，用于模糊搜索优化',
    KEY `idx_phone` (`phone`) COMMENT '商家电话索引',
    KEY `idx_email` (`email`) COMMENT '商家邮箱索引',
    KEY `idx_real_name` (`real_name`(20)) COMMENT '商家真实姓名前缀索引，用于模糊搜索优化'
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT ='商家表';

-- 创建商家商品关联表（因为商家和商品是多对多关系）
CREATE TABLE IF NOT EXISTS `merchant_product`
(
    `merchantID`  bigint(20) NOT NULL AUTO_INCREMENT COMMENT '关联ID',
    `merchant_id` bigint(20) NOT NULL COMMENT '商家ID',
    `product_id`  bigint(20) NOT NULL COMMENT '商品ID',
    `create_at`   bigint(20) NOT NULL COMMENT '创建时间',
    PRIMARY KEY (`merchantID`),
    UNIQUE KEY `uk_merchant_product` (`merchant_id`, `product_id`),
    KEY `idx_product_id` (`product_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT ='商家商品关联表';


-- 创建复合索引，用于按分类和评分排序查询
CREATE INDEX `idx_cate_score` ON `merchant` (`cate_id`, `score` DESC);

-- 创建复合索引，用于按分类和创建时间排序查询
CREATE INDEX `idx_cate_create` ON `merchant` (`cate_id`, `create_at` DESC);

-- 创建全文索引，用于更高效的商家名称搜索（如果MySQL版本支持）
ALTER TABLE `merchant`
    ADD FULLTEXT INDEX `ft_name` (`name`);

-- 添加外键约束
ALTER TABLE `merchant_product`
    ADD CONSTRAINT `fk_merchant_product_merchant` FOREIGN KEY (`merchant_id`) REFERENCES `merchant` (`id`) ON DELETE CASCADE;
ALTER TABLE `merchant_product`
    ADD CONSTRAINT `fk_merchant_product_product` FOREIGN KEY (`product_id`) REFERENCES `product` (`id`) ON DELETE CASCADE;

-- 确保商品表中有shop_id字段并添加外键约束（如果商品表已存在）
-- 注意：取消注释以下语句前，请确保product表已存在
# ALTER TABLE `product` ADD CONSTRAINT `fk_product_merchant` FOREIGN KEY (`shop_id`) REFERENCES `merchant` (`merchantID`) ON DELETE SET NULL;

use shop;

drop table if exists orders;
CREATE TABLE `orders` (
                          `id` BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '订单ID',
                          `user_id` BIGINT NOT NULL COMMENT '用户ID',
                          `order_status` VARCHAR(50) NOT NULL COMMENT '订单状态',
                          `total_price` BIGINT NOT NULL DEFAULT 0 COMMENT '总价（单位：分）',
                          `pay_price` BIGINT NOT NULL DEFAULT 0 COMMENT '实际支付价格（单位：分）',
                          `created_at` BIGINT NOT NULL COMMENT '创建时间戳',
                          `updated_at` BIGINT NOT NULL COMMENT '更新时间戳',
                          `payment_id` BIGINT NOT NULL COMMENT '支付方式ID',
                          `shipping_id` BIGINT NOT NULL COMMENT '配送地址ID',
                          INDEX `idx_user_id` (`user_id`),
                          INDEX `idx_order_status` (`order_status`),
                          INDEX `created_at` (`created_at`),
                          INDEX `updated_at` (`updated_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='订单主表';

drop table if exists order_items;
CREATE TABLE `order_items` (
                               `item_id` BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '订单项ID',
                               `order_id` BIGINT NOT NULL COMMENT '关联订单ID',
                               `product_id` BIGINT NOT NULL COMMENT '商品ID',
                               `product_title` VARCHAR(255) NOT NULL COMMENT '商品名称',
                               `unit_price` BIGINT NOT NULL DEFAULT 0 COMMENT '单价（单位：分）',
                               `quantity` BIGINT NOT NULL DEFAULT 1 COMMENT '购买数量',
                               `subtotal` BIGINT NOT NULL DEFAULT 0 COMMENT '小计金额（单位：分）',
                               INDEX `idx_order_id` (`order_id`),
                               INDEX `idx_product_id` (`product_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='订单商品明细表';
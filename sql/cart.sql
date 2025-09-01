use shop;

drop table if exists cart;
CREATE TABLE cart (
                      id BIGINT PRIMARY KEY AUTO_INCREMENT, -- 分布式ID建议改用Snowflake/UUID
                      user_id VARCHAR(36) NOT NULL COMMENT '用户全局唯一标识',
                      created_at int NOT NULL COMMENT '创建时间',
                      updated_at int  COMMENT '更新时间'
);

drop table if exists cart_item;
CREATE TABLE cart_item (
                           id BIGINT PRIMARY KEY AUTO_INCREMENT,
                           cart_id BIGINT NOT NULL COMMENT '关联cart.id',
                           product_id BIGINT NOT NULL COMMENT '商品冗余快照ID',
                           product_title VARCHAR(255) NOT NULL COMMENT '加入时的商品标题快照',
                           create_price DECIMAL(10,2) NOT NULL COMMENT '加入时的价格快照',
                           now_price DECIMAL(10,2) NOT NULL COMMENT '现在价格',
                           product_image_oss VARCHAR(512) NOT NULL COMMENT 'OSS图片路径',
                           quantity INT NOT NULL DEFAULT 1 CHECK (quantity > 0) comment '商品数量',
                           specs JSON NOT NULL COMMENT '商品规格JSON快照',
                           added_at int  NOT NULL,
                           INDEX idx_cart_id (cart_id),
                           INDEX idx_product_id (product_id)
);

use shop;

drop table if exists inventory;
CREATE TABLE inventory
(
    product_id          BIGINT,
    available_stock     INT UNSIGNED NOT NULL comment '可用库存',
    reserved_stock      INT UNSIGNED NOT NULL comment '锁定库存',
    low_stock_threshold INT UNSIGNED NOT NULL comment '库存预警值',
    version             BIGINT DEFAULT 1, -- 乐观锁版本号（防超卖）
    create_at           int comment '创建时间',
    update_at           int comment '更新时间'
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;
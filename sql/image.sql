use shop;

drop table if exists images;
CREATE TABLE if not exists `images`
(
    `imageID`      bigint       NOT NULL PRIMARY KEY COMMENT '主键',
    `ownerType`    int          COMMENT '拥有者类型',
    `ownerID`      bigint       COMMENT '拥有者ID',
    `path`         VARCHAR(255) NOT NULL COMMENT '存储路径',
    `sha256hash`   varchar(255) comment '哈希值',
    `isCompressed` TINYINT(1)   DEFAULT 0 COMMENT '是否压缩 (0-否, 1-是)',
    `content_type` int          NOT NULL COMMENT '文件类型',
    `create_at`    timestamp  default current_timestamp comment '创建时间',
    `status`       int          DEFAULT 60 comment '状态 (60-正常,61-删除)',
    INDEX idx_owner (imageID, ownerType, ownerID)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;

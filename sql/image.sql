use shop;

drop table if exists images;
CREATE TABLE if not exists `images`
(
    `imageID`      bigint       NOT NULL PRIMARY KEY COMMENT '主键',
    `ownerType`    bigint       NOT NULL COMMENT '拥有者类型',
    `ownerID`      int          NOT NULL COMMENT '拥有者ID',
    `ossPath`      VARCHAR(255) NOT NULL COMMENT 'OSS存储路径',
    `SHA256Hash`   VARCHAR(64)  NOT NULL COMMENT 'SHA256哈希值',
    `isCompressed` TINYINT(1)   NOT NULL DEFAULT 0 COMMENT '是否压缩 (0-否, 1-是)',
    `file_path`          mediumblob   NOT NULL COMMENT '二进制数据',
    `content_type`        varchar(255) NOT NULL COMMENT '文件类型',
    INDEX idx_owner (ownerType, ownerID)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;

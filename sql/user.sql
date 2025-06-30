use shop;

drop table if exists user;
CREATE TABLE IF NOT EXISTS `user`
(
    `id`          INT NOT NULL AUTO_INCREMENT COMMENT '用户ID',
    `name`        VARCHAR(255) DEFAULT '' COMMENT '用户名',
    `password`    VARCHAR(255) DEFAULT '' COMMENT '密码',
    `email`       VARCHAR(100) DEFAULT '' COMMENT '邮箱',
    `phone`       VARCHAR(20)  DEFAULT '' COMMENT '手机号',
    `sex`         int          default 0 comment '性别',
    `tags`        json                   comment '标签',
    `create_time` bigint(20)   DEFAULT 0 COMMENT '创建时间',
    `update_time` bigint(20)   DEFAULT 0 COMMENT '更新时间',
    `delete_time` bigint(20)   DEFAULT 0 COMMENT '删除时间',
    `status`      INT          DEFAULT 0 COMMENT '状态',
    `last_ip`     VARCHAR(45)  DEFAULT '' COMMENT '最后登录IP',
    `image`       int          default 0 COMMENT '头像图片id',
    `is_vip`      TINYINT(1)   DEFAULT 0 COMMENT '是否VIP用户 (0: 否, 1: 是)',
    `role`        int          default 0 COMMENT '角色权限ID',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  AUTO_INCREMENT = 1000001
  DEFAULT CHARSET = utf8mb4 COMMENT = '用户表';

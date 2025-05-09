create database if not exists shop;
use shop;

create table if not exists  product(
    id int(10) not null auto_increment,
    title varchar(100) default '' comment '标题',
    sub_title varchar(100) default '' comment '子标题',
    product_sn varchar(50) default '' ,
    cate_id int(10) default '0' comment '分类id',
    click_count int(10) default '0' comment '点击量',
    product_num int(10) default '0' comment '商品编号',
    price decimal(10,2) default '0.00' comment '价格',
    market_price decimal(10,2) default '0.00' comment '市场价',
    relation_product varchar(100) default '' comment '关联商品',
    product_attr varchar(100) default '' comment '商品属性',
    product_version varchar(100) default '' comment '商品版本',
    product_img varchar(100) default '' comment '商品图片',
    product_gift varchar(100) default '' comment '商品赠品',
    product_fitting varchar(100) default '' comment '商品搭配',
    product_color varchar(100) default '' comment '商品颜色',
    product_keyword varchar(100) default '' comment '商品关键字',
    product_desc varchar(50) default '' comment '商品描述',
    product_content varchar(100) default '' comment '商品内容',
    product_type_id tinyint(4) default '0' comment '商品类型id',
    is_delete tinyint(4) default '0' comment '是否删除',
    is_hot tinyint(4) default '0' comment '是否热门',
    is_best tinyint(4) default '0' comment '是否畅销',
    is_new tinyint(4) default '0' comment '是否新品',
    add_time int(10) default '0' comment '添加时间',
    primary key (id)
)engine = innodb auto_increment=4 default charset = utf8mb4 comment = '商品表';

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


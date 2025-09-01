package database

import "github.com/jmoiron/sqlx"

type Database interface {
	// GetDB interface{}
	GetDB() *sqlx.DB
	// Ping 健康检查
	Ping() error
	// Close 关闭数据库连接
	Close() error
}

package database

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/star-find-cloud/star-mall/conf"
	log "github.com/star-find-cloud/star-mall/pkg/logger"
	"sync"
)

type MySQL struct {
	Conn *sqlx.DB
}

var (
	_mysql = &MySQL{}
	once   sync.Once
)

func initMysql() (*sqlx.DB, error) {
	var (
		_db *sqlx.DB
		err error
	)
	once.Do(func() {
		c := conf.GetConfig()
		user := c.Database.MySQL.User
		passwd := c.Database.MySQL.Password
		Host := c.Database.MySQL.MasterHost
		Port := c.Database.MySQL.MasterPort
		timeout := c.Database.MySQL.Timeout
		DSN := fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=utf8mb4&parseTime=True&timeout=%ss", user, passwd, Host, Port, timeout)
		_db, err = sqlx.Connect("mysql", DSN)
		if err != nil {
			fmt.Println("Database error, please check the logs.")
			//fmt.Println(c.Database.MySQL)
			log.MySQLLogger.Errorf("MySQL master connect faild: %s \n", err)
		} else {
			log.MySQLLogger.Infof("MySQL master connection successful: %s\n", Host)
		}
		_db.SetMaxOpenConns(c.Database.MySQL.MaxOpenConns)
		_db.SetMaxIdleConns(c.Database.MySQL.MaxIdleConns)
	})
	return _db, err
}

func NewMySQL() (*MySQL, error) {
	db, err := initMysql()
	_mysql.Conn = db
	return _mysql, err
}

func (sql MySQL) GetDB() *sqlx.DB {
	return sql.Conn
}

func (sql MySQL) Ping() error {
	return sql.Conn.Ping()
}

func (sql MySQL) Close() error {
	return sql.Conn.Close()
}

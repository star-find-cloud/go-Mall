package database

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/star-find-cloud/star-mall/conf"
	log "github.com/star-find-cloud/star-mall/pkg/logger"
	"go.uber.org/zap"
)

type MySQL struct {
	_db    *sqlx.DB
	logger *zap.SugaredLogger
}

var (
	master = &MySQL{}
	slaves = &[]MySQL{}
)

func init() {
	c := conf.GetConfig()
	user := c.Database.MySQL.User
	passwd := c.Database.MySQL.Password
	masterHost := c.Database.MySQL.MasterHost
	masterPort := c.Database.MySQL.MasterPort
	databasename := c.Database.MySQL.DBName
	timeout := c.Database.MySQL.Timeout
	masterDSN := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&timeout=%ss", user, passwd, masterHost, masterPort, databasename, timeout)
	m_db, err := sqlx.Connect("mysql", masterDSN)
	if err != nil {
		fmt.Println("Database error, please check the logs.")
		log.MySQLLogger.Errorf("MySQL master connect faild: %s \n", err)
	} else {
		log.MySQLLogger.Infof("MySQL master connection successful: %s\n", masterHost)
	}
	m_db.SetMaxOpenConns(c.Database.MySQL.MaxOpenConns)
	m_db.SetMaxIdleConns(c.Database.MySQL.MaxIdleConns)

	for _, slave := range c.Database.MySQL.Slaves {
		slaveDSN := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&timeout=%ss", user, passwd, slave.Host, slave.Port, databasename, timeout)
		s_db, err1 := sqlx.Connect("mysql", slaveDSN)
		if err1 != nil {
			log.MySQLLogger.Errorf("MySQL slaves connect faild: %s slaveIP: %s  \n", err, slave.Host)
		} else {
			log.MySQLLogger.Infof("MySQL slaves connection successful: %s\n", slave.Host)
		}
		s_db.SetMaxOpenConns(c.Database.MySQL.MaxOpenConns)
		s_db.SetMaxIdleConns(c.Database.MySQL.MaxIdleConns)

		if s_db != nil {
			*slaves = append(*slaves, MySQL{_db: s_db})
		}
	}

	master._db = m_db
}

func GetMySQLMaster() *sqlx.DB {
	return master._db
}

func GetMySQLSlaves() *sqlx.DB {
	for _, slave := range *slaves {
		return slave._db
	}
	return nil
}

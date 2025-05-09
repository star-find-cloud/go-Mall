package database

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/star-find-cloud/star-mall/conf"
	log "github.com/star-find-cloud/star-mall/pkg/logger"
	"go.uber.org/zap"
)

type Redis struct {
	rdb    *redis.Client
	logger *zap.SugaredLogger
}

var redisDB Redis

func init() {
	var c = conf.GetConfig()
	masterHost := c.Database.Redis.MasterHost
	masterport := c.Database.Redis.MasterPort
	passwd := c.Database.Redis.Password

	masterURL := fmt.Sprintf("redis://root:%s@%s:%d/0", passwd, masterHost, masterport)
	master, err := redis.ParseURL(masterURL)
	if err != nil {
		log.RedisLogger.Errorf("Redis Master connect faild: %s\n", err)
	} else {
		log.RedisLogger.Infof("Redis Master connect success: %s\n", masterURL)
	}
	redisDB.rdb = redis.NewClient(master)

	//slaveURL1 := fmt.Sprintf("redis://root:%s@%s:%d/0", passwd, slaveHost1, port)
	//slave1, err := redis.ParseURL(slaveURL1)
	//if err != nil {
	//	rdbLog.Errorf("Redis Slave1 connect faild: %s\n", err)
	//} else {
	//	rdbLog.Infof("Redis Slave1 connect success: %s\n", slaveURL1)
	//}
	//rdbSlave1.rdb = redis.NewClient(slave1)

	//slaveURL2 := fmt.Sprintf("redis://root:%s@%s:%d/0", passwd, slaveHost2, port)
	//slave2, err := redis.ParseURL(slaveURL2)
	//if err != nil {
	//	rdbLog.Errorf("Redis Slave2 connect faild: %s\n", err)
	//} else {
	//	rdbLog.Infof("Redis Slave2 connect success: %s\n", slaveURL2)
	//}
	//rdbSlave2.rdb = redis.NewClient(slave2)
}

func GetRedis() *redis.Client {
	return redisDB.rdb
}

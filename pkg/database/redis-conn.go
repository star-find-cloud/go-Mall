package database

import (
	"github.com/redis/go-redis/v9"
	"github.com/star-find-cloud/star-mall/conf"
	applog "github.com/star-find-cloud/star-mall/pkg/logger"
	"time"
)

type Redis struct {
	rdb        *redis.Client
	slaves     []*redis.Client
	slaveIndex uint32
}

var _redisDB = &Redis{}

func init() {
	var c = conf.GetConfig()
	masterAddr := c.Database.Redis.MasterHost + ":" + c.Database.Redis.MasterPort

	rdb := redis.NewClient(&redis.Options{
		Addr:     masterAddr,
		Password: c.Database.Redis.Password,
		DB:       c.Database.Redis.Database,
		PoolSize: c.Database.Redis.PoolSize,
		// 连接超时时间
		DialTimeout:    c.Database.Redis.DialTimeout * time.Second,
		ReadTimeout:    c.Database.Redis.ReadTimeout * time.Second,
		MaxIdleConns:   c.Database.Redis.MaxIdleConns,
		MaxActiveConns: c.Database.Redis.MaxActiveConns,
	})

	slaves := make([]*redis.Client, 0)
	for _, slaveConf := range c.Database.Redis.Slaves {
		slaveAddr := slaveConf.Host + ":" + slaveConf.Port
		slave := redis.NewClient(&redis.Options{
			Addr:     slaveAddr,
			Password: c.Database.Redis.Password,
			DB:       0,
			PoolSize: 100,
			// 连接超时时间
			DialTimeout:    10 * time.Second,
			ReadTimeout:    10 * time.Second,
			MaxIdleConns:   50,
			MaxActiveConns: 200,
		})
		slaves = append(slaves, slave)
	}

	_redisDB = &Redis{
		rdb:    rdb,
		slaves: slaves,
	}
}

func GetRedis() *Redis {
	return _redisDB
}

func (r *Redis) GetWriteCache() *redis.Client {
	if _redisDB.rdb != nil {
		applog.RedisLogger.Infoln("redis connect success")
	} else {
		applog.RedisLogger.Errorln("redis connect failed")
	}
	return _redisDB.rdb
}

func (r *Redis) GetReadCache() *redis.Client {
	//	if len(_redisDB.slaves) == 0 {
	//		return _redisDB.rdb
	//	}
	//	index := atomic.AddUint32(&r.slaveIndex, 1) % uint32(len(r.slaves))
	//	if r.slaves[index] == nil {
	//		applog.RedisLogger.Errorln("slave connect failed")
	//	}
	return _redisDB.rdb
}

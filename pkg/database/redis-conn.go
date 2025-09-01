package database

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/star-find-cloud/star-mall/conf"
	log "github.com/star-find-cloud/star-mall/pkg/logger"
	"sync"

	"time"
)

type Redis struct {
	//rdb *redis.ClusterClient
	Cache *redis.Client
}

var (
	_redis    = &Redis{}
	redisOnce sync.Once
)

// 生产环境取消以下注释, 改用集群模式
func initRedis() (*redis.Client, error) {
	var (
		rdb *redis.Client
		err error
	)
	redisOnce.Do(func() {
		var (
			c = conf.GetConfig()
			//addrs []string
		)
		masterAddr := c.Database.Redis.MasterHost + ":" + c.Database.Redis.MasterPort
		//for _, slave := range c.Database.Redis.Slaves {
		//	slaveAddr := slave.Host + ":" + slave.Port
		//	addrs = append(addrs, slaveAddr)
		//}
		//addrs = append(addrs, masterAddr)

		//rdb := redis.NewClusterClient(&redis.ClusterOptions{
		//	Addrs:    addrs,
		//	Password: c.Database.Redis.Password,
		//	PoolSize: c.Database.Redis.PoolSize,
		//	// 将读命令发送到从节点
		//	ReadOnly: true,
		//	// 优先选择延迟最低的节点
		//	RouteByLatency: true,
		//	// 连接超时时间
		//	DialTimeout:    c.Database.Redis.DialTimeout * time.Second,
		//	WriteTimeout:   c.Database.Redis.WriteTimeout * time.Second,
		//	ReadTimeout:    c.Database.Redis.ReadTimeout * time.Second,
		//	MaxIdleConns:   c.Database.Redis.MaxIdleConns,
		//	MaxActiveConns: c.Database.Redis.MaxActiveConns,
		//})

		rdb = redis.NewClient(&redis.Options{
			Addr:         masterAddr,
			Password:     c.Database.Redis.Password,
			PoolSize:     c.Database.Redis.PoolSize,
			DB:           c.Database.Redis.Database,
			DialTimeout:  c.Database.Redis.DialTimeout * time.Second,
			ReadTimeout:  c.Database.Redis.ReadTimeout * time.Second,
			WriteTimeout: c.Database.Redis.WriteTimeout * time.Second,
		})

		err = rdb.Ping(context.Background()).Err()
		if err != nil {
			log.RedisLogger.Errorf("redis连接失败: %v", err)
		}
	})
	return rdb, err
}

func NewRedis() (*Redis, error) {
	cache, err := initRedis()
	_redis.Cache = cache
	return _redis, err
}

//func (r *Redis) GetCache() *redis.ClusterClient {
//	return _redis.rdb
//}

func (r Redis) GetCache() *redis.Client {
	return _redis.Cache
}

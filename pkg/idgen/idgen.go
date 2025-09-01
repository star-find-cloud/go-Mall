package idgen

import (
	"context"
	"github.com/bwmarrin/snowflake"
	"github.com/star-find-cloud/star-mall/conf"
	"github.com/star-find-cloud/star-mall/pkg/database"
	log "github.com/star-find-cloud/star-mall/pkg/logger"
	"strconv"
	"sync"
)

var Node = &snowflake.Node{}

func init() {
	var c = conf.GetConfig()
	snowNode, err := snowflake.NewNode(c.App.NodeID)
	if err != nil {
		log.AppLogger.Fatalf("Failed to create snowflake node: %v/n", err)

	}
	Node = snowNode
}

// GenerateUid 生成 uid
func GenerateUid() (int64, error) {
	return Node.Generate().Int64(), nil
}

// GetUid 获取 uid
func GetUid(ctx context.Context, cache *database.Redis) (int64, error) {
	// 弹出 uid
	uidStr, err := cache.GetCache().LPop(ctx, "uid_pool").Result()
	uid, err := strconv.ParseInt(uidStr, 10, 64)
	if err == nil {
		return uid, nil
	}

	uids := generateUIDs()
	strUids := make([]string, len(uids))
	for i, uid := range uids {
		strUids[i] = strconv.FormatInt(uid, 10) // int64 转 string
	}
	err = cache.GetCache().RPush(ctx, "uid_pool", strUids).Err()
	if err != nil {
		log.RedisLogger.Errorf("generate uids error: %v", err)
		return 0, err
	}

	uidStr, err = cache.GetCache().LPop(ctx, "uid_pool").Result()
	uid, err = strconv.ParseInt(uidStr, 10, 64)
	if err != nil {
		log.RedisLogger.Errorf("get uid error: %v", err)
		return 0, err
	}
	return uid, nil
}

func generateUIDs() []int64 {
	var (
		c        = conf.GetConfig()
		uidCount = c.App.UidCount
		uids     = make([]int64, 0, uidCount)
		wg       = &sync.WaitGroup{}
		mu       = &sync.Mutex{}
		errChan  = make(chan error)
	)

	wg.Add(uidCount)
	for i := 0; i < uidCount; i++ {
		go func() {
			defer wg.Done()

			// 生成 uid
			var uid, err = GenerateUid()
			if err != nil {
				errChan <- err
				return
			}

			// 加锁保证并发安全
			mu.Lock()
			uids = append(uids, uid)
			mu.Unlock()
		}()
	}

	go func() {
		for err := range errChan {
			log.AppLogger.Errorf("UID generation failed: %v/n", err)
		}
	}()

	wg.Wait()
	close(errChan)
	return uids
}

// MakeUidCache 创建 uid 缓存
func MakeUidCache(ctx context.Context, uid int64, cache *database.Redis, metadata map[string]interface{}) error {
	key := "pending_uid:" + strconv.FormatInt(uid, 10)
	return cache.GetCache().HSet(ctx, key, metadata).Err()
}

func GetUidCache(ctx context.Context, uid int64, cache *database.Redis) (map[string]string, error) {
	key := "pending_uid:" + strconv.FormatInt(uid, 10)
	return cache.GetCache().HGetAll(ctx, key).Result()
}

func DeleteUidCache(ctx context.Context, uid int64, cache *database.Redis) error {
	key := "pending_uid:" + strconv.FormatInt(uid, 10)
	return cache.GetCache().Del(ctx, key).Err()
}

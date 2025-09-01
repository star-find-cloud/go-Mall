package database

import (
	"context"
	"fmt"
	"github.com/star-find-cloud/star-mall/conf"
	"testing"
)

func TestRedis_GetCache(t *testing.T) {
	var c = conf.GetConfig()
	fmt.Println(c.Database.Redis)

	rdb := &Redis{}
	cache := rdb.GetCache()
	fmt.Println("1", cache)
	ctx := context.Background()
	cache.Set(ctx, "k2", "v2", 0)
	v, err := cache.Get(ctx, "k2").Result()
	fmt.Println(v, err)
}

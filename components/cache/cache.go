package cache

import (
	"github.com/foolishnoob/go-xkratos/config"
	"github.com/go-redis/redis/v8"
)

type Cache interface {
	GetInterface() *redis.Client
}

type cache struct {
	Cache
	client *redis.Client
}

func (c *cache) GetInterface() *redis.Client {
	return c.client
}

func NewCache(conf *config.BootConfig) Cache {
	if "" == conf.GetRedis().GetAddr() {
		return nil
	}
	var option = redis.Options{
		Addr: conf.GetRedis().GetAddr(),
	}
	return &cache{
		client: redis.NewClient(&option),
	}
}

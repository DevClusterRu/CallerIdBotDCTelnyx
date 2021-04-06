package engine

import (
	"context"
	"github.com/go-redis/redis/v8"
)

type RedisObject struct {
	RD *redis.Client
}

func (f *RedisObject) Init()  {
	f.RD =  redis.NewClient(&redis.Options{
		Addr:     Conf.Variables["REDIS_HOST"] + ":" + Conf.Variables["REDIS_PORT"],
		Password: Conf.Variables["REDIS_PASSWORD"], // no password set
		DB:       0,                                // use default DB
	})
}

func (f *RedisObject) GetRow(header string) string  {
	val, _ := f.RD.Keys(context.Background(), header+"*").Result()
	if len(val) == 0 {
		return ""
	}
	v := val[0]
	value:=f.RD.Get(context.Background(), val[0]).Val()
	f.RD.Del(context.Background(), v).Result()
	return value
}

func (f *RedisObject) ClearByKeys(keys []string)  {
	for _,key:=range keys{
		val, _ := f.RD.Keys(context.Background(), key+"*").Result()
		for _, v := range val {
			f.RD.Del(context.Background(), v).Result()
		}
	}
}
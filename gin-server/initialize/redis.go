package initialize

import (
	"fmt"
	"github.com/go-redis/redis"
	"slb-admin/global"
)

func Redis() {
	redisCfg := global.CONFIG.Redis
	client := redis.NewClient(&redis.Options{
		Addr:     redisCfg.Addr,
		Password: redisCfg.Password, // no password set
		DB:       redisCfg.DB,       // use default DB
	})
	_, err := client.Ping().Result()
	if err != nil {
		fmt.Println(err)
	} else {
		global.REDIS = client
	}
}

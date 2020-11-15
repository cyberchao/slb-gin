package service

import (
	"slb-admin/global"
	"time"
)

func SetRedisJWT(jwt string, userName string) (err error) {
	timer := 60 * 60 * 24 * 3 * time.Second
	err = global.REDIS.Set(userName, jwt, timer).Err()
	return err
}

func GetRedisJWT(userName string) (err error, redisJWT string) {
	redisJWT, err = global.REDIS.Get(userName).Result()
	return err, redisJWT
}

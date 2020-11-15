package global

import (
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"slb-admin/config"
)

var (
	DB     *gorm.DB
	REDIS  *redis.Client
	VP     *viper.Viper
	CONFIG config.Server
)

package global

import (
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
	"slb-admin/config"
)

var (
	DB     *gorm.DB
	REDIS  *redis.Client
	Mogo   *mongo.Client
	VP     *viper.Viper
	CONFIG config.Server
)

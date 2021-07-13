package main

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/go-redis/redis"
	"github.com/robfig/cron"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/natefinch/lumberjack.v2"
)

// set结构:key=**nginx-ssl_uat_bx:m.yqb.com** value=**['21.64.24.23:8080:1','21.64.24.24:8080:1']**

const (
	username    = "admin"
	password    = "123"
	host        = "10.0.0.10"
	port        = "27017"
	dbname      = "slb"
	redis_pwd   = "Qwer4321!"
	redis_host  = "127.0.0.1"
	redis_port  = "6379"
	redis_db    = 1
	logfile     = "/Users/pangru/Documents/slb-admin/gin-server/tmp/health_check/syncer.log"
	interval    = "5"
	tcp_timeout = 2
)

var (
	Redis  *redis.Client
	Mongo  *mongo.Client
	Logger *zap.SugaredLogger
)

type ServerInfo struct {
	Ip     string `json:"ip"`
	Port   string `json:"port"`
	Weight string `json:"weight"`
	Status string `json:"status"`
}

type UpstreamDoc struct {
	Env        string       `json:"env"`
	Cluster    string       `json:"cluster"`
	Name       string       `json:"name"`
	ServerList []ServerInfo `json:"serverList"`
	Forward    string       `json:"forward"'`
	Version    int          `json:"version"`
	Time       time.Time    `json:"time"`
	FilePath   string       `json:"filepath"`
}

type HostEl struct {
	Key    string
	Member string
}

func InitLogger() {
	writeSyncer := getLogWriter()
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)

	logger := zap.New(core, zap.AddCaller())
	Logger = logger.Sugar()
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename: logfile,
		Compress: false,
	}
	return zapcore.AddSync(lumberJackLogger)
}

func initRedis() {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redis_host, redis_port),
		Password: redis_pwd, // no password set
		DB:       redis_db,  // use default DB
	})
	_, err := client.Ping().Result()
	if err != nil {
		Logger.Errorf("redis ping failed [%s]", err.Error())
	} else {
		Redis = client
	}
}

func initMongodb() {
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s", username, password, host, port, dbname)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// Connect to MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		Logger.Errorf("mongodb connect [%s] failed [%s]", uri, err.Error())
	} else {
		if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
			Logger.Errorf("mongodb ping [%s] failed [%s]", uri, err.Error())
		}
		Mongo = client
	}
}

func connect(key string, serverinfo ServerInfo, ch chan<- HostEl) {
	var ipPort string
	var hostEl HostEl
	var member string

	if serverinfo.Status == "up" {
		ipPort = fmt.Sprintf("%s:%s", serverinfo.Ip, serverinfo.Port)
	} else {
		hostEl.Key = ""
		hostEl.Member = ""
		ch <- hostEl
		return
	}

	timeout := time.Second * tcp_timeout
	_, err := net.DialTimeout("tcp", ipPort, timeout)

	if err != nil {
		hostEl.Key = ""
		hostEl.Member = ""
		ch <- hostEl
		// Logger.Errorf("upstream[%s] tcp connect failed[%s]", key, serverinfo)
		return
	} else {
		member = fmt.Sprintf("%s:%s:%s", serverinfo.Ip, serverinfo.Port, serverinfo.Weight)
		hostEl.Key = key
		hostEl.Member = member
		ch <- hostEl
		Logger.Infof("upstream[%s] tcp connect succeed[%s]", key, serverinfo)
		return
	}
}

func mongo2redis() {

	pipe := Redis.Pipeline()
	collection := Mongo.Database("slb").Collection("upstream")
	filterCursor, _ := collection.Find(context.TODO(), bson.M{})

	var upstreamList []UpstreamDoc
	if err := filterCursor.All(context.TODO(), &upstreamList); err != nil {

		Logger.Errorf("get upstream info from mongodb failed[%s]", err.Error())
	}

	// set结构:key=**nginx-ssl_uat_bx:m.yqb.com** value=**['21.64.24.23:8080:1','21.64.24.24:8080:1']**
	ch := make(chan HostEl)
	for _, upstream := range upstreamList {
		key := fmt.Sprintf("%s_%s:%s", upstream.Cluster, upstream.Env, upstream.Name)
		for _, serverinfo := range upstream.ServerList {
			go connect(key, serverinfo, ch)
		}
	}
	for _, upstream := range upstreamList {
		for range upstream.ServerList {
			hostEl := <-ch
			if hostEl.Key != "" {
				pipe.SAdd(hostEl.Key, hostEl.Member)
			}
		}
	}
	Redis.FlushDB()
	_, err := pipe.Exec()
	if err != nil {
		Logger.Errorf("pipeline exec failed[%s]", err.Error())
	}
}

func main() {
	initRedis()
	initMongodb()
	InitLogger()
	// mongo2redis()
	cron := cron.New()

	cron.AddFunc(fmt.Sprintf("@every %ss", interval), func() {
		Logger.Infof("start cron [%s]", time.Now())
		mongo2redis()
		Logger.Infof("end cron [%s]", time.Now())
	})

	cron.Start()
	fmt.Println("server stated")
	// 保持主goroutine运行
	for {
		time.Sleep(time.Second * 10)
	}
}

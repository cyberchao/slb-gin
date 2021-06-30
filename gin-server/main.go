package main

import (
	"net/http"
	"slb-admin/core"
	"slb-admin/global"
	"slb-admin/initialize"
)

func main() {
	// 初始化配置文件
	global.VP = core.Viper("./config.yaml")
	// 初始化mongodb
	global.Mogo = initialize.Mongo()
	// 初始化mysql
	global.DB = initialize.GormMysql()
	initialize.MysqlTables(global.DB)
	db, _ := global.DB.DB()
	defer db.Close()
	// 初始化redis
	initialize.Redis()
	// 初始化日志
	initialize.InitLogger()
	defer global.Logger.Sync()
	// 初始化路由
	Router := initialize.Routers()
	// 启动http服务
	http.ListenAndServe(":8080", Router)
}

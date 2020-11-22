package main

import (
	"net/http"
	"slb-admin/core"
	"slb-admin/global"
	"slb-admin/initialize"
)
// @title Swagger Example API
// @version 0.0.1
// @description This is a sample Server pets
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name x-token
// @BasePath /
func main() {
	global.VP = core.Viper("./config.yaml")
	global.DB = initialize.GormMysql()
	initialize.MysqlTables(global.DB)
	initialize.Redis()
	db, _ := global.DB.DB()
	defer db.Close()

	Router := initialize.Routers()
	//s := endless.NewServer(":8080", Router)
	//s.ReadHeaderTimeout = 10 * time.Millisecond
	//s.WriteTimeout = 10 * time.Second
	//s.MaxHeaderBytes = 1 << 20
	//s.ListenAndServe()
	//Router.Run(":8082")
	http.ListenAndServe(":8080", Router)
}

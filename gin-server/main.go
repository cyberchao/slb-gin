package main

import (
	"slb-admin/core"
	"slb-admin/global"
	"slb-admin/initialize"
)

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
	Router.Run(":8080")
	//http.ListenAndServe(":8080", Router)
}
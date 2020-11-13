package main

import (
	"github.com/fvbock/endless"
	"slb-admin/initialize"
	"time"
)

func main() {
	Router := initialize.Routers()
	s := endless.NewServer(":8081", Router)
	s.ReadHeaderTimeout = 10 * time.Millisecond
	s.WriteTimeout = 10 * time.Second
	s.MaxHeaderBytes = 1 << 20
	s.ListenAndServe()
}

package model

import "time"

type Host struct {
	ID          int    `json:"id" gorm:"AUTO_INCREMENT"`
	Cluster     string `json:"cluster" gorm:"cluster:集群"`
	Env         string `json:"env" gorm:"env:环境"`
	Ip          string `json:"ip" gorm:"ip:ip"`
	Description string `json:"description" gorm:"description:描述"`
	Reload_time string `json:"reload_time" gorm:"reload_time:reload时间"`
}

type VhostDoc struct {
	Env         string      `json:"env"`
	Cluster     string      `json:"cluster"`
	Ngx         interface{} `json:"ngx"`
	Src         string      `json:"src"`
	Description string      `json:"description"'`
	Version     int         `json:"version"`
	Time        time.Time   `json:"time"`
	FilePath    string      `json:"filepath"`
}

type ServerList struct {
	Ip     string `json:"ip"`
	Port   string `json:"port"`
	Weight string `json:"weight"`
	Status string `json:"status"`
}

type UpstreamDoc struct {
	Env        string       `json:"env"`
	Cluster    string       `json:"cluster"`
	Name       string       `json:"name"`
	ServerList []ServerList `json:"serverList"`
	Forward    string       `json:"forward"'`
	Version    int          `json:"version"`
	Time       time.Time    `json:"time"`
	FilePath   string       `json:"filepath"`
}

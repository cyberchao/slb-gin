package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	username = "admin"
	password = "123"
	host     = "10.0.0.10"
	port     = "27017"
	dbname   = "slb"
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

var Mongo *mongo.Client

func initMongodb() {
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s", username, password, host, port, dbname)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// Connect to MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	} else {
		if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
			panic(err)
		}
		Mongo = client
	}
}

func main() {
	initMongodb()
	collection := Mongo.Database("slb").Collection("upstream")
	for k := 1; k <= 8000; k++ {
		var serverinfo ServerInfo
		serverinfo.Ip = "10.0.0.10"
		serverinfo.Port = strconv.Itoa(k)
		serverinfo.Weight = "1"
		serverinfo.Status = "up"

		var upstreamdoc UpstreamDoc
		upstreamdoc.Env = "uat"
		upstreamdoc.Cluster = "nginx-bc"
		upstreamdoc.Name = fmt.Sprintf("tmp%v", k)
		upstreamdoc.ServerList = []ServerInfo{serverinfo}
		upstreamdoc.Forward = "chash"
		upstreamdoc.Version = 1
		upstreamdoc.Time = time.Now()
		upstreamdoc.FilePath = "123"
		res, err := collection.InsertOne(context.TODO(), upstreamdoc)
		if err != nil {
			panic(err)
		}
		fmt.Println(k, res)

	}
}

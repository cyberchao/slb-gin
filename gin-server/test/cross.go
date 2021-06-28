package test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aluttik/go-crossplane"
	"github.com/tidwall/gjson"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"slb-admin/global"
)

func Cross() {
	payload, err := crossplane.Parse("yqb.com.conf", &crossplane.ParseOptions{SingleFile: true, StopParsingOnError: true})
	if err != nil {
		panic(err)
	}
	b, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}
	collection := global.Mogo.Database("slb").Collection("vhost")

	server := gjson.Get(string(b), "config.0.parsed.0.block.0")
	fmt.Println(server.String())

	//insert 参考：https://groups.google.com/g/mgo-users/c/tLvYfXExkRk
	//var docs interface{}
	//json.Unmarshal([]byte(server.String()),&docs)
	//collection.InsertOne(context.TODO(), docs)

	//find
	filterCursor, err := collection.Find(context.TODO(), bson.M{"directive": "server"})
	var result []bson.M
	if err = filterCursor.All(context.TODO(), &result); err != nil {
		log.Fatal(err)
	}
	for _, m := range result {
		fmt.Println(m["directive"].(string))
	}
}

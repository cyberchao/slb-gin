package initialize

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"reflect"
	"slb-admin/global"
	"time"
)

func Mongo() *mongo.Client {
	mongoc := global.CONFIG.Mongo
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s", mongoc.Username, mongoc.Password, mongoc.Host, mongoc.Port, mongoc.Dbname)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// Connect to MongoDB
	if mogo, err := mongo.Connect(ctx, options.Client().ApplyURI(uri)); err != nil {
		panic(err)
	} else {
		fmt.Println(reflect.TypeOf(mogo))
		if err = mogo.Ping(context.TODO(), readpref.Primary()); err != nil {
			panic(err)
		}
		return mogo
	}
}



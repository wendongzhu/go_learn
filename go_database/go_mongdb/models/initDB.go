package models

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Database struct {
	Mongo *mongo.Client
}

var DB *Database

// Init 初始化
func Init() {
	DB = &Database{
		Mongo: SetConnect(),
	}
}

// SetConnect 连接设置
func SetConnect() *mongo.Client {
	//uri := "mongodb+srv://用户名:密码@官方给的.mongodb.net"
	uri := "mongodb://localhost:27017"
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri).SetMaxPoolSize(20)) // 连接池
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("connection sussful\n")
	return client
}

package models

import (
	"context"
	"github.com/go-redis/redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var Mongodb *mongo.Database
var Redis *redis.Client

func init() {
	Mongodb = getMongoDB()
	Redis = getRedis()
}

// 获取一个mongodb实例
func getMongoDB() *mongo.Database {
	options := options.Client().SetAuth(options.Credential{
		Username: "admin",
		Password: "123456",
	}).ApplyURI("mongodb://127.0.0.1:27017")

	client, err := mongo.Connect(context.Background(), options)
	if err != nil {
		log.Fatalln(err)
	}
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatalln(err)
	}
	return client.Database("im")
}

// 获取一个redis实例
func getRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Username: "default",
		Password: "123456",
		DB:       0,
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalln(err)
	}
	return client
}

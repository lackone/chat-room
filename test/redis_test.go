package test

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v9"
	"testing"
	"time"
)

func TestRedis(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Username: "default", //ACL默认用户名
		Password: "123456",
		DB:       0,
	})
	result, err := client.Ping(context.Background()).Result()
	if err != nil {
		t.Log(err)
	}
	t.Log(result)

	fmt.Println(client.Set(context.Background(), "test", "test", 30*time.Second))
}

package test

import (
	"context"
	"github.com/redis/go-redis/v9"
	"testing"
	"time"
)

var ctx = context.Background()
var rdb = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "", // no password set
	DB:       0,  // use default DB
})

func TestSetValue(t *testing.T) {
	err := rdb.Set(ctx, "key", "value", time.Second*10).Err()
	if err != nil {
		t.Error(err)
	}
}

func TestGetValue(t *testing.T) {
	result, err := rdb.Get(ctx, "key").Result()
	if err != nil {
		t.Error(err)
	}
	t.Log(result)
}

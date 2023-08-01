package models

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"
	"log"
	"xorm.io/xorm"
)

var Engine = Init()
var RDB = InitRedis()

func Init() *xorm.Engine {
	engine, err := xorm.NewEngine("mysql", "root:20000620@tcp(127.0.0.1:3306)/cloud-disk?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		log.Printf("Xorm new engine error!%v", err)
		return nil
	}
	return engine

}

func InitRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

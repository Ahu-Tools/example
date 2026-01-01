package redis

import (
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	//@ahum: imports
)

func GetClient(db int) *redis.Client {
	host := viper.GetString("infras.redis.host")
	port := viper.GetString("infras.redis.port")

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Username: viper.GetString("infras.redis.username"),
		Password: viper.GetString("infras.redis.password"),
		DB:       db,
	})

	return rdb
}

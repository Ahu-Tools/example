package asynq

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
)

var client *asynq.Client
var clientOnce sync.Once

func GetClient() *asynq.Client {
	clientOnce.Do(func() {
		client = asynq.NewClientFromRedisClient(getRedis())
	})

	return client
}

func NewTask(moduleName, taskName string, payload interface{}, opts ...asynq.Option) (*asynq.Task, error) {
	pload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	modCfg := ModuleConfig(moduleName)
	verName := modCfg["version"].(string)
	return asynq.NewTask(fmt.Sprintf("%s:%s:%s", verName, moduleName, taskName), pload, opts...), nil
}

func getRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", ConfigVar(HOST_CONFIG).(string), ConfigVar(PORT_CONFIG).(int)),
		Username: ConfigVar(USER_CONFIG).(string),
		Password: ConfigVar(PASSWORD_CONFIG).(string),
		DB:       ConfigVar(DB_CONFIG).(int),
	})

	return rdb
}

package asynq

import (
	"fmt"

	"github.com/spf13/viper"
)

type ConfigName uint

const (
	HOST_CONFIG ConfigName = iota
	PORT_CONFIG
	USER_CONFIG
	PASSWORD_CONFIG
	DB_CONFIG
)

var configVars = map[ConfigName]string{
	HOST_CONFIG:     "infras.asynq.redis.host",
	PORT_CONFIG:     "infras.asynq.redis.port",
	USER_CONFIG:     "infras.asynq.redis.username",
	PASSWORD_CONFIG: "infras.asynq.redis.password",
	DB_CONFIG:       "infras.asynq.redis.db",
}

var configMap = make(map[ConfigName]any)

func Configure() error {
	for cfgName, cfgPath := range configVars {
		if !viper.IsSet(cfgPath) {
			return fmt.Errorf("'%s' not set", cfgPath)
		}
		configMap[cfgName] = viper.Get(cfgPath)
	}

	return nil
}

func ConfigVar(name ConfigName) any {
	return configVars[name]
}

func ModuleConfig(moduleName string) map[string]any {
	cfgPath := fmt.Sprintf("infras.asynq.%s", moduleName)
	return viper.GetStringMap(cfgPath)
}

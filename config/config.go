package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"

	"github.com/Ahu-Tools/example/infrastructure/postgres"

	"github.com/Ahu-Tools/example/infrastructure/asynq"
	//@ahum: imports
)

// NewConfig loads configuration from environment variables or .env file.
func CheckConfigs() {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath("./config")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			panic("No config file found!")
		} else {
			panic(fmt.Errorf("Error happened during loading config file: %e", err))
		}
	}

	secretKey := viper.Get("app.secret_key")
	if _, ok := secretKey.(string); !ok {
		log.Fatal("app.secret_key not set. Please provide it.")
	}

}

func ConfigInfras() error {

	// @ahum:edges.group

	// @ahum:gin.load

	// @ahum:end.gin.load

	// @ahum:connect.load

	// @ahum:end.connect.load

	// @ahum:asynq.load

	// @ahum:end.asynq.load

	// @ahum:end.edges.group

	// @ahum:infras.group

	// @ahum:postgres.load
	err := postgres.Configure()
	if err != nil {
		return err
	}
	// @ahum:end.postgres.load

	// @ahum:redis.load

	// @ahum:end.redis.load

	// @ahum:asynq.load
	err = asynq.Configure()
	if err != nil {
		return err
	}
	// @ahum:end.asynq.load

	// @ahum:end.infras.group

	//@ahum: loads

	return nil
}

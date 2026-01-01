package postgres

import (
	"log"

	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var Db *gorm.DB

func Configure() error {
	ps := NewPostgresConfig(
		viper.GetString("infras.postgres.user"),
		viper.GetString("infras.postgres.password"),
		viper.GetString("infras.postgres.db_name"),
		viper.GetString("infras.postgres.host"),
		viper.GetString("infras.postgres.sslmode"),
		viper.GetInt("infras.postgres.port"),
	)

	viper.Set("infras.postgres.url", ps.Url())
	postgresURL := viper.Get("infras.postgres.url")
	if _, ok := postgresURL.(string); !ok {
		log.Fatal("infras.postgres.url not set. Please provide it.")
	}

	var err error
	Db, err = NewConnection(viper.GetString("infras.postgres.url"))

	return err
}

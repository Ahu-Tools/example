package postgres

import (
	"log"

	"github.com/Ahu-Tools/example/infrastructure/postgres/security"
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

	// Register the callback to run BEFORE every Create and Update
	// 1. Run before the INSERT statement is generated
	Db.Callback().Create().Before("gorm:create").Register("blind_index", security.BlindIndexCallback)

	// 2. Run before the UPDATE statement is generated
	Db.Callback().Update().Before("gorm:update").Register("blind_index", security.BlindIndexCallback)
	return err
}

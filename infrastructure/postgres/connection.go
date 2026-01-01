package postgres

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresConfig struct {
	user     string
	password string
	db       string
	host     string
	port     int
	sslMode  string
}

func NewPostgresConfig(user, password, db, host, sslMode string, port int) *PostgresConfig {
	return &PostgresConfig{
		user,
		password,
		db,
		host,
		port,
		sslMode,
	}
}

func NewConnection(dsn string) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func (c PostgresConfig) Url() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=%s",
		c.user,
		c.password,
		c.host,
		c.port,
		c.db,
		c.sslMode,
	)
}

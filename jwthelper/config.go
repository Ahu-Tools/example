package jwthelper

import (
	"errors"
	"log/slog"

	"github.com/spf13/viper"
)

var (
	ErrPublicKeyNotSet  = errors.New("jwt.public_key_file not set")
	ErrPrivateKeyNotSet = errors.New("jwt.private_key_file not set")
)

var logger *slog.Logger

func Configure(lger *slog.Logger) error {
	logger = lger
	jwtPrivateKeyFile := viper.Get("jwt.private_key_file")
	if _, ok := jwtPrivateKeyFile.(string); !ok {
		logger.Error("config.CheckConfigs: jwt.private_key_file not set. Please provide it.", "code_level", "config", "method", "CheckConfigs")
		return ErrPrivateKeyNotSet
	}
	if _, err := GetPrivateKey(); err != nil {
		logger.Error("config.CheckConfigs: jwt.private_key_file is invalid", "code_level", "config", "method", "CheckConfigs", "error", err)
		return ErrPrivateKeyNotSet
	}

	jwtPublicKeyFile := viper.Get("jwt.public_key_file")
	if _, ok := jwtPublicKeyFile.(string); !ok {
		logger.Error("config.CheckConfigs: jwt.public_key_file not set. Please provide it.", "code_level", "config", "method", "CheckConfigs")
		return ErrPublicKeyNotSet
	}
	if _, err := GetPublicKey(); err != nil {
		logger.Error("config.CheckConfigs: jwt.public_key_file is invalid", "code_level", "config", "method", "CheckConfigs", "error", err)
		return ErrPublicKeyNotSet
	}

	return nil
}

package crypto

import (
	"errors"

	"github.com/Ahu-Tools/example/crypto/mock"
	"github.com/spf13/viper"
)

var (
	ErrMissingBlindPepper = errors.New("missing blind pepper")
)

var GlobalEncrypter Encrypter

func Configure() error {
	if viper.GetString("app.env") == "dev" {
		GlobalEncrypter = mock.NewRotationManager()
	} else {
		panic("unimplemented main encrypter")
	}

	blindPepper := viper.GetString("app.secret_key")
	if len(blindPepper) == 0 {
		return ErrMissingBlindPepper
	}

	return nil
}

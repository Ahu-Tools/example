package jwthelper

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

var (
	ErrJWTGenFail     = errors.New("unable to generate JWT")
	ErrPrivateKeyFail = errors.New("unable to get private key")
	ErrPublicKeyFail  = errors.New("unable to get public key")
	ErrTokenInvalid   = errors.New("invalid token")
)

func ParseWithClaims(token string, claimsType jwt.Claims) (*jwt.Token, error) {
	//# Check and get key
	key, err := ParsePubKey(token)
	if err != nil {
		logger.Error("jwt.ParseRSAPublicKeyFromPEM() failed", "code_level", "jwthelper.ParseTParseWithClaimsoken()", "error", err)
		return nil, ErrPublicKeyFail
	}

	//# Parse token
	parsedToken, err := jwt.ParseWithClaims(token, claimsType, func(t *jwt.Token) (interface{}, error) {
		return key, nil
	})

	if err != nil {
		logger.Error("jwt.Parse failed", "code_level", "jwthelper.ParseWithClaims()", "error", err)
		return nil, err
	}
	return parsedToken, nil
}

func ParseToken(token string) (*jwt.Token, error) {
	//# Check and get key
	key, err := ParsePubKey(token)
	if err != nil {
		logger.Error("jwt.ParseRSAPublicKeyFromPEM() failed", "code_level", "jwthelper.ParseToken()", "error", err)
		return nil, ErrPublicKeyFail
	}

	//# Parse token
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return key, nil
	})

	if err != nil {
		logger.Error("jwt.Parse failed", "code_level", "jwthelper.ParseToken()", "error", err)
		return nil, err
	}

	return parsedToken, nil
}

func ParsePubKey(token string) (*rsa.PublicKey, error) {
	//# Get public key
	publicKey, err := GetPublicKey()
	if err != nil {
		logger.Error("GetPublicKey() failed", "code_level", "jwthelper.GetPubKey()", "error", err)
		return nil, ErrPublicKeyFail
	}

	key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(publicKey))
	if err != nil {
		logger.Error("jwt.ParseRSAPublicKeyFromPEM() failed", "code_level", "jwthelper.GetPubKey()", "error", err)
		return nil, ErrPublicKeyFail
	}
	return key, nil
}

func GenerateToken(claims jwt.MapClaims) (string, error) {
	privateKey, err := GetPrivateKey()
	if err != nil {
		logger.Error("unable to get private key: GetPrivateKey()", "code_level", "jwthelper.Infra.GenerateToken", "error", err)
		return "", ErrPrivateKeyFail
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKey))
	if err != nil {
		logger.Error("unable to parse private key: ParseRSAPrivateKeyFromPEM()", "code_level", "jwthelper.Infra.GenerateToken", "error", err)
		return "", ErrPrivateKeyFail
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signedToken, err := token.SignedString(key)
	if err != nil {
		logger.Error("unable to sign token: token.SignedString(key)", "code_level", "jwthelper.Infra.GenerateToken", "error", err)
		return "", ErrJWTGenFail
	}

	return signedToken, nil
}

func GetPublicKey() (string, error) {
	publicKeyPath := viper.GetString("jwt.public_key_file")
	publicKeyBytes, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return "", fmt.Errorf("could not read public key file: %w", err)
	}
	return string(publicKeyBytes), nil
}

func GetPrivateKey() (string, error) {
	privateKeyPath := viper.GetString("jwt.private_key_file")
	privateKeyBytes, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return "", fmt.Errorf("could not read private key file: %w", err)
	}
	return string(privateKeyBytes), nil
}

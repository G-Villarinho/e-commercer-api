package config

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"log/slog"
	"os"

	"github.com/Netflix/go-env"
	"github.com/joho/godotenv"
)

var Env Environment

type Environment struct {
	ConnectionString string `env:"CONNECTION_STRING"`
	RedisAdress      string `env:"REDIS_ADRESS"`
	RedisPassword    string `env:"REDIS_PASSWORD"`
	RedisDB          int    `env:"REDIS_DB"`
	MongoURI         string `env:"MONGO_URI"`
	APIPort          string `env:"API_PORT"`
	TokenExp         int    `env:"TOKEN_EXP"`
	ResendKey        string `env:"RESEND_KEY"`
	URLFront         string `env:"FRONT_URL"`
	OTPEmailSize     int8   `env:"OTP_EMAIL_SIZE"`
	OTPExp           uint8  `env:"OTP_EXP"`
	PrivateKey       *ecdsa.PrivateKey
	PublicKey        *ecdsa.PublicKey
}

func LoadEnvironments() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	_, err = env.UnmarshalFromEnviron(&Env)
	if err != nil {
		panic(err)
	}

	Env.PrivateKey, err = loadPrivateKey()
	if err != nil {
		panic(err)
	}

	Env.PublicKey, err = loadPublicKey()
	if err != nil {
		panic(err)
	}
}

func ConfigureLogger() {
	handler := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: false,
	}))
	slog.SetDefault(handler)
}

func loadPrivateKey() (*ecdsa.PrivateKey, error) {
	keyData, err := os.ReadFile("ec_private_key.pem")
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(keyData)
	if block == nil || block.Type != "EC PRIVATE KEY" {
		return nil, errors.New("failed to decode PEM block containing private key")
	}

	privateKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func loadPublicKey() (*ecdsa.PublicKey, error) {
	keyData, err := os.ReadFile("ec_public_key.pem")
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(keyData)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, errors.New("failed to decode PEM block containing public key")
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	ecdsaPubKey, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("not ECDSA public key")
	}

	return ecdsaPubKey, nil
}

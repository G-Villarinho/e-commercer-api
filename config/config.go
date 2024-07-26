package config

import (
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
	SecretKeyPath    string `env:"SECRET_KEY_PATH"`
	OTPExp           uint8  `env:"OTP_EXP"`
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
}

func ConfigureLogger() {
	handler := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: false,
	}))
	slog.SetDefault(handler)
}

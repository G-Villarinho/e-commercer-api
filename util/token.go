package util

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
	"strings"

	"github.com/OVillas/e-commercer-api/config"
	"github.com/OVillas/e-commercer-api/domain"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func CreateToken(user domain.User) (string, error) {
	privateKey, err := LoadPrivateKey(config.Env.SecretKeyPath)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"id":        user.ID,
		"name":      user.Name,
		"email":     user.Email,
		"avatarURL": user.AvatarURL,
	})

	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ExtractToken(ctx echo.Context) (string, error) {
	token := ctx.Request().Header.Get("Authorization")

	length := len(strings.Split(token, " "))
	if length == 2 {
		return strings.Split(token, " ")[1], nil
	}

	return "", domain.ErrSessionNotFound
}

func ExtractUserIDFromToken(tokenString string) (string, error) {
	publicKey, err := LoadPublicKey(config.Env.SecretKeyPath)
	if err != nil {
		return "", err
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, domain.ErrorUnexpectedMethod
		}
		return publicKey, nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", domain.ErrTokenInvalid
	}

	userID, ok := claims["id"].(string)
	if !ok {
		return "", domain.ErrTokenInvalid
	}

	return userID, nil
}

func LoadPrivateKey(path string) (*ecdsa.PrivateKey, error) {
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

func LoadPublicKey(path string) (*ecdsa.PublicKey, error) {
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

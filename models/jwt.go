package models

import (
	"encoding/base64"
	"fmt"
	"time"

	"mikiwa/utils"

	"github.com/beego/beego/v2/core/logs"
	"github.com/dgrijalva/jwt-go"
)

func CreateToken(ttl time.Duration, payload interface{}, privateKey string) (string, string, error) {
	logs.Info("Create Token")

	decodedPrivateKey, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		return "", "", fmt.Errorf("could not decode key: %w", err)
	}
	key, err := jwt.ParseRSAPrivateKeyFromPEM(decodedPrivateKey)

	if err != nil {
		return "", "", fmt.Errorf("create: parse key: %w", err)
	}

	now := utils.GetSvrDate()

	logs.Info("Server date :", now)

	claims := make(jwt.MapClaims)
	claims["sub"] = payload
	claims["exp"] = now.Add(ttl).Unix()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()
	claims["expired"] = time.Unix(now.Add(ttl).Unix(), 0)

	logs.Info("Token will be expired at ", time.Unix(now.Add(ttl).Unix(), 0))

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(key)

	if err != nil {
		return "", "", fmt.Errorf("create: sign token: %w", err)
	}
	expdate := time.Unix(now.Add(ttl).Unix(), 0)
	return token, expdate.Format("2006-01-02 15:04:05"), nil
}

func VerifyToken(token string, publicKey string) (interface{}, error) {
	logs.Info("Verify Token")

	decodedPublicKey, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return nil, fmt.Errorf("could not decode: %w", err)
	}

	key, err := jwt.ParseRSAPublicKeyFromPEM(decodedPublicKey)

	if err != nil {
		return "", fmt.Errorf("validate: parse key: %w", err)
	}

	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected method: %s", t.Header["alg"])
		}
		return key, nil
	})

	if err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return nil, fmt.Errorf("validate: invalid token")
	}

	logs.Info("End Verify Token")

	return claims["sub"], nil
}

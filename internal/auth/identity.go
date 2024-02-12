package auth

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"log"
	"os"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateUserToken(id int) string {

	privateKeyRaw, _ := os.ReadFile(".ssh/key")

	block, _ := pem.Decode(privateKeyRaw)
	if block == nil || block.Type != "PRIVATE KEY" {
		log.Fatal("failed to decode PEM block containing private key")
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		log.Fatal(err)
	}

	token := jwt.NewWithClaims(&jwt.SigningMethodEd25519{}, jwt.MapClaims{
		"iss": "commute-and-mute",
		"sub": strconv.FormatInt(int64(id), 10),
	})

	tokenString, _ := token.SignedString(key)

	return tokenString
}

func GetConnectedUserId(tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return GetPublicKey(), nil
	})
	if err != nil {
		return -1, err
	}
	if !token.Valid {
		return -1, errors.New("Token invalid")
	}
	strId, _ := token.Claims.GetSubject()
	return strconv.Atoi(strId)
}

func GetPublicKey() ed25519.PublicKey {
	publicKeyRaw, _ := os.ReadFile(".ssh/key.pub")

	block, _ := pem.Decode(publicKeyRaw)
	if block == nil || block.Type != "PUBLIC KEY" {
		log.Fatal("failed to decode PEM block containing private key")
	}

	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		log.Fatal(err)
	}
	return key.(ed25519.PublicKey)
}

package auth

import (
	"context"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/asn1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/golang-jwt/jwt/v5"
)

type JwtHeader struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

type JwtPayload struct {
	Sub int    `json:"sub"`
	Iss string `json:"iss"`
	Exp int    `json:"exp"`
}

type TokenGenerator struct {
	client KmsApi
}

type KmsApi interface {
	Sign(ctx context.Context, params *kms.SignInput, optFns ...func(*kms.Options)) (*kms.SignOutput, error)
	GetPublicKey(ctx context.Context, params *kms.GetPublicKeyInput, optFns ...func(*kms.Options)) (*kms.GetPublicKeyOutput, error)
}

func NewTokenGenerator(config aws.Config) TokenGenerator {
	return TokenGenerator{client: kms.NewFromConfig(config)}
}

func (tokenGenerator TokenGenerator) GenerateForId(ctx context.Context, id int) string {
	now := time.Now()

	header := JwtHeader{Alg: "ES256", Typ: "jwt"}
	headerBytes, _ := json.Marshal(header)
	headerBase64 := base64.RawURLEncoding.EncodeToString(headerBytes)

	payload := JwtPayload{Sub: id, Iss: "commute-and-mute", Exp: int(now.Add(time.Hour).Unix())}
	payloadBytes, _ := json.Marshal(payload)
	payloadBase64 := base64.RawURLEncoding.EncodeToString(payloadBytes)

	unsignedString := fmt.Sprintf("%v.%v", headerBase64, payloadBase64)
	keyId := "ddd8b1bf-b47b-4a62-8ea2-108df35a9f12"
	signInput := kms.SignInput{
		KeyId:            &keyId,
		SigningAlgorithm: types.SigningAlgorithmSpecEcdsaSha256,
		Message:          []byte(unsignedString),
		MessageType:      types.MessageTypeRaw,
	}
	signOutput, err := tokenGenerator.client.Sign(ctx, &signInput)

	if err != nil {
		log.Panicf("error signing new token: %v", err)
	}

	signatureB64, err := kmsResponseToJwtSignature(signOutput.Signature)
	if err != nil {
		log.Panic("could not generate jwt signature", err)
	}

	signedToken := fmt.Sprintf("%v.%v", unsignedString, signatureB64)
	return signedToken
}

func (tokenGenerator TokenGenerator) Validate(ctx context.Context, tokenString string) (int, error) {

	tokenJwt, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return tokenGenerator.getPublicKey(ctx)
	})

	if err != nil {
		return -1, fmt.Errorf("could not validate token: %w", err)
	}

	subjectString, err := tokenJwt.Claims.GetSubject()
	if err != nil {
		return -1, err
	}
	subjectInt, err := strconv.Atoi(subjectString)
	if err != nil {
		return -1, fmt.Errorf("jwt contains non-numeric subject")
	}
	return subjectInt, nil
}

func (tokenGenerator TokenGenerator) getPublicKey(ctx context.Context) (*ecdsa.PublicKey, error) {
	keyId := "ddd8b1bf-b47b-4a62-8ea2-108df35a9f12"
	input := kms.GetPublicKeyInput{
		KeyId:       &keyId,
		GrantTokens: []string{},
	}
	output, err := tokenGenerator.client.GetPublicKey(ctx, &input)
	if err != nil {
		return nil, fmt.Errorf("error validating token. could not get public key from kms: %w", err)
	}
	keyBytes := output.PublicKey
	key, err := x509.ParsePKIXPublicKey(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("could not parse public key retrieved from kms: %w", err)
	}
	return key.(*ecdsa.PublicKey), nil
}

func kmsResponseToJwtSignature(sigBytes []byte) (string, error) {
	var signatureStruct EcdsaSigValue
	_, err := asn1.Unmarshal(sigBytes, &signatureStruct)
	if err != nil {
		return "", err
	}
	rBytes := signatureStruct.R.Bytes()
	rBytesPadded := make([]byte, 32)
	copy(rBytesPadded[32-len(rBytes):], rBytes)

	sBytes := signatureStruct.S.Bytes()
	sBytesPadded := make([]byte, 32)
	copy(sBytesPadded[32-len(sBytes):], sBytes)

	out := append(rBytesPadded, sBytesPadded...)
	outString := base64.RawURLEncoding.EncodeToString(out)

	return outString, nil
}

type EcdsaSigValue struct {
	R *big.Int
	S *big.Int
}

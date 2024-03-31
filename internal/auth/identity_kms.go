package auth

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
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

type KmsUtils struct {
	client KmsApi
}

type KmsApi interface {
	Sign(ctx context.Context, params *kms.SignInput, optFns ...func(*kms.Options)) (*kms.SignOutput, error)
}

func NewKmsUtils(config aws.Config) KmsUtils {
	return KmsUtils{client: kms.NewFromConfig(config)}
}

func (kmsUtils KmsUtils) GenerateForId(ctx context.Context, id int) string {
	now := time.Now()

	header := JwtHeader{Alg: "ES256", Typ: "jwt"}
	headerBytes, _ := json.Marshal(header)
	headerBase64 := base64.RawStdEncoding.EncodeToString(headerBytes)

	payload := JwtPayload{Sub: id, Iss: "commute-and-mute", Exp: int(now.Add(time.Hour).Unix())}
	payloadBytes, _ := json.Marshal(payload)
	payloadBase64 := base64.RawStdEncoding.EncodeToString(payloadBytes)

	unsignedString := fmt.Sprintf("%v.%v", headerBase64, payloadBase64)
	keyId := "ddd8b1bf-b47b-4a62-8ea2-108df35a9f12"
	signInput := kms.SignInput{
		KeyId:            &keyId,
		SigningAlgorithm: types.SigningAlgorithmSpecEcdsaSha256,
		Message:          []byte(unsignedString),
	}
	signOutput, err := kmsUtils.client.Sign(ctx, &signInput)

	if err != nil {
		log.Panicf("error signing new token: %v", err)
	}

	signedToken := fmt.Sprintf("%v.%v", unsignedString, string(signOutput.Signature))
	return signedToken
}

package auth

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"
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

type KmsUtils struct{}

func (kmsUtils KmsUtils) GenerateForId(id int) string {
	// config, _ := config.LoadDefaultConfig(context.TODO(), config.WithRegion("eu-north-1"))
	// client := kms.NewFromConfig(config)

	now := time.Now()

	header := JwtHeader{Alg: "ES256", Typ: "jwt"}
	headerBytes, _ := json.Marshal(header)
	headerBase64 := base64.RawStdEncoding.EncodeToString(headerBytes)

	payload := JwtPayload{Sub: id, Iss: "commute-and-mute", Exp: int(now.Add(time.Hour).Unix())}
	payloadBytes, _ := json.Marshal(payload)
	payloadBase64 := base64.RawStdEncoding.EncodeToString(payloadBytes)

	unsignedString := fmt.Sprintf("%v.%v", headerBase64, payloadBase64)

	return unsignedString

}

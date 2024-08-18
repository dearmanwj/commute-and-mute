package auth

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/asn1"
	"math/big"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/golang-jwt/jwt/v5"
)

type mockKmsApi struct {
	mockSign         func(ctx context.Context, params *kms.SignInput, optFns ...func(*kms.Options)) (*kms.SignOutput, error)
	mockGetPublicKey func(ctx context.Context, params *kms.GetPublicKeyInput, optFns ...func(*kms.Options)) (*kms.GetPublicKeyOutput, error)
}

func (m mockKmsApi) Sign(ctx context.Context, params *kms.SignInput, optFns ...func(*kms.Options)) (*kms.SignOutput, error) {
	return m.mockSign(ctx, params, optFns...)
}

func (m mockKmsApi) GetPublicKey(ctx context.Context, params *kms.GetPublicKeyInput, optFns ...func(*kms.Options)) (*kms.GetPublicKeyOutput, error) {
	return m.mockGetPublicKey(ctx, params, optFns...)
}

func TestGenerateSigned(t *testing.T) {
	// Given
	id := 123
	mockSignature := getMockSignature()
	mock := mockKmsApi{
		mockSign: func(ctx context.Context, params *kms.SignInput, optFns ...func(*kms.Options)) (*kms.SignOutput, error) {
			if params.SigningAlgorithm != types.SigningAlgorithmSpecEcdsaSha256 {
				t.Error("incorrect signing algorithm")
			}
			output := kms.SignOutput{
				Signature: mockSignature,
			}
			return &output, nil
		},
		mockGetPublicKey: nil,
	}
	generator := TokenGenerator{client: mock}
	ctx := context.Background()

	// When
	token := generator.GenerateForId(ctx, id)

	// Then
	tokenParts := strings.Split(token, ".")
	// Can definitely think of a better assert than this...
	if len(tokenParts) != 3 {
		t.Error("incorrect signature")
	}
}

func TestValidateToken(t *testing.T) {
	// Given
	id := 1234567
	privateKey := generateKeyPair()
	mock := mockKmsApi{
		mockSign: nil,
		mockGetPublicKey: func(ctx context.Context, params *kms.GetPublicKeyInput, optFns ...func(*kms.Options)) (*kms.GetPublicKeyOutput, error) {
			publicKey := &privateKey.PublicKey
			keyBytes, _ := x509.MarshalPKIXPublicKey(publicKey)
			output := kms.GetPublicKeyOutput{
				PublicKey: keyBytes,
			}
			return &output, nil
		},
	}

	generator := TokenGenerator{
		client: mock,
	}
	ctx := context.Background()
	expTime := time.Now().Add(time.Hour)
	token := jwt.NewWithClaims(jwt.SigningMethodES256,
		jwt.RegisteredClaims{
			Subject:   strconv.Itoa(id),
			ExpiresAt: jwt.NewNumericDate(expTime),
		})
	tokenString, _ := token.SignedString(privateKey)

	// When
	result, err := generator.GetIdIfValid(ctx, tokenString)

	// Then
	if err != nil {
		t.Error("could not validate token", err)
	}
	if result != strconv.Itoa(id) {
		t.Error("incorrect sub")
	}
}

func getMockSignature() []byte {
	sig := EcdsaSigValue{
		R: big.NewInt(1234567),
		S: big.NewInt(7654321),
	}
	encoded, _ := asn1.Marshal(sig)
	return encoded
}

func generateKeyPair() *ecdsa.PrivateKey {
	curve := elliptic.P256()
	key, _ := ecdsa.GenerateKey(curve, rand.Reader)
	return key
}

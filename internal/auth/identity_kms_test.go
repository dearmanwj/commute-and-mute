package auth

import (
	"context"
	"encoding/asn1"
	"encoding/base64"
	"math/big"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
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
	mock := mockKmsApi{
		mockSign: nil,
		mockGetPublicKey: func(ctx context.Context, params *kms.GetPublicKeyInput, optFns ...func(*kms.Options)) (*kms.GetPublicKeyOutput, error) {
			keyString := "MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAETz8/nVxGqJ0QwrDsnuG6EOEEUZ3jQ1rbwJo7G8IJ2zGb8p2Xjrph+90p6T5ityqoYW+inVJ2vh+kmbdb9jzcBA=="
			keyBytes, _ := base64.StdEncoding.DecodeString(keyString)
			output := kms.GetPublicKeyOutput{
				PublicKey: keyBytes,
			}
			return &output, nil
		},
	}

	token := "eyJhbGciOiJFUzI1NiIsInR5cCI6Imp3dCJ9.eyJzdWIiOiIxMDUwMzgxMiIsImlzcyI6ImNvbW11dGUtYW5kLW11dGUiLCJleHAiOjE3MTI2MDM3NDJ9.K_UlrCKa_SKrWVDTAfvOtZPGIct2BHCHHVO_T4OihQTvA6R2-_-3qsB4bMdrPghkrl4kM4p-0Lgc6vZ29Btrtg"

	generator := TokenGenerator{
		client: mock,
	}
	ctx := context.Background()

	// When
	result, err := generator.GetIdIfValid(ctx, token)

	// Then
	if err != nil {
		t.Error("could not validate token")
	}
	if result != "10503812" {
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

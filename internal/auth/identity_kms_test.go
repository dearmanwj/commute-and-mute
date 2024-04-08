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

	token := "eyJhbGciOiJFUzI1NiIsInR5cCI6Imp3dCJ9.eyJzdWIiOjEwNTAzODEyLCJpc3MiOiJjb21tdXRlLWFuZC1tdXRlIiwiZXhwIjoxNzEyNTk5ODI5fQ.p5srwusF05SV_F7Z7p-NTlKJCafL7lYBHNmJ0DeNwyt_0aXXlV6COdI2O6eCs06ZmldyaCYPhc5Uc47jXNtlhA"

	generator := TokenGenerator{
		client: mock,
	}
	ctx := context.Background()

	// When
	result, err := generator.Validate(ctx, token)

	// Then
	if err != nil {
		t.Error("could not validate token")
	}
	if result != 10503812 {
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

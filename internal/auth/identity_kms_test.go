package auth

import (
	"context"
	"encoding/base64"
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
	mockSignature := []byte("signature")
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
	expectedSignature := base64.RawURLEncoding.EncodeToString(mockSignature)
	tokenParts := strings.Split(token, ".")
	if tokenParts[2] != expectedSignature {
		t.Error("incorrect signature")
	}
}

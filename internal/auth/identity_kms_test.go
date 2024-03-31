package auth

import (
	"context"
	"encoding/base64"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
)

type mockKmsApi func(ctx context.Context, params *kms.SignInput, optFns ...func(*kms.Options)) (*kms.SignOutput, error)

func (m mockKmsApi) Sign(ctx context.Context, params *kms.SignInput, optFns ...func(*kms.Options)) (*kms.SignOutput, error) {
	return m(ctx, params, optFns...)
}

func TestGenerateSigned(t *testing.T) {
	// Given
	id := 123
	mockSignature := []byte("signature")
	mockKms := mockKmsApi(func(ctx context.Context, params *kms.SignInput, optFns ...func(*kms.Options)) (*kms.SignOutput, error) {
		if params.SigningAlgorithm != types.SigningAlgorithmSpecEcdsaSha256 {
			t.Error("incorrect signing algorithm")
		}
		output := kms.SignOutput{
			Signature: mockSignature,
		}
		return &output, nil
	})
	generator := TokenGenerator{client: mockKms}
	ctx := context.Background()

	// When
	token := generator.GenerateForId(ctx, id)

	// Then
	expectedSignature := base64.RawStdEncoding.EncodeToString(mockSignature)
	tokenParts := strings.Split(token, ".")
	if tokenParts[2] != expectedSignature {
		t.Error("incorrect signature")
	}
}

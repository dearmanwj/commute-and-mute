package auth

import (
	"context"
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
	mockKms := mockKmsApi(func(ctx context.Context, params *kms.SignInput, optFns ...func(*kms.Options)) (*kms.SignOutput, error) {
		if params.SigningAlgorithm != types.SigningAlgorithmSpecEcdsaSha256 {
			t.Error("incorrect signing algorithm")
		}
		output := kms.SignOutput{
			Signature: []byte("signature"),
		}
		return &output, nil
	})
	utils := KmsUtils{client: mockKms}
	ctx := context.Background()

	// When
	token := utils.GenerateForId(ctx, id)

	// Then
	tokenParts := strings.Split(token, ".")
	if tokenParts[2] != "signature" {
		t.Error("incorrect signature")
	}
}

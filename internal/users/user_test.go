package users

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type MockDynamoClient struct {
	get    func(context.Context, *dynamodb.GetItemInput, ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	put    func(context.Context, *dynamodb.PutItemInput, ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	delete func(context.Context, *dynamodb.DeleteItemInput, ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error)
}

func (m MockDynamoClient) GetItem(context context.Context, input *dynamodb.GetItemInput, opts ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
	return m.get(context, input, opts...)
}

func (m MockDynamoClient) PutItem(context context.Context, input *dynamodb.PutItemInput, opts ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	return m.put(context, input, opts...)
}

func (m MockDynamoClient) DeleteItem(context context.Context, input *dynamodb.DeleteItemInput, opts ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error) {
	return m.delete(context, input, opts...)
}

func TestGetUser(t *testing.T) {

	mockDbClient := MockDynamoClient{
		get: func(ctx context.Context, gii *dynamodb.GetItemInput, f ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
			return new(dynamodb.GetItemOutput), nil
		},
		put:    nil,
		delete: nil,
	}

	db := TableBasics{
		DynamoDbClient: mockDbClient,
		TableName:      "user",
	}

	user, err := db.GetUser(context.Background(), 123)

	expected := User{
		HomeLat: -1,
		HomeLng: -1,
		WorkLat: -1,
		WorkLng: -1,
	}

	if err != nil {
		t.Error(err)
	}

	if user != expected {
		t.Error("empty user not expected")
	}

}

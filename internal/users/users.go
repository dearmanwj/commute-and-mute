package users

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type User struct {
	ID           int
	AccessToken  string
	RefreshToken string
	HomeLat      float64
	HomeLng      float64
	WorkLat      float64
	WorkLng      float64
	ExpiresAt    int64
}

type TableBasics struct {
	DynamoDbClient DynamoClient
	TableName      string
}

type DynamoClient interface {
	GetItem(context.Context, *dynamodb.GetItemInput, ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	PutItem(context.Context, *dynamodb.PutItemInput, ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	DeleteItem(context.Context, *dynamodb.DeleteItemInput, ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error)
}

func (t TableBasics) UpdateUser(ctx context.Context, user User) error {

	item, err := attributevalue.MarshalMap(user)

	if err != nil {
		panic(err)
	}

	_, err = t.DynamoDbClient.PutItem(ctx, &dynamodb.PutItemInput{
		Item: item, TableName: aws.String(t.TableName),
	},
	)

	if err != nil {
		log.Printf("Could not update user record: %v\n", err)
	}
	return err
}

func (t TableBasics) GetUser(ctx context.Context, id int) (User, error) {
	val, err := attributevalue.Marshal(id)
	if err != nil {
		panic(err)
	}
	key := map[string]types.AttributeValue{"ID": val}

	var user User

	response, err := t.DynamoDbClient.GetItem(ctx, &dynamodb.GetItemInput{Key: key, TableName: aws.String(t.TableName)})
	if err != nil {
		log.Printf("Error getting user: %v,\n", err)
		return User{}, err
	}

	if response.Item == nil {
		return User{HomeLat: -1, HomeLng: -1, WorkLat: -1, WorkLng: -1}, nil
	}

	err = attributevalue.UnmarshalMap(response.Item, &user)

	if err != nil {
		log.Printf("Error mapping user to go type: %v,\n", user)
		return User{}, err
	}

	return user, nil
}

func (t TableBasics) DeleteUser(context context.Context, id int) error {
	val, err := attributevalue.Marshal(id)
	if err != nil {
		panic(err)
	}
	key := map[string]types.AttributeValue{"ID": val}
	_, err = t.DynamoDbClient.DeleteItem(context, &dynamodb.DeleteItemInput{Key: key, TableName: aws.String(t.TableName)})
	if err != nil {
		return fmt.Errorf("error deleting user from db %v", err)
	}
	return nil
}

func GetDbConnection(ctx context.Context) (TableBasics, error) {
	tableName := os.Getenv("USERS_TABLE_NAME")
	config, err := config.LoadDefaultConfig(ctx, config.WithRegion("eu-north-1"))
	if err != nil {
		log.Printf("Error getting db connection: %v\n", err)
		return TableBasics{}, err
	}

	return TableBasics{TableName: tableName, DynamoDbClient: dynamodb.NewFromConfig(config)}, nil
}

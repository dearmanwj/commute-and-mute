package users

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var DB TableBasics

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
	DynamoDbClient *dynamodb.Client
	TableName      string
}

func UpdateUser(ctx context.Context, user User) error {

	item, err := attributevalue.MarshalMap(user)

	if err != nil {
		panic(err)
	}

	_, err = DB.DynamoDbClient.PutItem(ctx, &dynamodb.PutItemInput{
		Item: item, TableName: aws.String(DB.TableName),
	},
	)

	if err != nil {
		log.Printf("Could not update user record: %v\n", err)
	}
	return err
}

func GetUser(ctx context.Context, id int) (User, error) {
	val, err := attributevalue.Marshal(id)
	if err != nil {
		panic(err)
	}
	key := map[string]types.AttributeValue{"ID": val}

	var user User

	response, err := DB.DynamoDbClient.GetItem(ctx, &dynamodb.GetItemInput{Key: key, TableName: aws.String(DB.TableName)})
	if err != nil {
		log.Printf("Error getting user: %v,\n", err)
		return user, err
	}

	err = attributevalue.UnmarshalMap(response.Item, &user)

	if err != nil {
		log.Printf("Error mapping user to go type: %v,\n", user)
		return user, err
	}

	return user, nil
}

func GetDbConnection(ctx context.Context) error {
	tableName := os.Getenv("USERS_TABLE_NAME")
	config, err := config.LoadDefaultConfig(ctx, config.WithRegion("eu-north-1"))
	if err != nil {
		log.Printf("Error getting db connection: %v\n", err)
		return err
	}

	DB = TableBasics{TableName: tableName, DynamoDbClient: dynamodb.NewFromConfig(config)}
	return nil
}

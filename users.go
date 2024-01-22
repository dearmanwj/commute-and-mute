package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type User struct {
	UserName     string
	AccessToken  string
	RefreshToken string
	HomeLat      float64
	HomeLng      float64
	WorkLat      float64
	WorkLng      float64
}

type TableBasics struct {
	DynamoDbClient *dynamodb.Client
	TableName      string
}

func UpdateUser(user User) error {
	tableName := os.Getenv("USERS_TABLE_NAME")
	config, _ := config.LoadDefaultConfig(context.TODO(), config.WithRegion("eu-north-1"))
	tableBasics := TableBasics{TableName: tableName,
		DynamoDbClient: dynamodb.NewFromConfig(config)}

	item, err := attributevalue.MarshalMap(user)

	if err != nil {
		panic(err)
	}

	_, err = tableBasics.DynamoDbClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		Item: item, TableName: aws.String(tableName),
	},
	)

	if err != nil {
		log.Printf("Could not update user record: %v\n", err)
	}

	return err

}

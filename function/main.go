package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var table string
var client *dynamodb.Client

func init() {
	table = os.Getenv("TABLE_NAME")
	if table == "" {
		log.Fatal("missing environment variable TABLE_NAME")
	}
	cfg, _ := config.LoadDefaultConfig(context.Background())
	client = dynamodb.NewFromConfig(cfg)

}

func handler(ctx context.Context, kinesisEvent events.KinesisEvent) error {

	for _, record := range kinesisEvent.Records {
		fmt.Println("received message from kinesis. partition key", record.Kinesis.PartitionKey)
		fmt.Println("storing info to dynamodb table", table)

		data := record.Kinesis.Data

		var user CreateUserInfo
		err := json.Unmarshal(data, &user)

		if err != nil {
			return err
		}

		item, err := attributevalue.MarshalMap(user)
		if err != nil {
			return err
		}

		item["email"] = &types.AttributeValueMemberS{Value: record.Kinesis.PartitionKey}

		_, err = client.PutItem(context.Background(), &dynamodb.PutItemInput{
			TableName: aws.String(table),
			Item:      item,
		})

		if err != nil {
			return err
		}

		fmt.Println("item added to table")
	}

	return nil
}

func main() {
	lambda.Start(handler)
}

type CreateUserInfo struct {
	Name string `json:"name"`
	City string `json:"city"`
}

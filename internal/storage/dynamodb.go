package storage

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	_ "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DynamoDBStorage struct {
	client    *dynamodb.Client
	tableName string
}

func NewDynamoDBStorage(client *dynamodb.Client, tableName string) *DynamoDBStorage {
	return &DynamoDBStorage{
		client:    client,
		tableName: tableName,
	}
}
func (d *DynamoDBStorage) SaveURLMapping(shortCode, longURL string) error {
	_, err := d.client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: &d.tableName,
		Item: map[string]types.AttributeValue{
			"short_code": &types.AttributeValueMemberS{Value: shortCode},
			"long_url":   &types.AttributeValueMemberS{Value: longURL},
		},
	})
	return err
}

func (d *DynamoDBStorage) GetLongURL(shortCode string) (string, error) {
	output, err := d.client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: &d.tableName,
		Key: map[string]types.AttributeValue{
			"short_code": &types.AttributeValueMemberS{Value: shortCode},
		},
	})
	if err != nil {
		return "", err
	}
	if output.Item == nil {
		return "", errors.New("short code not found")
	}

	return output.Item["long_url"].(*types.AttributeValueMemberS).Value, nil
}

func (d *DynamoDBStorage) CheckShortCodeExists(shortCode string) (bool, error) {
	_, err := d.GetLongURL(shortCode)
	if err != nil {
		return false, nil
	}
	return true, nil
}

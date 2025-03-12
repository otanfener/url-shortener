package storage

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	_ "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/otanfener/url-shortener/internal/domain"
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
	if err != nil {
		return fmt.Errorf("%w: failed to save URL mapping for short code: %s", domain.ErrStorageFailure, shortCode)
	}
	return nil
}

func (d *DynamoDBStorage) GetLongURL(shortCode string) (string, error) {
	output, err := d.client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: &d.tableName,
		Key: map[string]types.AttributeValue{
			"short_code": &types.AttributeValueMemberS{Value: shortCode},
		},
	})
	if err != nil {
		return "", fmt.Errorf("%w: failed to get long URL for short code: %s", domain.ErrStorageFailure, shortCode)
	}
	if output.Item == nil {
		return "", domain.ErrShortCodeNotFound
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

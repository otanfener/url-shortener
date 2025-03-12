package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/otanfener/url-shortener/internal/counter"
	"github.com/otanfener/url-shortener/internal/logger"
	"github.com/otanfener/url-shortener/internal/server"
	"github.com/otanfener/url-shortener/internal/service"
	"github.com/otanfener/url-shortener/internal/storage"
	"github.com/redis/go-redis/v9"
)

func main() {
	// Local configuration values (hardcoded for now)
	const (
		appPort       = "8080"
		awsRegion     = "us-east-1"
		dynamoDBTable = "urls"
		dynamoDBLocal = "http://localhost:8000"
		redisAddr     = "localhost:6379"
		redisPassword = "" // No password for local Redis
		redisDB       = 0  // Default DB index
	)

	// Load AWS SDK with local DynamoDB configuration
	awsCFG, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(awsRegion),
		config.WithCredentialsProvider(aws.CredentialsProviderFunc(
			func(ctx context.Context) (aws.Credentials, error) {
				return aws.Credentials{
					AccessKeyID:     "fake", // Fake credentials for DynamoDB Local
					SecretAccessKey: "fake",
					SessionToken:    "",
				}, nil
			})),
		config.WithEndpointResolver(aws.EndpointResolverFunc(
			func(service, region string) (aws.Endpoint, error) {
				if service == dynamodb.ServiceID {
					return aws.Endpoint{URL: dynamoDBLocal}, nil
				}
				return aws.Endpoint{}, fmt.Errorf("unknown endpoint requested")
			})),
	)
	if err != nil {
		log.Fatalf("Failed to load AWS config: %v", err)
	}

	// Initialize DynamoDB Local storage
	dynamoDBClient := dynamodb.NewFromConfig(awsCFG)
	dynamoDBStorage := storage.NewDynamoDBStorage(dynamoDBClient, dynamoDBTable)

	// Initialize Redis counter
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDB,
	})
	redisCounter := counter.NewRedisCounter(redisClient)

	appLogger := logger.NewLogger()
	// Initialize URL Shortener Service
	urlShortenerService := service.NewService(dynamoDBStorage, redisCounter, appLogger)

	// Initialize HTTP server
	httpServer := server.NewServer(urlShortenerService, appLogger)
	appLogger.Info("starting server", map[string]interface{}{"port": appPort})
	if err := httpServer.Open(":" + appPort); err != nil {
		appLogger.Error("server error", map[string]interface{}{"error": err.Error()})
	}

}

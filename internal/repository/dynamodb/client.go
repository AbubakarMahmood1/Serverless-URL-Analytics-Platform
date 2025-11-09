package dynamodb

import (
	"context"

	"github.com/AbubakarMahmood1/Serverless-URL-Analytics-Platform/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// NewClient creates a new DynamoDB client
func NewClient(cfg *config.Config) (*dynamodb.Client, error) {
	// Load AWS configuration
	awsCfg, err := awsconfig.LoadDefaultConfig(context.TODO(),
		awsconfig.WithRegion(cfg.AWSRegion),
		awsconfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				cfg.AWSAccessKeyID,
				cfg.AWSSecretAccessKey,
				"",
			),
		),
	)
	if err != nil {
		return nil, err
	}

	// Create DynamoDB client
	client := dynamodb.NewFromConfig(awsCfg)

	return client, nil
}

// NewLocalClient creates a DynamoDB client for local development (DynamoDB Local)
func NewLocalClient(endpoint string) (*dynamodb.Client, error) {
	awsCfg, err := awsconfig.LoadDefaultConfig(context.TODO(),
		awsconfig.WithRegion("us-east-1"),
		awsconfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider("dummy", "dummy", ""),
		),
	)
	if err != nil {
		return nil, err
	}

	client := dynamodb.NewFromConfig(awsCfg, func(o *dynamodb.Options) {
		o.BaseEndpoint = aws.String(endpoint)
	})

	return client, nil
}

package dynamodb

import (
	"context"
	"errors"
	"fmt"

	"github.com/AbubakarMahmood1/Serverless-URL-Analytics-Platform/internal/models"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var (
	ErrURLNotFound      = errors.New("URL not found")
	ErrURLAlreadyExists = errors.New("short code already exists")
)

type URLRepository struct {
	client    *dynamodb.Client
	tableName string
}

func NewURLRepository(client *dynamodb.Client, tableName string) *URLRepository {
	return &URLRepository{
		client:    client,
		tableName: tableName,
	}
}

// Create creates a new shortened URL entry
func (r *URLRepository) Create(ctx context.Context, url *models.URL) error {
	item, err := attributevalue.MarshalMap(url)
	if err != nil {
		return fmt.Errorf("failed to marshal URL: %w", err)
	}

	// Use conditional expression to prevent overwriting
	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           aws.String(r.tableName),
		Item:                item,
		ConditionExpression: aws.String("attribute_not_exists(short_code)"),
	})

	if err != nil {
		var condErr *types.ConditionalCheckFailedException
		if errors.As(err, &condErr) {
			return ErrURLAlreadyExists
		}
		return fmt.Errorf("failed to create URL: %w", err)
	}

	return nil
}

// GetByShortCode retrieves a URL by its short code
func (r *URLRepository) GetByShortCode(ctx context.Context, shortCode string) (*models.URL, error) {
	result, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"short_code": &types.AttributeValueMemberS{Value: shortCode},
		},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get URL: %w", err)
	}

	if result.Item == nil {
		return nil, ErrURLNotFound
	}

	var url models.URL
	if err := attributevalue.UnmarshalMap(result.Item, &url); err != nil {
		return nil, fmt.Errorf("failed to unmarshal URL: %w", err)
	}

	return &url, nil
}

// Update updates an existing URL entry
func (r *URLRepository) Update(ctx context.Context, url *models.URL) error {
	item, err := attributevalue.MarshalMap(url)
	if err != nil {
		return fmt.Errorf("failed to marshal URL: %w", err)
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	})

	if err != nil {
		return fmt.Errorf("failed to update URL: %w", err)
	}

	return nil
}

// Delete deactivates a URL (soft delete)
func (r *URLRepository) Delete(ctx context.Context, shortCode string) error {
	_, err := r.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"short_code": &types.AttributeValueMemberS{Value: shortCode},
		},
		UpdateExpression: aws.String("SET is_active = :inactive"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":inactive": &types.AttributeValueMemberBOOL{Value: false},
		},
	})

	if err != nil {
		return fmt.Errorf("failed to delete URL: %w", err)
	}

	return nil
}

// GetByCreatedBy retrieves all URLs created by a specific user/API key
func (r *URLRepository) GetByCreatedBy(ctx context.Context, createdBy string) ([]*models.URL, error) {
	result, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		IndexName:              aws.String("created_by-index"),
		KeyConditionExpression: aws.String("created_by = :createdBy"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":createdBy": &types.AttributeValueMemberS{Value: createdBy},
		},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to query URLs: %w", err)
	}

	urls := make([]*models.URL, 0, len(result.Items))
	for _, item := range result.Items {
		var url models.URL
		if err := attributevalue.UnmarshalMap(item, &url); err != nil {
			return nil, fmt.Errorf("failed to unmarshal URL: %w", err)
		}
		urls = append(urls, &url)
	}

	return urls, nil
}

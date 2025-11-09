package dynamodb

import (
	"context"
	"fmt"
	"strconv"

	"github.com/AbubakarMahmood1/Serverless-URL-Analytics-Platform/internal/models"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type AnalyticsRepository struct {
	client    *dynamodb.Client
	tableName string
}

func NewAnalyticsRepository(client *dynamodb.Client, tableName string) *AnalyticsRepository {
	return &AnalyticsRepository{
		client:    client,
		tableName: tableName,
	}
}

// RecordClick records a click event
func (r *AnalyticsRepository) RecordClick(ctx context.Context, event *models.ClickEvent) error {
	// Create composite sort key: timestamp#click_id
	sortKey := fmt.Sprintf("%d#%s", event.Timestamp, event.ClickID)

	item, err := attributevalue.MarshalMap(event)
	if err != nil {
		return fmt.Errorf("failed to marshal click event: %w", err)
	}

	// Add the composite sort key
	item["sort_key"] = &types.AttributeValueMemberS{Value: sortKey}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	})

	if err != nil {
		return fmt.Errorf("failed to record click: %w", err)
	}

	return nil
}

// GetClicksByShortCode retrieves all clicks for a short code
func (r *AnalyticsRepository) GetClicksByShortCode(ctx context.Context, shortCode string, limit int) ([]*models.ClickEvent, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		KeyConditionExpression: aws.String("short_code = :shortCode"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":shortCode": &types.AttributeValueMemberS{Value: shortCode},
		},
		ScanIndexForward: aws.Bool(false), // Sort descending (newest first)
	}

	if limit > 0 {
		input.Limit = aws.Int32(int32(limit))
	}

	result, err := r.client.Query(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to query clicks: %w", err)
	}

	clicks := make([]*models.ClickEvent, 0, len(result.Items))
	for _, item := range result.Items {
		var click models.ClickEvent
		if err := attributevalue.UnmarshalMap(item, &click); err != nil {
			return nil, fmt.Errorf("failed to unmarshal click: %w", err)
		}
		clicks = append(clicks, &click)
	}

	return clicks, nil
}

// GetClicksByTimeRange retrieves clicks within a time range
func (r *AnalyticsRepository) GetClicksByTimeRange(ctx context.Context, shortCode string, startTime, endTime int64) ([]*models.ClickEvent, error) {
	result, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		KeyConditionExpression: aws.String("short_code = :shortCode AND sort_key BETWEEN :start AND :end"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":shortCode": &types.AttributeValueMemberS{Value: shortCode},
			":start":     &types.AttributeValueMemberS{Value: strconv.FormatInt(startTime, 10)},
			":end":       &types.AttributeValueMemberS{Value: strconv.FormatInt(endTime, 10)},
		},
		ScanIndexForward: aws.Bool(false),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to query clicks by time range: %w", err)
	}

	clicks := make([]*models.ClickEvent, 0, len(result.Items))
	for _, item := range result.Items {
		var click models.ClickEvent
		if err := attributevalue.UnmarshalMap(item, &click); err != nil {
			return nil, fmt.Errorf("failed to unmarshal click: %w", err)
		}
		clicks = append(clicks, &click)
	}

	return clicks, nil
}

// GetTotalClicks returns the total number of clicks for a short code
func (r *AnalyticsRepository) GetTotalClicks(ctx context.Context, shortCode string) (int, error) {
	result, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		KeyConditionExpression: aws.String("short_code = :shortCode"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":shortCode": &types.AttributeValueMemberS{Value: shortCode},
		},
		Select: types.SelectCount,
	})

	if err != nil {
		return 0, fmt.Errorf("failed to count clicks: %w", err)
	}

	return int(result.Count), nil
}

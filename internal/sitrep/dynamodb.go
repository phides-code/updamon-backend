// DynamoDB implementation of Repository for the sitreps table.
package sitrep

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/phides-code/updamon-backend/internal/domain"
)

// DynamoDB attribute names. Keep aligned with Sitrep json / dynamodbav tags.
const (
	AttrID        = "id"
	AttrHostname  = "hostname"
	AttrAptLog    = "aptlog"
	AttrCreatedOn = "createdOn"
)

// Condition expressions on the partition key: create requires absence.
const ConditionIDNotExists = "attribute_not_exists(" + AttrID + ")"

type dynamoRepository struct {
	client dynamoAPI
}

type dynamoAPI interface {
	GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	Scan(ctx context.Context, params *dynamodb.ScanInput, optFns ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error)
}

func NewRepository(client dynamoAPI) Repository {
	return &dynamoRepository{client: client}
}

func idKey(id string) map[string]types.AttributeValue {
	return map[string]types.AttributeValue{
		AttrID: &types.AttributeValueMemberS{Value: id},
	}
}

func isConditionalCheckFailed(err error) bool {
	var conditionalCheck *types.ConditionalCheckFailedException
	return errors.As(err, &conditionalCheck)
}

func unmarshalSitrep(item map[string]types.AttributeValue) (Sitrep, error) {
	var s Sitrep
	if err := attributevalue.UnmarshalMap(item, &s); err != nil {
		return Sitrep{}, fmt.Errorf("unmarshal sitrep: %w", err)
	}
	return s, nil
}

func (r *dynamoRepository) Create(ctx context.Context, s Sitrep) (Sitrep, error) {
	item, err := attributevalue.MarshalMap(s)
	if err != nil {
		return Sitrep{}, fmt.Errorf("marshal sitrep: %w", err)
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           aws.String(TableName),
		Item:                item,
		ConditionExpression: aws.String(ConditionIDNotExists),
	})

	if err != nil {
		if isConditionalCheckFailed(err) {
			return Sitrep{}, domain.ErrAlreadyExists
		}
		return Sitrep{}, fmt.Errorf("put item: %w", err)
	}

	return s, nil
}

func (r *dynamoRepository) GetByID(ctx context.Context, id string) (Sitrep, error) {
	out, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(TableName),
		Key:       idKey(id),
	})
	if err != nil {
		return Sitrep{}, fmt.Errorf("get item: %w", err)
	}
	if out.Item == nil {
		return Sitrep{}, domain.ErrNotFound
	}

	return unmarshalSitrep(out.Item)
}

func (r *dynamoRepository) List(ctx context.Context) ([]Sitrep, error) {
	var items []Sitrep
	var startKey map[string]types.AttributeValue

	for {
		input := &dynamodb.ScanInput{
			TableName: aws.String(TableName),
		}
		if startKey != nil {
			input.ExclusiveStartKey = startKey
		}

		out, err := r.client.Scan(ctx, input)
		if err != nil {
			return nil, fmt.Errorf("scan items: %w", err)
		}

		for _, item := range out.Items {
			s, err := unmarshalSitrep(item)
			if err != nil {
				return nil, err
			}
			items = append(items, s)
		}

		if out.LastEvaluatedKey == nil {
			break
		}
		startKey = out.LastEvaluatedKey
	}

	return items, nil
}

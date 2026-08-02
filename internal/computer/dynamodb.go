// DynamoDB implementation of Repository for the computers table.
package computer

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

// DynamoDB attribute names. Keep aligned with Computer json / dynamodbav tags.
const (
	AttrID        = "id"
	AttrHostname  = "hostname"
	AttrIP        = "ip"
	AttrRating    = "rating"
	AttrCreatedOn = "createdOn"
)

// Condition expressions on the partition key: create requires absence, update requires presence.
const (
	ConditionIDNotExists = "attribute_not_exists(" + AttrID + ")"
	ConditionIDExists    = "attribute_exists(" + AttrID + ")"
)

type dynamoRepository struct {
	client dynamoAPI
}

type dynamoAPI interface {
	GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	Scan(ctx context.Context, params *dynamodb.ScanInput, optFns ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error)
	UpdateItem(ctx context.Context, params *dynamodb.UpdateItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.UpdateItemOutput, error)
	DeleteItem(ctx context.Context, params *dynamodb.DeleteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error)
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

func unmarshalComputer(item map[string]types.AttributeValue) (Computer, error) {
	var computer Computer
	if err := attributevalue.UnmarshalMap(item, &computer); err != nil {
		return Computer{}, fmt.Errorf("unmarshal computer: %w", err)
	}
	return computer, nil
}

func (r *dynamoRepository) Create(ctx context.Context, computer Computer) (Computer, error) {
	item, err := attributevalue.MarshalMap(computer)
	if err != nil {
		return Computer{}, fmt.Errorf("marshal computer: %w", err)
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           aws.String(TableName),
		Item:                item,
		ConditionExpression: aws.String(ConditionIDNotExists),
	})

	if err != nil {
		if isConditionalCheckFailed(err) {
			return Computer{}, domain.ErrAlreadyExists
		}
		return Computer{}, fmt.Errorf("put item: %w", err)
	}

	return computer, nil
}

func (r *dynamoRepository) GetByID(ctx context.Context, id string) (Computer, error) {
	out, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(TableName),
		Key:       idKey(id),
	})
	if err != nil {
		return Computer{}, fmt.Errorf("get item: %w", err)
	}
	if out.Item == nil {
		return Computer{}, domain.ErrNotFound
	}

	return unmarshalComputer(out.Item)
}

func (r *dynamoRepository) List(ctx context.Context) ([]Computer, error) {
	var items []Computer
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
			computer, err := unmarshalComputer(item)
			if err != nil {
				return nil, err
			}
			items = append(items, computer)
		}

		if out.LastEvaluatedKey == nil {
			break
		}
		startKey = out.LastEvaluatedKey
	}

	return items, nil
}

func (r *dynamoRepository) Update(ctx context.Context, computer Computer) (Computer, error) {
	out, err := r.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(TableName),
		Key:       idKey(computer.ID),
		UpdateExpression: aws.String(fmt.Sprintf(
			"SET #%s = :%s, "+
				"#%s = :%s, "+
				"#%s = :%s",
			AttrHostname, AttrHostname,
			AttrIP, AttrIP,
			AttrRating, AttrRating,
		)),
		ConditionExpression: aws.String(ConditionIDExists),
		ExpressionAttributeNames: map[string]string{
			"#" + AttrHostname: AttrHostname,
			"#" + AttrIP:       AttrIP,
			"#" + AttrRating:   AttrRating,
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":" + AttrHostname: &types.AttributeValueMemberS{Value: computer.Hostname},
			":" + AttrIP:       &types.AttributeValueMemberS{Value: computer.IP},
			":" + AttrRating:   &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", computer.Rating)},
		},
		ReturnValues: types.ReturnValueAllNew,
	})
	if err != nil {
		if isConditionalCheckFailed(err) {
			return Computer{}, domain.ErrNotFound
		}
		return Computer{}, fmt.Errorf("update item: %w", err)
	}

	return unmarshalComputer(out.Attributes)
}

func (r *dynamoRepository) Delete(ctx context.Context, id string) (Computer, error) {
	out, err := r.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName:    aws.String(TableName),
		Key:          idKey(id),
		ReturnValues: types.ReturnValueAllOld,
	})
	if err != nil {
		return Computer{}, fmt.Errorf("delete item: %w", err)
	}
	if out.Attributes == nil {
		return Computer{}, domain.ErrNotFound
	}

	return unmarshalComputer(out.Attributes)
}

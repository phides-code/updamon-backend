// Repository interface; DynamoDB implementation is in dynamodb.go.
package computer

import "context"

type Repository interface {
	Create(ctx context.Context, computer Computer) (Computer, error)
	GetByID(ctx context.Context, id string) (Computer, error)
	List(ctx context.Context) ([]Computer, error)
	Update(ctx context.Context, computer Computer) (Computer, error)
	Delete(ctx context.Context, id string) (Computer, error)
}

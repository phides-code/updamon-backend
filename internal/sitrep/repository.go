// Repository interface for sitrep persistence. DynamoDB implements this in dynamodb.go.
package sitrep

import "context"

type Repository interface {
	Create(ctx context.Context, sitrep Sitrep) (Sitrep, error)
	GetByID(ctx context.Context, id string) (Sitrep, error)
	List(ctx context.Context) ([]Sitrep, error)
}

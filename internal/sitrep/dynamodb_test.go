// Unit tests for the sitrep DynamoDB repository using a mocked DynamoDB client.
package sitrep_test

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	awsdynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/phides-code/updamon-backend/internal/sitrep"
	"github.com/phides-code/updamon-backend/internal/domain"
	"github.com/phides-code/updamon-backend/internal/testutil"
)

// errDynamoUnavailable is the stub SDK failure used across repository mock cases.
var errDynamoUnavailable = errors.New("dynamo unavailable")

type mockDynamoClient struct {
	getItemFn func(ctx context.Context, params *awsdynamodb.GetItemInput, optFns ...func(*awsdynamodb.Options)) (*awsdynamodb.GetItemOutput, error)
	putItemFn func(ctx context.Context, params *awsdynamodb.PutItemInput, optFns ...func(*awsdynamodb.Options)) (*awsdynamodb.PutItemOutput, error)
	scanFn    func(ctx context.Context, params *awsdynamodb.ScanInput, optFns ...func(*awsdynamodb.Options)) (*awsdynamodb.ScanOutput, error)
}

func (m *mockDynamoClient) GetItem(ctx context.Context, params *awsdynamodb.GetItemInput, optFns ...func(*awsdynamodb.Options)) (*awsdynamodb.GetItemOutput, error) {
	return m.getItemFn(ctx, params, optFns...)
}

func (m *mockDynamoClient) PutItem(ctx context.Context, params *awsdynamodb.PutItemInput, optFns ...func(*awsdynamodb.Options)) (*awsdynamodb.PutItemOutput, error) {
	return m.putItemFn(ctx, params, optFns...)
}

func (m *mockDynamoClient) Scan(ctx context.Context, params *awsdynamodb.ScanInput, optFns ...func(*awsdynamodb.Options)) (*awsdynamodb.ScanOutput, error) {
	return m.scanFn(ctx, params, optFns...)
}



func scanItems(t *testing.T, sitreps []sitrep.Sitrep) []map[string]types.AttributeValue {
	t.Helper()
	items := make([]map[string]types.AttributeValue, len(sitreps))
	for i, b := range sitreps {
		item, err := attributevalue.MarshalMap(b)
		if err != nil {
			t.Fatal(err)
		}
		items[i] = item
	}
	return items
}

func TestSitrepRepositoryGetByID(t *testing.T) {
	t.Parallel()

	validId, validSitrep, item := storedSitrepFixture(t)
	tests := []struct {
		name       string
		setupMock  func(t *testing.T) *mockDynamoClient
		wantSitrep sitrep.Sitrep
		wantErr    error
	}{
		{
			name: "found",
			setupMock: func(_ *testing.T) *mockDynamoClient {
				return &mockDynamoClient{
					getItemFn: func(_ context.Context, _ *awsdynamodb.GetItemInput, _ ...func(*awsdynamodb.Options)) (*awsdynamodb.GetItemOutput, error) {
						return &awsdynamodb.GetItemOutput{Item: item}, nil
					},
				}
			},
			wantSitrep: validSitrep,
			wantErr:    nil,
		},
		{
			name: "not found",
			setupMock: func(_ *testing.T) *mockDynamoClient {
				return &mockDynamoClient{
					getItemFn: func(_ context.Context, _ *awsdynamodb.GetItemInput, _ ...func(*awsdynamodb.Options)) (*awsdynamodb.GetItemOutput, error) {
						return &awsdynamodb.GetItemOutput{Item: nil}, nil
					},
				}
			},
			wantSitrep: sitrep.Sitrep{},
			wantErr:    domain.ErrNotFound,
		},
		{
			name: "sdk error",
			setupMock: func(_ *testing.T) *mockDynamoClient {
				return &mockDynamoClient{
					getItemFn: func(_ context.Context, _ *awsdynamodb.GetItemInput, _ ...func(*awsdynamodb.Options)) (*awsdynamodb.GetItemOutput, error) {
						return nil, errDynamoUnavailable
					},
				}
			},
			wantSitrep: sitrep.Sitrep{},
			wantErr:    errDynamoUnavailable,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := sitrep.NewRepository(tt.setupMock(t))
			got, err := repo.GetByID(context.Background(), validId)

			assertSitrepRepoResult(t, "GetByID", got, err, tt.wantSitrep, tt.wantErr)
		})
	}
}



func TestSitrepRepositoryCreate(t *testing.T) {
	t.Parallel()

	_, want := testutil.SitrepWithID(
		testutil.ValidSitrepBody(),
		testutil.TestStoredSitrepCreatedOn,
	)
	tests := []struct {
		name       string
		setupMock  func(t *testing.T) *mockDynamoClient
		wantSitrep sitrep.Sitrep
		wantErr    error
	}{
		{
			name: "success",
			setupMock: func(t *testing.T) *mockDynamoClient {
				return &mockDynamoClient{
					putItemFn: func(_ context.Context, params *awsdynamodb.PutItemInput, _ ...func(*awsdynamodb.Options)) (*awsdynamodb.PutItemOutput, error) {
						assertSitrepPutItem(t, params, want)
						return &awsdynamodb.PutItemOutput{}, nil
					},
				}
			},
			wantSitrep: want,
			wantErr:    nil,
		},
		{
			name: "duplicate id",
			setupMock: func(_ *testing.T) *mockDynamoClient {
				return &mockDynamoClient{
					putItemFn: func(_ context.Context, _ *awsdynamodb.PutItemInput, _ ...func(*awsdynamodb.Options)) (*awsdynamodb.PutItemOutput, error) {
						return nil, &types.ConditionalCheckFailedException{}
					},
				}
			},
			wantSitrep: sitrep.Sitrep{},
			wantErr:    domain.ErrAlreadyExists,
		},
		{
			name: "sdk error",
			setupMock: func(_ *testing.T) *mockDynamoClient {
				return &mockDynamoClient{
					putItemFn: func(_ context.Context, _ *awsdynamodb.PutItemInput, _ ...func(*awsdynamodb.Options)) (*awsdynamodb.PutItemOutput, error) {
						return nil, errDynamoUnavailable
					},
				}
			},
			wantSitrep: sitrep.Sitrep{},
			wantErr:    errDynamoUnavailable,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := sitrep.NewRepository(tt.setupMock(t))
			got, err := repo.Create(context.Background(), want)

			assertSitrepRepoResult(t, "Create", got, err, tt.wantSitrep, tt.wantErr)
		})
	}
}

func TestSitrepRepositoryList(t *testing.T) {
	t.Parallel()

	b1, b2, b3 := testutil.ListSitreps(true)
	wantItems := []sitrep.Sitrep{b1, b2}
	page2 := []sitrep.Sitrep{b3}
	scanOutputItems := scanItems(t, wantItems)
	page2ScanItems := scanItems(t, page2)

	tests := []struct {
		name      string
		setupMock func(t *testing.T) *mockDynamoClient
		wantItems []sitrep.Sitrep
		wantErr   bool
	}{
		{
			name: "returns items",
			setupMock: func(_ *testing.T) *mockDynamoClient {
				return &mockDynamoClient{
					scanFn: func(_ context.Context, params *awsdynamodb.ScanInput, _ ...func(*awsdynamodb.Options)) (*awsdynamodb.ScanOutput, error) {
						if params.Limit != nil {
							t.Errorf("Limit = %v, want nil", params.Limit)
						}
						return &awsdynamodb.ScanOutput{Items: scanOutputItems}, nil
					},
				}
			},
			wantItems: wantItems,
		},
		{
			name: "empty",
			setupMock: func(_ *testing.T) *mockDynamoClient {
				return &mockDynamoClient{
					scanFn: func(_ context.Context, _ *awsdynamodb.ScanInput, _ ...func(*awsdynamodb.Options)) (*awsdynamodb.ScanOutput, error) {
						return &awsdynamodb.ScanOutput{Items: nil}, nil
					},
				}
			},
			wantItems: []sitrep.Sitrep{},
		},
		{
			name: "scans all pages",
			setupMock: func(_ *testing.T) *mockDynamoClient {
				calls := 0
				return &mockDynamoClient{
					scanFn: func(_ context.Context, params *awsdynamodb.ScanInput, _ ...func(*awsdynamodb.Options)) (*awsdynamodb.ScanOutput, error) {
						calls++
						switch calls {
						case 1:
							if params.ExclusiveStartKey != nil {
								t.Fatal("expected first scan without ExclusiveStartKey")
							}
							return &awsdynamodb.ScanOutput{
								Items: scanOutputItems,
								LastEvaluatedKey: map[string]types.AttributeValue{
									sitrep.AttrID: &types.AttributeValueMemberS{Value: b2.ID},
								},
							}, nil
						case 2:
							if params.ExclusiveStartKey == nil {
								t.Fatal("expected second scan with ExclusiveStartKey")
							}
							return &awsdynamodb.ScanOutput{Items: page2ScanItems}, nil
						default:
							t.Fatal("unexpected extra scan")
							return nil, nil
						}
					},
				}
			},
			wantItems: append(wantItems, page2...),
		},
		{
			name: "sdk error",
			setupMock: func(_ *testing.T) *mockDynamoClient {
				return &mockDynamoClient{
					scanFn: func(_ context.Context, _ *awsdynamodb.ScanInput, _ ...func(*awsdynamodb.Options)) (*awsdynamodb.ScanOutput, error) {
						return nil, errDynamoUnavailable
					},
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := sitrep.NewRepository(tt.setupMock(t))
			items, err := repo.List(context.Background())

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("List: %v", err)
			}

			if len(items) != len(tt.wantItems) {
				t.Fatalf("len(items) = %d, want %d", len(items), len(tt.wantItems))
			}

			for i := range tt.wantItems {
				if items[i] != tt.wantItems[i] {
					t.Fatalf("items[%d] = %+v, want %+v", i, items[i], tt.wantItems[i])
				}
			}
		})
	}
}

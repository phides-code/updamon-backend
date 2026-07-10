// Computer-specific test helpers for HTTP wire shape and DynamoDB repository assertions.
package computer_test

import (
	"encoding/json"
	"errors"
	"maps"
	"testing"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	awsdynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/phides-code/go-multi-api/internal/computer"
	"github.com/phides-code/go-multi-api/internal/platform"
)

func decodeComputerData(t *testing.T, envelope platform.APIResponse) computer.Computer {
	t.Helper()
	if envelope.Error != nil {
		t.Fatalf("unexpected error: %s", *envelope.Error)
	}
	data, err := json.Marshal(envelope.Data)
	if err != nil {
		t.Fatalf("marshal data: %v", err)
	}
	var b computer.Computer
	if err := json.Unmarshal(data, &b); err != nil {
		t.Fatalf("unmarshal computer: %v", err)
	}
	return b
}

func decodeComputerListData(t *testing.T, envelope platform.APIResponse) []computer.Computer {
	t.Helper()
	if envelope.Error != nil {
		t.Fatalf("unexpected error: %s", *envelope.Error)
	}
	data, err := json.Marshal(envelope.Data)
	if err != nil {
		t.Fatalf("marshal data: %v", err)
	}
	var items []computer.Computer
	if err := json.Unmarshal(data, &items); err != nil {
		t.Fatalf("unmarshal list: %v", err)
	}
	return items
}

func assertComputerDataKeys(t *testing.T, envelope platform.APIResponse) {
	t.Helper()

	raw, err := json.Marshal(envelope.Data)
	if err != nil {
		t.Fatalf("marshal data: %v", err)
	}

	var keys map[string]json.RawMessage
	if err := json.Unmarshal(raw, &keys); err != nil {
		t.Fatalf("unmarshal data keys: %v", err)
	}

	want := []string{"hostname", "createdOn", "id"}
	if len(keys) != len(want) {
		t.Fatalf("data has %d keys %v, want exactly %v", len(keys), maps.Keys(keys), want)
	}
	for _, k := range want {
		if _, ok := keys[k]; !ok {
			t.Fatalf("missing data key %q; got %v", k, maps.Keys(keys))
		}
	}
}

func assertComputerPutItem(t *testing.T, params *awsdynamodb.PutItemInput, want computer.Computer) {
	t.Helper()

	if params.ConditionExpression == nil || *params.ConditionExpression != "attribute_not_exists(id)" {
		t.Fatalf("ConditionExpression = %v, want attribute_not_exists(id)", params.ConditionExpression)
	}

	var got computer.Computer
	if err := attributevalue.UnmarshalMap(params.Item, &got); err != nil {
		t.Fatalf("unmarshal item: %v", err)
	}
	if got != want {
		t.Fatalf("Item = %+v, want %+v", got, want)
	}
}

func assertComputerRepoResult(t *testing.T, op string, got computer.Computer, err error, want computer.Computer, wantErr error) {
	t.Helper()

	if wantErr != nil {
		if !errors.Is(err, wantErr) {
			t.Fatalf("err = %v, want %v", err, wantErr)
		}
		if got != want {
			t.Fatalf("got %+v, want %+v", got, want)
		}
		return
	}

	if err != nil {
		t.Fatalf("%s: %v", op, err)
	}
	if got != want {
		t.Fatalf("got %+v, want %+v", got, want)
	}
}

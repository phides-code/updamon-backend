// Sitrep-specific asserts for HTTP response wire shape and DynamoDB repository mocks.
package sitrep_test

import (
	"encoding/json"
	"errors"
	"maps"
	"testing"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	awsdynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/phides-code/updamon-backend/internal/sitrep"
	"github.com/phides-code/updamon-backend/internal/platform"
)

func requireEnvelopeDataJSON(t *testing.T, envelope platform.APIResponse) []byte {
	t.Helper()
	if envelope.Error != nil {
		t.Fatalf("unexpected error: %s", *envelope.Error)
	}
	data, err := json.Marshal(envelope.Data)
	if err != nil {
		t.Fatalf("marshal data: %v", err)
	}
	return data
}

func decodeSitrepData(t *testing.T, envelope platform.APIResponse) sitrep.Sitrep {
	t.Helper()
	var b sitrep.Sitrep
	if err := json.Unmarshal(requireEnvelopeDataJSON(t, envelope), &b); err != nil {
		t.Fatalf("unmarshal sitrep: %v", err)
	}
	return b
}

func decodeSitrepListData(t *testing.T, envelope platform.APIResponse) []sitrep.Sitrep {
	t.Helper()
	var items []sitrep.Sitrep
	if err := json.Unmarshal(requireEnvelopeDataJSON(t, envelope), &items); err != nil {
		t.Fatalf("unmarshal list: %v", err)
	}
	return items
}

func assertSitrepDataKeys(t *testing.T, envelope platform.APIResponse) {
	t.Helper()

	var keys map[string]json.RawMessage
	if err := json.Unmarshal(requireEnvelopeDataJSON(t, envelope), &keys); err != nil {
		t.Fatalf("unmarshal data keys: %v", err)
	}

	// Expected JSON keys for a sitrep item (Attr* must stay aligned with Sitrep tags).
	// Keep this list alphabetical so missing/extra keys are easy to spot.
	want := []string{
		sitrep.AttrAptLog,
		sitrep.AttrCreatedOn,
		sitrep.AttrHostname,
		sitrep.AttrID,
	}
	if len(keys) != len(want) {
		t.Fatalf("data has %d keys %v, want exactly %v", len(keys), maps.Keys(keys), want)
	}
	for _, k := range want {
		if _, ok := keys[k]; !ok {
			t.Fatalf("missing data key %q; got %v", k, maps.Keys(keys))
		}
	}
}

func assertSitrepPutItem(t *testing.T, params *awsdynamodb.PutItemInput, want sitrep.Sitrep) {
	t.Helper()

	if params.ConditionExpression == nil || *params.ConditionExpression != sitrep.ConditionIDNotExists {
		t.Fatalf("ConditionExpression = %v, want %s", params.ConditionExpression, sitrep.ConditionIDNotExists)
	}

	var got sitrep.Sitrep
	if err := attributevalue.UnmarshalMap(params.Item, &got); err != nil {
		t.Fatalf("unmarshal item: %v", err)
	}
	if got != want {
		t.Fatalf("Item = %+v, want %+v", got, want)
	}
}

func assertSitrepRepoResult(t *testing.T, op string, got sitrep.Sitrep, err error, want sitrep.Sitrep, wantErr error) {
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

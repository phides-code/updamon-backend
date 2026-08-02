// DynamoDB fixture: a persisted sitrep plus its marshaled AttributeValue map.
package sitrep_test

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/phides-code/updamon-backend/internal/sitrep"
	"github.com/phides-code/updamon-backend/internal/testutil"
)

func storedSitrepFixture(t *testing.T) (id string, b sitrep.Sitrep, item map[string]types.AttributeValue) {
	t.Helper()
	id, b = testutil.SitrepWithID(testutil.ValidSitrepBody(), testutil.TestStoredSitrepCreatedOn)
	var err error
	item, err = attributevalue.MarshalMap(b)
	if err != nil {
		t.Fatal(err)
	}
	return
}

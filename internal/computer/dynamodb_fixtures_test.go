// DynamoDB fixture: a persisted computer plus its marshaled AttributeValue map.
package computer_test

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/phides-code/updamon-backend/internal/computer"
	"github.com/phides-code/updamon-backend/internal/testutil"
)

func storedComputerFixture(t *testing.T) (id string, b computer.Computer, item map[string]types.AttributeValue) {
	t.Helper()
	id, b = testutil.ComputerWithID(testutil.ValidComputerBody(), testutil.TestStoredComputerCreatedOn)
	var err error
	item, err = attributevalue.MarshalMap(b)
	if err != nil {
		t.Fatal(err)
	}
	return
}

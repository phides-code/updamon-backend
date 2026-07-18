// DynamoDB test fixture: persisted computer row plus marshaled item for Get/Delete mocks.
package computer_test

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/phides-code/go-multi-api/internal/computer"
	"github.com/phides-code/go-multi-api/internal/testutil"
)

func storedComputerFixture(t *testing.T) (id string, b computer.Computer, item map[string]types.AttributeValue) {
	t.Helper()
	id, b = testutil.ComputerWithID(testutil.TestComputerHostname, testutil.TestComputerIP, testutil.TestComputerOS, testutil.TestStoredComputerCreatedOn)
	var err error
	item, err = attributevalue.MarshalMap(b)
	if err != nil {
		t.Fatal(err)
	}
	return
}

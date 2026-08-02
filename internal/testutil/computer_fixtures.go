// Shared computer fixtures used by handler and DynamoDB tests outside package computer.
package testutil

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/phides-code/updamon-backend/internal/computer"
)

// TestComputerHostname is the canonical valid hostname in handler and DynamoDB tests.
const TestComputerHostname = "cavendish"

// TestComputerIP is the canonical valid IPv4 address in handler and DynamoDB tests.
const TestComputerIP = "192.168.1.10"

// TestComputerRating is the canonical valid rating in handler and DynamoDB tests.
const TestComputerRating = 50

// TestStoredComputerCreatedOn is a fixed timestamp for persisted-computer repository tests.
const TestStoredComputerCreatedOn uint64 = 12345

const (
	ListComputerHostnameFirst  = TestComputerHostname
	ListComputerHostnameSecond = "plantain"
	ListComputerHostnameThird  = "burro"

	ListComputerIPFirst  = TestComputerIP
	ListComputerIPSecond = "10.0.0.2"
	ListComputerIPThird  = "172.16.0.3"

	ListComputerRatingFirst  = 10
	ListComputerRatingSecond = 20
	ListComputerRatingThird  = 30
)

// ComputerBody is the create/update JSON payload for computer HTTP tests.
// Kept separate from computer.Computer so json-tag drift fails tests instead of silently matching.
type ComputerBody struct {
	Hostname string `json:"hostname"`
	IP       string `json:"ip"`
	Rating   int    `json:"rating"`
}

// ValidComputerBody returns a ComputerBody with canonical valid field values.
func ValidComputerBody() ComputerBody {
	return ComputerBody{
		Hostname: TestComputerHostname,
		IP:       TestComputerIP,
		Rating:   TestComputerRating,
	}
}

// JSON marshals the body to a request payload string.
func (b ComputerBody) JSON(t *testing.T) string {
	t.Helper()
	data, err := json.Marshal(b)
	if err != nil {
		t.Fatalf("marshal computer body: %v", err)
	}
	return string(data)
}

// ComputerWithID builds a computer from a request body and fixed createdOn; returns the same id.
func ComputerWithID(body ComputerBody, createdOn uint64) (id string, b computer.Computer) {
	id = uuid.NewString()
	b = computer.Computer{
		ID:        id,
		Hostname:  body.Hostname,
		IP:        body.IP,
		Rating:    body.Rating,
		CreatedOn: createdOn,
	}
	return
}

// ListComputers returns three distinct list fixtures.
// When withTimestamps is true, CreatedOn is 1, 2, and 3.
func ListComputers(withTimestamps bool) (first, second, third computer.Computer) {
	first = computer.Computer{
		ID:       uuid.NewString(),
		Hostname: ListComputerHostnameFirst,
		IP:       ListComputerIPFirst,
		Rating:   ListComputerRatingFirst,
	}
	second = computer.Computer{
		ID:       uuid.NewString(),
		Hostname: ListComputerHostnameSecond,
		IP:       ListComputerIPSecond,
		Rating:   ListComputerRatingSecond,
	}
	third = computer.Computer{
		ID:       uuid.NewString(),
		Hostname: ListComputerHostnameThird,
		IP:       ListComputerIPThird,
		Rating:   ListComputerRatingThird,
	}
	if withTimestamps {
		first.CreatedOn = 1
		second.CreatedOn = 2
		third.CreatedOn = 3
	}
	return
}

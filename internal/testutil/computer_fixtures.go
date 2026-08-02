// Shared computer fixtures used by handler and DynamoDB tests outside package computer.
package testutil

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/phides-code/updamon-backend/internal/computer"
)

// TestComputerDescriptor is the canonical valid descriptor in handler and DynamoDB tests.
const TestComputerDescriptor = "cavendish"

// TestComputerRating is the canonical valid rating in handler and DynamoDB tests.
const TestComputerRating = 50

// TestStoredComputerCreatedOn is a fixed timestamp for persisted-computer repository tests.
const TestStoredComputerCreatedOn uint64 = 12345

const (
	ListComputerDescriptorFirst  = TestComputerDescriptor
	ListComputerDescriptorSecond = "plantain"
	ListComputerDescriptorThird  = "burro"

	ListComputerRatingFirst  = 10
	ListComputerRatingSecond = 20
	ListComputerRatingThird  = 30
)

// ComputerBody is the create/update JSON payload for computer HTTP tests.
// Kept separate from computer.Computer so json-tag drift fails tests instead of silently matching.
type ComputerBody struct {
	Descriptor string `json:"descriptor"`
	Rating     int    `json:"rating"`
}

// ValidComputerBody returns a ComputerBody with canonical valid field values.
func ValidComputerBody() ComputerBody {
	return ComputerBody{
		Descriptor: TestComputerDescriptor,
		Rating:     TestComputerRating,
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
		ID:         id,
		Descriptor: body.Descriptor,
		Rating:     body.Rating,
		CreatedOn:  createdOn,
	}
	return
}

// ListComputers returns three distinct list fixtures.
// When withTimestamps is true, CreatedOn is 1, 2, and 3.
func ListComputers(withTimestamps bool) (first, second, third computer.Computer) {
	first = computer.Computer{
		ID:         uuid.NewString(),
		Descriptor: ListComputerDescriptorFirst,
		Rating:     ListComputerRatingFirst,
	}
	second = computer.Computer{
		ID:         uuid.NewString(),
		Descriptor: ListComputerDescriptorSecond,
		Rating:     ListComputerRatingSecond,
	}
	third = computer.Computer{
		ID:         uuid.NewString(),
		Descriptor: ListComputerDescriptorThird,
		Rating:     ListComputerRatingThird,
	}
	if withTimestamps {
		first.CreatedOn = 1
		second.CreatedOn = 2
		third.CreatedOn = 3
	}
	return
}

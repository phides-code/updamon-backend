// Shared computer test fixtures for handler and DynamoDB tests.
package testutil

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/phides-code/go-multi-api/internal/computer"
)

// TestComputerContent is the canonical valid content in handler and DynamoDB tests.
const TestComputerContent = "ripe"

// TestStoredComputerCreatedOn is a fixed timestamp for persisted-computer repository tests.
const TestStoredComputerCreatedOn uint64 = 12345

const (
	ListComputerContentFirst  = "first"
	ListComputerContentSecond = "second"
	ListComputerContentThird  = "third"
)

// ComputerWithID returns a computer whose ID matches the returned id string.
func ComputerWithID(content string, createdOn uint64) (id string, b computer.Computer) {
	id = uuid.NewString()
	b = computer.Computer{ID: id, Content: content, CreatedOn: createdOn}
	return
}

// ComputerCreateBody returns JSON for a valid create/update request body.
func ComputerCreateBody(content string) string {
	return fmt.Sprintf(`{"content":%q}`, content)
}

// ListComputers returns three list items for repository list tests.
// When withTimestamps is true, CreatedOn is set to 1, 2, and 3 respectively.
func ListComputers(withTimestamps bool) (first, second, third computer.Computer) {
	first = computer.Computer{ID: uuid.NewString(), Content: ListComputerContentFirst}
	second = computer.Computer{ID: uuid.NewString(), Content: ListComputerContentSecond}
	third = computer.Computer{ID: uuid.NewString(), Content: ListComputerContentThird}
	if withTimestamps {
		first.CreatedOn = 1
		second.CreatedOn = 2
		third.CreatedOn = 3
	}
	return
}

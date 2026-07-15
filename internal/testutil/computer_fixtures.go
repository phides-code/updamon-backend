// Shared computer test fixtures for handler and DynamoDB tests.
package testutil

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/phides-code/go-multi-api/internal/computer"
)

// TestComputerHostname is the canonical valid hostname in handler and DynamoDB tests.
const TestComputerHostname = "ripe"

// TestComputerIP is the canonical valid IP in handler and DynamoDB tests.
const TestComputerIP = "192.0.2.10"

// TestComputerInvalidIP is a non-IPv4 string used in validation and handler tests.
const TestComputerInvalidIP = "not-an-ip"

// TestStoredComputerCreatedOn is a fixed timestamp for persisted-computer repository tests.
const TestStoredComputerCreatedOn uint64 = 12345

const (
	ListComputerHostnameFirst  = "first"
	ListComputerHostnameSecond = "second"
	ListComputerHostnameThird  = "third"

	ListComputerIPFirst  = "192.0.2.1"
	ListComputerIPSecond = "192.0.2.2"
	ListComputerIPThird  = "192.0.2.3"
)

// ComputerWithID returns a computer whose ID matches the returned id string.
func ComputerWithID(hostname string, ip string, createdOn uint64) (id string, b computer.Computer) {
	id = uuid.NewString()
	b = computer.Computer{ID: id, Hostname: hostname, IP: ip, CreatedOn: createdOn}
	return
}

// ComputerCreateBody returns JSON for a valid create/update request body.
func ComputerCreateBody(hostname, ip string) string {
	return fmt.Sprintf(`{"hostname":%q,"ip":%q}`, hostname, ip)
}

// ListComputers returns three list items for repository list tests.
// When withTimestamps is true, CreatedOn is set to 1, 2, and 3 respectively.
func ListComputers(withTimestamps bool) (first, second, third computer.Computer) {
	first = computer.Computer{ID: uuid.NewString(), Hostname: ListComputerHostnameFirst, IP: ListComputerIPFirst}
	second = computer.Computer{ID: uuid.NewString(), Hostname: ListComputerHostnameSecond, IP: ListComputerIPSecond}
	third = computer.Computer{ID: uuid.NewString(), Hostname: ListComputerHostnameThird, IP: ListComputerIPThird}
	if withTimestamps {
		first.CreatedOn = 1
		second.CreatedOn = 2
		third.CreatedOn = 3
	}
	return
}

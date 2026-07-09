// Package-local computer fixtures for handler tests (ID-linked entity and request bodies).
package computer_test

import (
	"github.com/phides-code/go-multi-api/internal/computer"
	"github.com/phides-code/go-multi-api/internal/testutil"
)

// existingComputerFixture returns an ID-linked computer and matching PUT body for get/update/delete tests.
func existingComputerFixture() (id string, b computer.Computer, updateBody string) {
	id, b = testutil.ComputerWithID(testutil.TestComputerContent, 0)
	updateBody = testutil.ComputerCreateBody(testutil.TestComputerContent)
	return
}

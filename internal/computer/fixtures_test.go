// Package-local fixtures for computer handler tests (existing item + invalid request bodies).
package computer_test

import (
	"strings"
	"testing"

	"github.com/phides-code/updamon-backend/internal/computer"
	"github.com/phides-code/updamon-backend/internal/domain"
	"github.com/phides-code/updamon-backend/internal/testutil"
)

// computerValidationBodies holds invalid create/update JSON for handler client-error tests.
// Names describe the invalidation shape (empty, whitespace, too long, …), not which
// entity field was mutated — keeps find-replace and new fields simple.
type computerValidationBodies struct {
	computerWithEmptyValue    string
	computerWithWhitespace    string
	computerWithValueTooLong  string
	computerWithValueBelowMin string
	computerWithValueAboveMax string
	computerWithInvalidIP     string
}

func newComputerValidationBodies(t *testing.T) computerValidationBodies {
	t.Helper()

	emptyValue := testutil.ValidComputerBody()
	emptyValue.Hostname = ""

	whitespace := testutil.ValidComputerBody()
	whitespace.Hostname = "   "

	valueTooLong := testutil.ValidComputerBody()
	valueTooLong.Hostname = strings.Repeat("a", domain.DefaultMaxStringLength+1)

	valueBelowMin := testutil.ValidComputerBody()
	valueBelowMin.Rating = domain.DefaultMinInt - 1

	valueAboveMax := testutil.ValidComputerBody()
	valueAboveMax.Rating = domain.DefaultMaxInt + 1

	invalidIP := testutil.ValidComputerBody()
	invalidIP.IP = "not-an-ip"

	return computerValidationBodies{
		computerWithEmptyValue:    emptyValue.JSON(t),
		computerWithWhitespace:    whitespace.JSON(t),
		computerWithValueTooLong:  valueTooLong.JSON(t),
		computerWithValueBelowMin: valueBelowMin.JSON(t),
		computerWithValueAboveMax: valueAboveMax.JSON(t),
		computerWithInvalidIP:     invalidIP.JSON(t),
	}
}

// existingComputerFixture returns a UUID, matching computer entity, and valid PUT JSON body.
func existingComputerFixture(t *testing.T) (id string, b computer.Computer, updateBody string) {
	t.Helper()
	id, b = testutil.ComputerWithID(testutil.ValidComputerBody(), 0)
	updateBody = testutil.ValidComputerBody().JSON(t)
	return
}

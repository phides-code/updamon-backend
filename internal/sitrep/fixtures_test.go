// Package-local fixtures for sitrep handler tests (existing item + invalid request bodies).
package sitrep_test

import (
	"strings"
	"testing"

	"github.com/phides-code/updamon-backend/internal/domain"
	"github.com/phides-code/updamon-backend/internal/sitrep"
	"github.com/phides-code/updamon-backend/internal/testutil"
)

// sitrepValidationBodies holds invalid create JSON for handler client-error tests.
// Names describe the invalidation shape (empty, whitespace, too long, …), not which
// entity field was mutated — keeps find-replace and new fields simple.
type sitrepValidationBodies struct {
	sitrepWithEmptyValue   string
	sitrepWithWhitespace   string
	sitrepWithValueTooLong string
	sitrepWithInvalidMAC   string
}

func newSitrepValidationBodies(t *testing.T) sitrepValidationBodies {
	t.Helper()

	emptyValue := testutil.ValidSitrepBody()
	emptyValue.Hostname = ""

	whitespace := testutil.ValidSitrepBody()
	whitespace.Hostname = "   "

	// Long-string fields use DefaultMaxLongStringLength; exercise that bound for the too-long shape.
	valueTooLong := testutil.ValidSitrepBody()
	valueTooLong.AptLog = strings.Repeat("a", domain.DefaultMaxLongStringLength+1)

	invalidMAC := testutil.ValidSitrepBody()
	invalidMAC.WAP = "not-a-mac"

	return sitrepValidationBodies{
		sitrepWithEmptyValue:   emptyValue.JSON(t),
		sitrepWithWhitespace:   whitespace.JSON(t),
		sitrepWithValueTooLong: valueTooLong.JSON(t),
		sitrepWithInvalidMAC:   invalidMAC.JSON(t),
	}
}

// existingSitrepFixture returns a UUID and matching sitrep entity.
func existingSitrepFixture(t *testing.T) (id string, s sitrep.Sitrep) {
	t.Helper()
	id, s = testutil.SitrepWithID(testutil.ValidSitrepBody(), 0)
	return
}

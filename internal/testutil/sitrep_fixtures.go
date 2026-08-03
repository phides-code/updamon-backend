// Shared sitrep fixtures used by handler and DynamoDB tests outside package sitrep.
package testutil

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/phides-code/updamon-backend/internal/sitrep"
)

// TestSitrepHostname is the canonical valid hostname in handler and DynamoDB tests.
const TestSitrepHostname = "cavendish"

// TestSitrepAptLog is the canonical valid aptlog in handler and DynamoDB tests.
const TestSitrepAptLog = "Start-Date: 2026-08-02 12:00:00+0000\nCommandline: apt upgrade"

// TestSitrepLast is the canonical valid last field in handler and DynamoDB tests.
const TestSitrepLast = "last login: Mon Aug  3 08:00:00 2026 from 10.0.0.1"

// TestStoredSitrepCreatedOn is a fixed timestamp for persisted-sitrep repository tests.
const TestStoredSitrepCreatedOn uint64 = 12345

const (
	ListSitrepHostnameFirst  = TestSitrepHostname
	ListSitrepHostnameSecond = "plantain"
	ListSitrepHostnameThird  = "burro"
)

// SitrepBody is the create JSON payload for sitrep HTTP tests.
// Kept separate from sitrep.Sitrep so json-tag drift fails tests instead of silently matching.
type SitrepBody struct {
	Hostname string `json:"hostname"`
	AptLog   string `json:"aptlog"`
	Last     string `json:"last"`
}

// ValidSitrepBody returns a SitrepBody with canonical valid field values.
func ValidSitrepBody() SitrepBody {
	return SitrepBody{
		Hostname: TestSitrepHostname,
		AptLog:   TestSitrepAptLog,
		Last:     TestSitrepLast,
	}
}

// JSON marshals the body to a request payload string.
func (b SitrepBody) JSON(t *testing.T) string {
	t.Helper()
	data, err := json.Marshal(b)
	if err != nil {
		t.Fatalf("marshal sitrep body: %v", err)
	}
	return string(data)
}

// SitrepWithID builds a sitrep from a request body and fixed createdOn; returns the same id.
func SitrepWithID(body SitrepBody, createdOn uint64) (id string, s sitrep.Sitrep) {
	id = uuid.NewString()
	s = sitrep.Sitrep{
		ID:        id,
		Hostname:  body.Hostname,
		AptLog:    body.AptLog,
		Last:      body.Last,
		CreatedOn: createdOn,
	}
	return
}

// ListSitreps returns three distinct list fixtures.
// When withTimestamps is true, CreatedOn is 1, 2, and 3.
func ListSitreps(withTimestamps bool) (first, second, third sitrep.Sitrep) {
	first = sitrep.Sitrep{
		ID:       uuid.NewString(),
		Hostname: ListSitrepHostnameFirst,
		AptLog:   TestSitrepAptLog,
		Last:     TestSitrepLast,
	}
	second = sitrep.Sitrep{
		ID:       uuid.NewString(),
		Hostname: ListSitrepHostnameSecond,
		AptLog:   TestSitrepAptLog,
		Last:     TestSitrepLast,
	}
	third = sitrep.Sitrep{
		ID:       uuid.NewString(),
		Hostname: ListSitrepHostnameThird,
		AptLog:   TestSitrepAptLog,
		Last:     TestSitrepLast,
	}
	if withTimestamps {
		first.CreatedOn = 1
		second.CreatedOn = 2
		third.CreatedOn = 3
	}
	return
}

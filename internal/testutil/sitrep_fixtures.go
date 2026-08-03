// Shared sitrep fixtures used by handler and DynamoDB tests outside package sitrep.
package testutil

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/phides-code/updamon-backend/internal/sitrep"
)

// TestSitrepHostname is the canonical valid hostname in handler and DynamoDB tests.
const TestSitrepHostname = "myHostName"

// TestSitrepAptLog is the canonical valid aptlog in handler and DynamoDB tests.
const TestSitrepAptLog = "Start-Date: 2026-08-02 12:00:00+0000\nCommandline: apt upgrade"

// TestSitrepLast is the canonical valid last field in handler and DynamoDB tests.
const TestSitrepLast = "last login: Mon Aug  3 08:00:00 2026 from 10.0.0.1"

// Canonical valid medium-string field values for handler and DynamoDB tests.
const (
	TestSitrepWAP  = "cc:f4:11:32:cb:ff"
	TestSitrepFree = "Mem: 16G total, 4G used, 12G free"
	TestSitrepDF   = "/dev/nvme0n1p2  500G  120G  380G  24%"
	TestSitrepWho  = "phil pts/0 2026-08-03 08:00"
)

// Canonical valid short-string field values for handler and DynamoDB tests.
const (
	TestSitrepTailscale = "100.64.0.1"
	TestSitrepBluetooth = "up"
)

// TestStoredSitrepCreatedOn is a fixed timestamp for persisted-sitrep repository tests.
const TestStoredSitrepCreatedOn uint64 = 12345

const (
	ListSitrepHostnameFirst  = TestSitrepHostname
	ListSitrepHostnameSecond = "mySecondHostname"
	ListSitrepHostnameThird  = "myThirdHostname"
)

// SitrepBody is the create JSON payload for sitrep HTTP tests.
// Kept separate from sitrep.Sitrep so json-tag drift fails tests instead of silently matching.
type SitrepBody struct {
	Hostname  string `json:"hostname"`
	AptLog    string `json:"aptlog"`
	Last      string `json:"last"`
	WAP       string `json:"wap"`
	Free      string `json:"free"`
	DF        string `json:"df"`
	Who       string `json:"who"`
	Tailscale string `json:"tailscale"`
	Bluetooth string `json:"bluetooth"`
}

// ValidSitrepBody returns a SitrepBody with canonical valid field values.
func ValidSitrepBody() SitrepBody {
	return SitrepBody{
		Hostname:  TestSitrepHostname,
		AptLog:    TestSitrepAptLog,
		Last:      TestSitrepLast,
		WAP:       TestSitrepWAP,
		Free:      TestSitrepFree,
		DF:        TestSitrepDF,
		Who:       TestSitrepWho,
		Tailscale: TestSitrepTailscale,
		Bluetooth: TestSitrepBluetooth,
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
		WAP:       body.WAP,
		Free:      body.Free,
		DF:        body.DF,
		Who:       body.Who,
		Tailscale: body.Tailscale,
		Bluetooth: body.Bluetooth,
		CreatedOn: createdOn,
	}
	return
}

func listSitrepFields(hostname string) sitrep.Sitrep {
	return sitrep.Sitrep{
		ID:        uuid.NewString(),
		Hostname:  hostname,
		AptLog:    TestSitrepAptLog,
		Last:      TestSitrepLast,
		WAP:       TestSitrepWAP,
		Free:      TestSitrepFree,
		DF:        TestSitrepDF,
		Who:       TestSitrepWho,
		Tailscale: TestSitrepTailscale,
		Bluetooth: TestSitrepBluetooth,
	}
}

// ListSitreps returns three distinct list fixtures.
// When withTimestamps is true, CreatedOn is 1, 2, and 3.
func ListSitreps(withTimestamps bool) (first, second, third sitrep.Sitrep) {
	first = listSitrepFields(ListSitrepHostnameFirst)
	second = listSitrepFields(ListSitrepHostnameSecond)
	third = listSitrepFields(ListSitrepHostnameThird)
	if withTimestamps {
		first.CreatedOn = 1
		second.CreatedOn = 2
		third.CreatedOn = 3
	}
	return
}

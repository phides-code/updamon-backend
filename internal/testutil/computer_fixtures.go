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

// Canonical valid hardware/OS field values for handler and DynamoDB tests.
const (
	TestComputerOS      = "Debian 13"
	TestComputerKernel  = "6.12.0"
	TestComputerModel   = "ThinkPad T14"
	TestComputerRAM     = "32GB"
	TestComputerCPU     = "AMD Ryzen 7"
	TestComputerStorage = "1TB NVMe"
)

// TestStoredComputerCreatedOn is a fixed timestamp for persisted-computer repository tests.
const TestStoredComputerCreatedOn uint64 = 12345

const (
	ListComputerHostnameFirst  = TestComputerHostname
	ListComputerHostnameSecond = "plantain"
	ListComputerHostnameThird  = "burro"

	ListComputerIPFirst  = TestComputerIP
	ListComputerIPSecond = "10.0.0.2"
	ListComputerIPThird  = "172.16.0.3"
)

// ComputerBody is the create/update JSON payload for computer HTTP tests.
// Kept separate from computer.Computer so json-tag drift fails tests instead of silently matching.
type ComputerBody struct {
	Hostname string `json:"hostname"`
	IP       string `json:"ip"`
	OS       string `json:"os"`
	Kernel   string `json:"kernel"`
	Model    string `json:"model"`
	RAM      string `json:"ram"`
	CPU      string `json:"cpu"`
	Storage  string `json:"storage"`
}

// ValidComputerBody returns a ComputerBody with canonical valid field values.
func ValidComputerBody() ComputerBody {
	return ComputerBody{
		Hostname: TestComputerHostname,
		IP:       TestComputerIP,
		OS:       TestComputerOS,
		Kernel:   TestComputerKernel,
		Model:    TestComputerModel,
		RAM:      TestComputerRAM,
		CPU:      TestComputerCPU,
		Storage:  TestComputerStorage,
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
		OS:        body.OS,
		Kernel:    body.Kernel,
		Model:     body.Model,
		RAM:       body.RAM,
		CPU:       body.CPU,
		Storage:   body.Storage,
		CreatedOn: createdOn,
	}
	return
}

func listComputerFields(hostname, ip string) computer.Computer {
	return computer.Computer{
		ID:       uuid.NewString(),
		Hostname: hostname,
		IP:       ip,
		OS:       TestComputerOS,
		Kernel:   TestComputerKernel,
		Model:    TestComputerModel,
		RAM:      TestComputerRAM,
		CPU:      TestComputerCPU,
		Storage:  TestComputerStorage,
	}
}

// ListComputers returns three distinct list fixtures.
// When withTimestamps is true, CreatedOn is 1, 2, and 3.
func ListComputers(withTimestamps bool) (first, second, third computer.Computer) {
	first = listComputerFields(ListComputerHostnameFirst, ListComputerIPFirst)
	second = listComputerFields(ListComputerHostnameSecond, ListComputerIPSecond)
	third = listComputerFields(ListComputerHostnameThird, ListComputerIPThird)
	if withTimestamps {
		first.CreatedOn = 1
		second.CreatedOn = 2
		third.CreatedOn = 3
	}
	return
}

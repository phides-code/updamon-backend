// Unit tests for computer create/update validation.
package computer_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/phides-code/updamon-backend/internal/computer"
	"github.com/phides-code/updamon-backend/internal/testutil"
)

func validCreateInput() computer.CreateInput {
	return computer.CreateInput{
		Hostname: testutil.TestComputerHostname,
		IP:       testutil.TestComputerIP,
		OS:       testutil.TestComputerOS,
		Kernel:   testutil.TestComputerKernel,
		Model:    testutil.TestComputerModel,
		RAM:      testutil.TestComputerRAM,
		CPU:      testutil.TestComputerCPU,
		Storage:  testutil.TestComputerStorage,
	}
}

func TestValidateCreateInput(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   computer.CreateInput
		wantErr bool
	}{
		{name: "valid", input: validCreateInput(), wantErr: false},
		{
			name: "empty hostname",
			input: func() computer.CreateInput {
				in := validCreateInput()
				in.Hostname = ""
				return in
			}(),
			wantErr: true,
		},
		{
			name: "empty ip",
			input: func() computer.CreateInput {
				in := validCreateInput()
				in.IP = ""
				return in
			}(),
			wantErr: true,
		},
		{
			name: "invalid ip",
			input: func() computer.CreateInput {
				in := validCreateInput()
				in.IP = "not-an-ip"
				return in
			}(),
			wantErr: true,
		},
		{
			name: "ipv6 rejected",
			input: func() computer.CreateInput {
				in := validCreateInput()
				in.IP = "2001:db8::1"
				return in
			}(),
			wantErr: true,
		},
		{
			name: "empty os",
			input: func() computer.CreateInput {
				in := validCreateInput()
				in.OS = ""
				return in
			}(),
			wantErr: true,
		},
		{
			name: "empty kernel",
			input: func() computer.CreateInput {
				in := validCreateInput()
				in.Kernel = ""
				return in
			}(),
			wantErr: true,
		},
		{
			name: "empty model",
			input: func() computer.CreateInput {
				in := validCreateInput()
				in.Model = ""
				return in
			}(),
			wantErr: true,
		},
		{
			name: "empty ram",
			input: func() computer.CreateInput {
				in := validCreateInput()
				in.RAM = ""
				return in
			}(),
			wantErr: true,
		},
		{
			name: "empty cpu",
			input: func() computer.CreateInput {
				in := validCreateInput()
				in.CPU = ""
				return in
			}(),
			wantErr: true,
		},
		{
			name: "empty storage",
			input: func() computer.CreateInput {
				in := validCreateInput()
				in.Storage = ""
				return in
			}(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			testutil.AssertWantErr(t, computer.ValidateCreateInput(tt.input), tt.wantErr)
		})
	}
}

func TestValidateUpdateInput(t *testing.T) {
	t.Parallel()

	validID := uuid.NewString()

	validUpdateInput := func() computer.UpdateInput {
		in := validCreateInput()
		return computer.UpdateInput{
			ID:       validID,
			Hostname: in.Hostname,
			IP:       in.IP,
			OS:       in.OS,
			Kernel:   in.Kernel,
			Model:    in.Model,
			RAM:      in.RAM,
			CPU:      in.CPU,
			Storage:  in.Storage,
		}
	}

	tests := []struct {
		name    string
		input   computer.UpdateInput
		wantErr bool
	}{
		{name: "valid", input: validUpdateInput(), wantErr: false},
		{
			name: "invalid id",
			input: func() computer.UpdateInput {
				in := validUpdateInput()
				in.ID = "bad"
				return in
			}(),
			wantErr: true,
		},
		{
			name: "empty hostname",
			input: func() computer.UpdateInput {
				in := validUpdateInput()
				in.Hostname = ""
				return in
			}(),
			wantErr: true,
		},
		{
			name: "empty ip",
			input: func() computer.UpdateInput {
				in := validUpdateInput()
				in.IP = ""
				return in
			}(),
			wantErr: true,
		},
		{
			name: "invalid ip",
			input: func() computer.UpdateInput {
				in := validUpdateInput()
				in.IP = "not-an-ip"
				return in
			}(),
			wantErr: true,
		},
		{
			name: "empty os",
			input: func() computer.UpdateInput {
				in := validUpdateInput()
				in.OS = ""
				return in
			}(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			testutil.AssertWantErr(t, computer.ValidateUpdateInput(tt.input), tt.wantErr)
		})
	}
}

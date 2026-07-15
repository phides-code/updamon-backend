// Unit tests for computer create/update validation.
package computer_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/phides-code/go-multi-api/internal/computer"
	"github.com/phides-code/go-multi-api/internal/testutil"
)

func TestValidateCreateInput(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   computer.CreateInput
		wantErr bool
	}{
		{
			name:    "valid",
			input:   computer.CreateInput{Hostname: testutil.TestComputerHostname, IP: testutil.TestComputerIP},
			wantErr: false,
		},
		{
			name:    "empty hostname",
			input:   computer.CreateInput{Hostname: "", IP: testutil.TestComputerIP},
			wantErr: true,
		},
		{
			name:    "invalid ip",
			input:   computer.CreateInput{Hostname: testutil.TestComputerHostname, IP: testutil.TestComputerInvalidIP},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := computer.ValidateCreateInput(tt.input)

			if tt.wantErr && err == nil {
				t.Fatal("expected error")
			}

			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestValidateUpdateInput(t *testing.T) {
	t.Parallel()

	validID := uuid.NewString()

	tests := []struct {
		name    string
		input   computer.UpdateInput
		wantErr bool
	}{
		{
			name:    "valid",
			input:   computer.UpdateInput{ID: validID, Hostname: testutil.TestComputerHostname, IP: testutil.TestComputerIP},
			wantErr: false,
		},
		{
			name:    "invalid id",
			input:   computer.UpdateInput{ID: "bad", Hostname: testutil.TestComputerHostname, IP: testutil.TestComputerIP},
			wantErr: true,
		},
		{
			name:    "empty hostname",
			input:   computer.UpdateInput{ID: validID, Hostname: "", IP: testutil.TestComputerIP},
			wantErr: true,
		},
		{
			name:    "invalid ip",
			input:   computer.UpdateInput{ID: validID, Hostname: testutil.TestComputerHostname, IP: testutil.TestComputerInvalidIP},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := computer.ValidateUpdateInput(tt.input)

			if tt.wantErr && err == nil {
				t.Fatal("expected error")
			}

			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

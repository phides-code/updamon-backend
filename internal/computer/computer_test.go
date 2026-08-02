// Unit tests for computer create/update validation.
package computer_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/phides-code/updamon-backend/internal/computer"
	"github.com/phides-code/updamon-backend/internal/domain"
	"github.com/phides-code/updamon-backend/internal/testutil"
)

func TestValidateCreateInput(t *testing.T) {
	t.Parallel()

	validCreateInput := func() computer.CreateInput {
		return computer.CreateInput{
			Hostname: testutil.TestComputerHostname,
			IP:       testutil.TestComputerIP,
			Rating:   testutil.TestComputerRating,
		}
	}

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
			name: "rating below min",
			input: func() computer.CreateInput {
				in := validCreateInput()
				in.Rating = domain.DefaultMinInt - 1
				return in
			}(),
			wantErr: true,
		},
		{
			name: "rating above max",
			input: func() computer.CreateInput {
				in := validCreateInput()
				in.Rating = domain.DefaultMaxInt + 1
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
		return computer.UpdateInput{
			ID:       validID,
			Hostname: testutil.TestComputerHostname,
			IP:       testutil.TestComputerIP,
			Rating:   testutil.TestComputerRating,
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
			name: "rating below min",
			input: func() computer.UpdateInput {
				in := validUpdateInput()
				in.Rating = domain.DefaultMinInt - 1
				return in
			}(),
			wantErr: true,
		},
		{
			name: "rating above max",
			input: func() computer.UpdateInput {
				in := validUpdateInput()
				in.Rating = domain.DefaultMaxInt + 1
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

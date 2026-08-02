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
			Descriptor: testutil.TestComputerDescriptor,
			Rating:     testutil.TestComputerRating,
		}
	}

	tests := []struct {
		name    string
		input   computer.CreateInput
		wantErr bool
	}{
		{name: "valid", input: validCreateInput(), wantErr: false},
		{
			name: "empty descriptor",
			input: func() computer.CreateInput {
				in := validCreateInput()
				in.Descriptor = ""
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
			ID:         validID,
			Descriptor: testutil.TestComputerDescriptor,
			Rating:     testutil.TestComputerRating,
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
			name: "empty descriptor",
			input: func() computer.UpdateInput {
				in := validUpdateInput()
				in.Descriptor = ""
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

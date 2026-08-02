// Unit tests for sitrep create validation.
package sitrep_test

import (
	"testing"

	"github.com/phides-code/updamon-backend/internal/sitrep"
	"github.com/phides-code/updamon-backend/internal/testutil"
)

func TestValidateCreateInput(t *testing.T) {
	t.Parallel()

	validCreateInput := func() sitrep.CreateInput {
		return sitrep.CreateInput{
			Hostname: testutil.TestSitrepHostname,
		}
	}

	tests := []struct {
		name    string
		input   sitrep.CreateInput
		wantErr bool
	}{
		{name: "valid", input: validCreateInput(), wantErr: false},
		{
			name: "empty hostname",
			input: func() sitrep.CreateInput {
				in := validCreateInput()
				in.Hostname = ""
				return in
			}(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			testutil.AssertWantErr(t, sitrep.ValidateCreateInput(tt.input), tt.wantErr)
		})
	}
}

// Unit tests for sitrep create validation.
package sitrep_test

import (
	"strings"
	"testing"

	"github.com/phides-code/updamon-backend/internal/domain"
	"github.com/phides-code/updamon-backend/internal/sitrep"
	"github.com/phides-code/updamon-backend/internal/testutil"
)

func TestValidateCreateInput(t *testing.T) {
	t.Parallel()

	validCreateInput := func() sitrep.CreateInput {
		return sitrep.CreateInput{
			Hostname:  testutil.TestSitrepHostname,
			AptLog:    testutil.TestSitrepAptLog,
			Last:      testutil.TestSitrepLast,
			WAP:       testutil.TestSitrepWAP,
			Free:      testutil.TestSitrepFree,
			DF:        testutil.TestSitrepDF,
			Who:       testutil.TestSitrepWho,
			Tailscale: testutil.TestSitrepTailscale,
			Bluetooth: testutil.TestSitrepBluetooth,
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
		{
			name: "empty aptlog",
			input: func() sitrep.CreateInput {
				in := validCreateInput()
				in.AptLog = ""
				return in
			}(),
			wantErr: true,
		},
		{
			name: "empty last",
			input: func() sitrep.CreateInput {
				in := validCreateInput()
				in.Last = ""
				return in
			}(),
			wantErr: true,
		},
		{
			name: "empty wap",
			input: func() sitrep.CreateInput {
				in := validCreateInput()
				in.WAP = ""
				return in
			}(),
			wantErr: true,
		},
		{
			name: "invalid wap",
			input: func() sitrep.CreateInput {
				in := validCreateInput()
				in.WAP = "not-a-mac"
				return in
			}(),
			wantErr: true,
		},
		{
			name: "empty free",
			input: func() sitrep.CreateInput {
				in := validCreateInput()
				in.Free = ""
				return in
			}(),
			wantErr: true,
		},
		{
			name: "empty df",
			input: func() sitrep.CreateInput {
				in := validCreateInput()
				in.DF = ""
				return in
			}(),
			wantErr: true,
		},
		{
			name: "empty who",
			input: func() sitrep.CreateInput {
				in := validCreateInput()
				in.Who = ""
				return in
			}(),
			wantErr: true,
		},
		{
			name: "empty tailscale",
			input: func() sitrep.CreateInput {
				in := validCreateInput()
				in.Tailscale = ""
				return in
			}(),
			wantErr: true,
		},
		{
			name: "empty bluetooth",
			input: func() sitrep.CreateInput {
				in := validCreateInput()
				in.Bluetooth = ""
				return in
			}(),
			wantErr: true,
		},
		{
			name: "aptlog at max",
			input: func() sitrep.CreateInput {
				in := validCreateInput()
				in.AptLog = strings.Repeat("a", domain.DefaultMaxLongStringLength)
				return in
			}(),
			wantErr: false,
		},
		{
			name: "aptlog too long",
			input: func() sitrep.CreateInput {
				in := validCreateInput()
				in.AptLog = strings.Repeat("a", domain.DefaultMaxLongStringLength+1)
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

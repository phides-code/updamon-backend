// Unit tests for shared string, IPv4, and MAC validation helpers.
package domain_test

import (
	"strings"
	"testing"

	"github.com/phides-code/updamon-backend/internal/domain"
	"github.com/phides-code/updamon-backend/internal/testutil"
)

func TestValidateRequiredString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{name: "valid", value: "hello", wantErr: false},
		{name: "empty", value: "", wantErr: true},
		{name: "whitespace", value: "   ", wantErr: true},
		{name: "max length", value: strings.Repeat("a", domain.DefaultMaxStringLength), wantErr: false},
		{name: "too long", value: strings.Repeat("a", domain.DefaultMaxStringLength+1), wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := domain.ValidateRequiredString(tt.value, domain.DefaultMinStringLength, domain.DefaultMaxStringLength)
			testutil.AssertWantErr(t, err, tt.wantErr)
		})
	}
}

func TestValidateIPv4(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{name: "valid", value: "192.168.1.10", wantErr: false},
		{name: "empty", value: "", wantErr: true},
		{name: "whitespace", value: "   ", wantErr: true},
		{name: "hostname", value: "example.com", wantErr: true},
		{name: "ipv6", value: "2001:db8::1", wantErr: true},
		{name: "ipv4 mapped ipv6", value: "::ffff:192.0.2.1", wantErr: true},
		{name: "octet too large", value: "256.0.0.1", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			testutil.AssertWantErr(t, domain.ValidateIPv4(tt.value), tt.wantErr)
		})
	}
}

func TestValidateMAC(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{name: "colon separated", value: "cc:f4:11:32:cb:ff", wantErr: false},
		{name: "hyphen separated", value: "cc-f4-11-32-cb-ff", wantErr: false},
		{name: "dot separated", value: "ccf4.1132.cbff", wantErr: false},
		{name: "empty", value: "", wantErr: true},
		{name: "whitespace", value: "   ", wantErr: true},
		{name: "not a mac", value: "not-a-mac", wantErr: true},
		{name: "eui64 rejected", value: "00:11:22:33:44:55:66:77", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			testutil.AssertWantErr(t, domain.ValidateMAC(tt.value), tt.wantErr)
		})
	}
}

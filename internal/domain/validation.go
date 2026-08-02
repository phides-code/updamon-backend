// Shared validation helpers for create/update payloads across resources.
package domain

import (
	"net"
	"strings"
	"unicode/utf8"
)

// Default rune length bounds for required strings unless a resource opts out.
const (
	DefaultMinStringLength = 1
	DefaultMaxStringLength = 100
)

// Default bounds for required integers unless a resource opts out.
const (
	DefaultMinInt = 0
	DefaultMaxInt = 100
)

// ValidateRequiredString rejects blank values (after trim) and enforces rune length bounds.
func ValidateRequiredString(s string, minLen, maxLen int) error {
	if strings.TrimSpace(s) == "" {
		return ErrValidationFailed
	}
	length := utf8.RuneCountInString(s)
	if length < minLen || length > maxLen {
		return ErrValidationFailed
	}
	return nil
}

// ValidateRequiredInt enforces inclusive integer bounds.
func ValidateRequiredInt(n, min, max int) error {
	if n < min || n > max {
		return ErrValidationFailed
	}
	return nil
}

// ValidateIPv4 rejects blank values and anything that is not a dotted-quad IPv4 address.
func ValidateIPv4(s string) error {
	s = strings.TrimSpace(s)
	if s == "" {
		return ErrValidationFailed
	}
	// Reject IPv4-mapped IPv6 forms that ParseIP would otherwise accept via To4.
	if strings.Contains(s, ":") {
		return ErrValidationFailed
	}
	ip := net.ParseIP(s)
	if ip == nil || ip.To4() == nil {
		return ErrValidationFailed
	}
	return nil
}

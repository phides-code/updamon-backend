// Admin key validation via the X-Admin-Key request header.
package platform

import (
	"crypto/subtle"
	"os"
)

const AdminKeyHeader = "X-Admin-Key"

// AdminKeyEnvVar is the process env key for the expected Admin key (SAM maps AdminKey here).
const AdminKeyEnvVar = "ADMIN_KEY"

func ExpectedAdminKey() string {
	return os.Getenv(AdminKeyEnvVar)
}

func ValidAdminKey(expected string, headers map[string]string) bool {
	if expected == "" {
		return false
	}

	provided := HeaderValue(headers, AdminKeyHeader)
	if provided == "" {
		return false
	}

	return subtle.ConstantTimeCompare([]byte(expected), []byte(provided)) == 1
}

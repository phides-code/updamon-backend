// Unit tests for X-Admin-Key header validation.
package platform_test

import (
	"testing"

	"github.com/phides-code/updamon-backend/internal/platform"
)

func TestValidAdminKey(t *testing.T) {
	t.Parallel()

	headers := map[string]string{platform.AdminKeyHeader: "secret"}

	if !platform.ValidAdminKey("secret", headers) {
		t.Fatal("expected key to match")
	}
	if platform.ValidAdminKey("wrong", headers) {
		t.Fatal("expected key mismatch")
	}
	if platform.ValidAdminKey("", headers) {
		t.Fatal("expected empty expected key to fail")
	}
	if platform.ValidAdminKey("secret", map[string]string{}) {
		t.Fatal("expected missing header to fail")
	}
}

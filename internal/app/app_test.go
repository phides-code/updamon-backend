// Composition smoke tests: built gateway serves registered routes without panicking.
package app

import (
	"context"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/phides-code/updamon-backend/internal/computer"
	"github.com/phides-code/updamon-backend/internal/gateway"
	"github.com/phides-code/updamon-backend/internal/platform"
	"github.com/phides-code/updamon-backend/internal/sitrep"
	"github.com/phides-code/updamon-backend/internal/testutil"
)

func TestWiringSmokeGETComputers(t *testing.T) {
	assertWiringSmokeGET(t, testGateway(t), "/"+computer.PathPrefix)
}

func TestWiringSmokeGETSitreps(t *testing.T) {
	assertWiringSmokeGET(t, testGateway(t), "/"+sitrep.PathPrefix)
}

func testGateway(t *testing.T) *gateway.Gateway {
	t.Helper()
	t.Setenv(platform.CFTTokenEnvVar, testutil.TestCFTToken)
	return buildGateway(platform.NewLogger(), stubComputerRepo{}, stubSitrepRepo{})
}

func assertWiringSmokeGET(t *testing.T, g *gateway.Gateway, path string) {
	t.Helper()

	resp, err := g.Handle(context.Background(), events.APIGatewayProxyRequest{
		HTTPMethod: http.MethodGet,
		Path:       path,
		Headers:    testutil.CFTokenHeaders(testutil.TestCFTToken),
	})
	if err != nil {
		t.Fatalf("handle: %v", err)
	}
	if resp.StatusCode >= http.StatusInternalServerError {
		t.Fatalf("status = %d, want < 500", resp.StatusCode)
	}
}

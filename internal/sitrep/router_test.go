// Gateway integration tests for the sitreps resource.
package sitrep_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
	"github.com/phides-code/updamon-backend/internal/sitrep"
	"github.com/phides-code/updamon-backend/internal/domain"
	"github.com/phides-code/updamon-backend/internal/gateway"
	"github.com/phides-code/updamon-backend/internal/platform"
	"github.com/phides-code/updamon-backend/internal/testutil"
)

func registeredSitrepGateway(repo sitrep.Repository) *gateway.Gateway {
	g := gateway.NewGatewayWithCFTToken(platform.NewLogger(), testutil.TestCFTToken)
	g.Register(sitrep.PathPrefix, sitrep.NewHandler(repo, platform.NewLogger()))
	return g
}

func TestGatewayRoutesSitreps(t *testing.T) {
	t.Parallel()

	id := uuid.NewString()
	g := registeredSitrepGateway(dispatchSitrepRepo())

	resp, err := g.Handle(context.Background(), events.APIGatewayProxyRequest{
		HTTPMethod:     http.MethodGet,
		Path:           "/" + sitrep.PathPrefix + "/" + id,
		PathParameters: map[string]string{sitrep.AttrID: id},
		Headers:        testutil.CFTokenHeaders(testutil.TestCFTToken),
	})
	testutil.RequireHandle(t, resp, err, http.StatusOK)
}

func TestGatewaySkipsCFTTokenUnderSAMLocal(t *testing.T) {
	t.Setenv(platform.SAMLocalEnvVar, "true")

	id := uuid.NewString()
	g := registeredSitrepGateway(dispatchSitrepRepo())

	resp, err := g.Handle(context.Background(), events.APIGatewayProxyRequest{
		HTTPMethod:     http.MethodGet,
		Path:           "/" + sitrep.PathPrefix + "/" + id,
		PathParameters: map[string]string{sitrep.AttrID: id},
	})
	testutil.RequireHandle(t, resp, err, http.StatusOK)
}

func TestGatewayAllowsOptionsWithoutCFTToken(t *testing.T) {
	t.Parallel()

	g := registeredSitrepGateway(emptySitrepRepo())

	resp, err := g.Handle(context.Background(), events.APIGatewayProxyRequest{
		HTTPMethod: http.MethodOptions,
		Path:       "/" + sitrep.PathPrefix,
	})
	envelope := testutil.RequireHandle(t, resp, err, http.StatusMethodNotAllowed)
	testutil.AssertAPIError(t, envelope, domain.ErrMethodNotAllowed.Error())
}

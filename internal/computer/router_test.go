// Gateway integration tests for the computers resource.
package computer_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
	"github.com/phides-code/updamon-backend/internal/computer"
	"github.com/phides-code/updamon-backend/internal/domain"
	"github.com/phides-code/updamon-backend/internal/gateway"
	"github.com/phides-code/updamon-backend/internal/platform"
	"github.com/phides-code/updamon-backend/internal/testutil"
)

func registeredComputerGateway(repo computer.Repository) *gateway.Gateway {
	g := gateway.NewGatewayWithCFTToken(platform.NewLogger(), testutil.TestCFTToken)
	g.Register(computer.PathPrefix, computer.NewHandler(repo, platform.NewLogger()))
	return g
}

func TestGatewayRoutesComputers(t *testing.T) {
	t.Parallel()

	id := uuid.NewString()
	g := registeredComputerGateway(dispatchComputerRepo())

	resp, err := g.Handle(context.Background(), events.APIGatewayProxyRequest{
		HTTPMethod:     http.MethodGet,
		Path:           "/" + computer.PathPrefix + "/" + id,
		PathParameters: map[string]string{computer.AttrID: id},
		Headers:        testutil.CFTokenHeaders(testutil.TestCFTToken),
	})
	testutil.RequireHandle(t, resp, err, http.StatusOK)
}

func TestGatewaySkipsCFTTokenUnderSAMLocal(t *testing.T) {
	t.Setenv(platform.SAMLocalEnvVar, "true")

	id := uuid.NewString()
	g := registeredComputerGateway(dispatchComputerRepo())

	resp, err := g.Handle(context.Background(), events.APIGatewayProxyRequest{
		HTTPMethod:     http.MethodGet,
		Path:           "/" + computer.PathPrefix + "/" + id,
		PathParameters: map[string]string{computer.AttrID: id},
	})
	testutil.RequireHandle(t, resp, err, http.StatusOK)
}

func TestGatewayAllowsOptionsWithoutCFTToken(t *testing.T) {
	t.Parallel()

	g := registeredComputerGateway(emptyComputerRepo())

	resp, err := g.Handle(context.Background(), events.APIGatewayProxyRequest{
		HTTPMethod: http.MethodOptions,
		Path:       "/" + computer.PathPrefix,
	})
	envelope := testutil.RequireHandle(t, resp, err, http.StatusMethodNotAllowed)
	testutil.AssertAPIError(t, envelope, domain.ErrMethodNotAllowed.Error())
}

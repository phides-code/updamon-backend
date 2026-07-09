// Composition smoke tests: verify the built gateway handles computer routes without panicking.
package app

import (
	"context"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/phides-code/go-multi-api/internal/computer"
	"github.com/phides-code/go-multi-api/internal/platform"
	"github.com/phides-code/go-multi-api/internal/testutil"
)

type stubComputerRepo struct{}

func (stubComputerRepo) Create(_ context.Context, _ computer.Computer) (computer.Computer, error) {
	return computer.Computer{}, nil
}
func (stubComputerRepo) GetByID(_ context.Context, _ string) (computer.Computer, error) {
	return computer.Computer{}, nil
}
func (stubComputerRepo) List(_ context.Context) ([]computer.Computer, error) {
	return nil, nil
}
func (stubComputerRepo) Update(_ context.Context, _ computer.Computer) (computer.Computer, error) {
	return computer.Computer{}, nil
}
func (stubComputerRepo) Delete(_ context.Context, _ string) (computer.Computer, error) {
	return computer.Computer{}, nil
}

func TestWiringSmokeGETComputers(t *testing.T) {
	t.Setenv("AWS_CF_TOKEN", testutil.TestCFTToken)

	g := buildGateway(platform.NewLogger(), stubComputerRepo{})

	resp, err := g.Handle(context.Background(), events.APIGatewayProxyRequest{
		HTTPMethod: http.MethodGet,
		Path:       "/computers",
		Headers:    map[string]string{platform.CFTTokenHeader: testutil.TestCFTToken},
	})
	if err != nil {
		t.Fatalf("handle: %v", err)
	}
	if resp.StatusCode >= http.StatusInternalServerError {
		t.Fatalf("status = %d, want < 500", resp.StatusCode)
	}
}

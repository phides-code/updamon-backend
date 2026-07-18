// Unit tests for computer HTTP handling using a mocked repository.
package computer_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/phides-code/go-multi-api/internal/computer"
	"github.com/phides-code/go-multi-api/internal/domain"
	"github.com/phides-code/go-multi-api/internal/platform"
	"github.com/phides-code/go-multi-api/internal/testutil"
)

func TestComputerHandlerCreate(t *testing.T) {
	t.Parallel()

	validCreateBody := testutil.ComputerCreateBody(testutil.TestComputerHostname, testutil.TestComputerIP, testutil.TestComputerOS)

	tests := []struct {
		name         string
		body         string
		setupRepo    func() *mockComputerRepository
		wantStatus   int
		wantErrorMsg string
		wantHostname string
		wantIP       string
		wantOS       string
	}{
		{
			name: "success",
			body: validCreateBody,
			setupRepo: func() *mockComputerRepository {
				return &mockComputerRepository{
					createFn: func(_ context.Context, computer computer.Computer) (computer.Computer, error) {
						return computer, nil
					},
				}
			},
			wantStatus:   http.StatusCreated,
			wantHostname: testutil.TestComputerHostname,
			wantIP:       testutil.TestComputerIP,
			wantOS:       testutil.TestComputerOS,
		},
		{
			name: "repo failure",
			body: validCreateBody,
			setupRepo: func() *mockComputerRepository {
				return &mockComputerRepository{
					createFn: func(_ context.Context, _ computer.Computer) (computer.Computer, error) {
						return computer.Computer{}, errors.New("db down")
					},
				}
			},
			wantStatus:   http.StatusInternalServerError,
			wantErrorMsg: platform.InternalServerErrorMessage,
		},
		{
			name: "duplicate id",
			body: validCreateBody,
			setupRepo: func() *mockComputerRepository {
				return &mockComputerRepository{
					createFn: func(_ context.Context, _ computer.Computer) (computer.Computer, error) {
						return computer.Computer{}, domain.ErrAlreadyExists
					},
				}
			},
			wantStatus:   http.StatusConflict,
			wantErrorMsg: "already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			h := computer.NewHandler(tt.setupRepo(), platform.NewLogger())

			resp, err := h.Handle(context.Background(), events.APIGatewayProxyRequest{
				HTTPMethod: "POST",
				Body:       tt.body,
			})
			if err != nil {
				t.Fatalf("handle: %v", err)
			}

			envelope := testutil.RequireStatusAndEnvelope(t, resp, tt.wantStatus)

			if tt.wantErrorMsg != "" {
				testutil.AssertAPIError(t, envelope, tt.wantErrorMsg)
				return
			}

			computer := decodeComputerData(t, envelope)
			assertComputerDataKeys(t, envelope)

			if computer.Hostname != tt.wantHostname {
				t.Fatalf("hostname = %q, want %q", computer.Hostname, tt.wantHostname)
			}
			if computer.IP != tt.wantIP {
				t.Fatalf("ip = %q, want %q", computer.IP, tt.wantIP)
			}
			if computer.OS != tt.wantOS {
				t.Fatalf("os = %q, want %q", computer.OS, tt.wantOS)
			}

			if err := domain.ValidateID(computer.ID); err != nil {
				t.Fatalf("expected generated uuid: %v", err)
			}
			if computer.CreatedOn == 0 {
				t.Fatal("expected createdOn in response")
			}
			now := uint64(time.Now().UnixMilli())
			if computer.CreatedOn > now || now-computer.CreatedOn > 5000 {
				t.Fatalf("createdOn = %d, expected within 5s of %d", computer.CreatedOn, now)
			}
		})
	}
}

func TestComputerHandlerDelete(t *testing.T) {
	t.Parallel()

	validUuid, deletedComputer, _ := existingComputerFixture()

	tests := []struct {
		name         string
		pathID       string
		wantStatus   int
		wantComputer *computer.Computer
		wantErrorMsg string
		setupRepo    func(pathID string) *mockComputerRepository
	}{
		{
			name:         "DELETE success",
			pathID:       validUuid,
			wantStatus:   http.StatusOK,
			wantComputer: &deletedComputer,
			wantErrorMsg: "",
			setupRepo: func(pathID string) *mockComputerRepository {
				return &mockComputerRepository{
					deleteFn: func(_ context.Context, id string) (computer.Computer, error) {
						if id != pathID {
							return computer.Computer{}, domain.ErrNotFound
						}
						return deletedComputer, nil
					},
				}
			},
		},
		{
			name:         "DELETE invalid ID",
			pathID:       "bad id",
			wantStatus:   http.StatusBadRequest,
			wantComputer: nil,
			wantErrorMsg: "invalid id",
			setupRepo:    func(pathID string) *mockComputerRepository { return emptyComputerRepo() },
		},
		{
			name:         "DELETE ID not found",
			pathID:       validUuid,
			wantStatus:   http.StatusNotFound,
			wantComputer: nil,
			wantErrorMsg: "not found",
			setupRepo: func(pathID string) *mockComputerRepository {
				return &mockComputerRepository{
					deleteFn: func(_ context.Context, id string) (computer.Computer, error) {
						if id == pathID {
							return computer.Computer{}, domain.ErrNotFound
						}
						return deletedComputer, nil
					},
				}
			},
		},
		{
			name:         "DELETE repo failure",
			pathID:       validUuid,
			wantStatus:   http.StatusInternalServerError,
			wantComputer: nil,
			wantErrorMsg: platform.InternalServerErrorMessage,
			setupRepo: func(pathID string) *mockComputerRepository {
				return &mockComputerRepository{
					deleteFn: func(_ context.Context, _ string) (computer.Computer, error) {
						return computer.Computer{}, errors.New("db down")
					},
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			h := computer.NewHandler(tt.setupRepo(tt.pathID), platform.NewLogger())

			req := events.APIGatewayProxyRequest{
				HTTPMethod: http.MethodDelete,
			}

			if tt.pathID != "" {
				req.PathParameters = map[string]string{"id": tt.pathID}
			}

			resp, err := h.Handle(context.Background(), req)
			if err != nil {
				t.Fatalf("handle: %v", err)
			}

			envelope := testutil.RequireStatusAndEnvelope(t, resp, tt.wantStatus)

			if tt.wantErrorMsg != "" {
				testutil.AssertAPIError(t, envelope, tt.wantErrorMsg)
				return
			}

			computer := decodeComputerData(t, envelope)

			if computer != *tt.wantComputer {
				t.Fatalf("computer = %+v, want %+v", computer, tt.wantComputer)
			}
		})
	}
}

func TestComputerHandlerGetByID(t *testing.T) {
	t.Parallel()

	validUuid, validComputer, _ := existingComputerFixture()

	tests := []struct {
		name         string
		pathID       string
		wantStatus   int
		wantComputer *computer.Computer
		wantErrorMsg string
		setupRepo    func(pathID string) *mockComputerRepository
	}{
		{
			name:         "GET by ID success",
			pathID:       validUuid,
			wantStatus:   http.StatusOK,
			wantComputer: &validComputer,
			wantErrorMsg: "",
			setupRepo: func(pathID string) *mockComputerRepository {
				return &mockComputerRepository{
					getFn: func(_ context.Context, id string) (computer.Computer, error) {
						if id != pathID {
							return computer.Computer{}, domain.ErrNotFound
						}
						return validComputer, nil
					},
				}
			},
		},
		{
			name:         "GET by ID invalid",
			pathID:       "bad id",
			wantStatus:   http.StatusBadRequest,
			wantComputer: nil,
			wantErrorMsg: "invalid id",
			setupRepo:    func(pathID string) *mockComputerRepository { return emptyComputerRepo() },
		},
		{
			name:         "GET by ID not found",
			pathID:       validUuid,
			wantStatus:   http.StatusNotFound,
			wantComputer: nil,
			wantErrorMsg: "not found",
			setupRepo: func(pathID string) *mockComputerRepository {
				return &mockComputerRepository{
					getFn: func(_ context.Context, id string) (computer.Computer, error) {
						if id == pathID {
							return computer.Computer{}, domain.ErrNotFound
						}
						return validComputer, nil
					},
				}
			},
		},
		{
			name:         "GET by ID repo failure",
			pathID:       validUuid,
			wantStatus:   http.StatusInternalServerError,
			wantComputer: nil,
			wantErrorMsg: platform.InternalServerErrorMessage,
			setupRepo: func(pathID string) *mockComputerRepository {
				return &mockComputerRepository{
					getFn: func(_ context.Context, _ string) (computer.Computer, error) {
						return computer.Computer{}, errors.New("db down")
					},
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			h := computer.NewHandler(tt.setupRepo(tt.pathID), platform.NewLogger())

			req := events.APIGatewayProxyRequest{
				HTTPMethod: http.MethodGet,
			}

			if tt.pathID != "" {
				req.PathParameters = map[string]string{"id": tt.pathID}
			}

			resp, err := h.Handle(context.Background(), req)
			if err != nil {
				t.Fatalf("handle: %v", err)
			}
			envelope := testutil.RequireStatusAndEnvelope(t, resp, tt.wantStatus)

			if tt.wantErrorMsg != "" {
				testutil.AssertAPIError(t, envelope, tt.wantErrorMsg)
				return
			}

			computer := decodeComputerData(t, envelope)
			assertComputerDataKeys(t, envelope)

			if computer != *tt.wantComputer {
				t.Fatalf("computer = %+v, want %+v", computer, tt.wantComputer)
			}
		})
	}
}

func TestComputerHandlerClientErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		method       string
		body         string
		wantStatus   int
		wantErrorMsg string
		setupRepo    func() *mockComputerRepository
	}{
		{
			name:         "POST invalid json",
			method:       "POST",
			body:         "{not json",
			wantStatus:   http.StatusBadRequest,
			wantErrorMsg: "invalid json",
		},
		{
			name:         "POST empty hostname",
			method:       "POST",
			body:         "{\"hostname\":\"\"}",
			wantStatus:   http.StatusBadRequest,
			wantErrorMsg: "validation failed",
			setupRepo:    panicComputerRepo,
		},
		{
			name:         "PATCH unsupported method",
			method:       "PATCH",
			body:         "",
			wantStatus:   http.StatusMethodNotAllowed,
			wantErrorMsg: "method not allowed",
		},
		{
			name:         "POST whitespace hostname",
			method:       "POST",
			body:         `{"hostname":"   "}`,
			wantStatus:   http.StatusBadRequest,
			wantErrorMsg: "validation failed",
			setupRepo:    panicComputerRepo,
		},
		{
			name:         "POST hostname too long",
			method:       "POST",
			body:         fmt.Sprintf(`{"hostname":%q}`, strings.Repeat("a", domain.DefaultMaxStringLength+1)),
			wantStatus:   http.StatusBadRequest,
			wantErrorMsg: "validation failed",
			setupRepo:    panicComputerRepo,
		},
		{
			name:         "POST invalid ip",
			method:       "POST",
			body:         fmt.Sprintf(`{"hostname":%q,"ip":%q,"os":%q}`, testutil.TestComputerHostname, testutil.TestComputerInvalidIP, testutil.TestComputerOS),
			wantStatus:   http.StatusBadRequest,
			wantErrorMsg: "validation failed",
			setupRepo:    panicComputerRepo,
		},
		{
			name:         "POST empty ip",
			method:       "POST",
			body:         fmt.Sprintf(`{"hostname":%q,"ip":%q,"os":%q}`, testutil.TestComputerHostname, "", testutil.TestComputerOS),
			wantStatus:   http.StatusBadRequest,
			wantErrorMsg: "validation failed",
			setupRepo:    panicComputerRepo,
		},
		{
			name:         "POST empty os",
			method:       "POST",
			body:         fmt.Sprintf(`{"hostname":%q,"ip":%q,"os":%q}`, testutil.TestComputerHostname, testutil.TestComputerIP, ""),
			wantStatus:   http.StatusBadRequest,
			wantErrorMsg: "validation failed",
			setupRepo:    panicComputerRepo,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := emptyComputerRepo()
			if tt.setupRepo != nil {
				repo = tt.setupRepo()
			}

			h := computer.NewHandler(repo, platform.NewLogger())

			req := events.APIGatewayProxyRequest{
				HTTPMethod: tt.method,
				Body:       tt.body,
			}

			resp, err := h.Handle(context.Background(), req)
			if err != nil {
				t.Fatalf("handle: %v", err)
			}

			testutil.AssertAPIError(t, testutil.RequireStatusAndEnvelope(t, resp, tt.wantStatus), tt.wantErrorMsg)
		})
	}
}

func TestComputerHandlerList(t *testing.T) {
	t.Parallel()

	computerOne, computerTwo, _ := testutil.ListComputers(false)
	wantItems := []computer.Computer{computerOne, computerTwo}

	tests := []struct {
		name         string
		wantStatus   int
		wantItems    []computer.Computer
		wantErrorMsg string
		setupRepo    func() *mockComputerRepository
	}{
		{
			name:       "GET list returns items",
			wantStatus: http.StatusOK,
			wantItems:  wantItems,
			setupRepo:  func() *mockComputerRepository { return listComputerRepo(wantItems) },
		},
		{
			name:       "GET list empty",
			wantStatus: http.StatusOK,
			wantItems:  []computer.Computer{},
			setupRepo:  func() *mockComputerRepository { return listComputerRepo([]computer.Computer{}) },
		},
		{
			name:         "GET list repo failure",
			wantStatus:   http.StatusInternalServerError,
			wantErrorMsg: platform.InternalServerErrorMessage,
			setupRepo: func() *mockComputerRepository {
				return &mockComputerRepository{
					listFn: func(_ context.Context) ([]computer.Computer, error) {
						return nil, errors.New("db down")
					},
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			h := computer.NewHandler(tt.setupRepo(), platform.NewLogger())

			resp, err := h.Handle(context.Background(), events.APIGatewayProxyRequest{
				HTTPMethod: http.MethodGet,
			})
			if err != nil {
				t.Fatalf("handle: %v", err)
			}

			envelope := testutil.RequireStatusAndEnvelope(t, resp, tt.wantStatus)

			if tt.wantErrorMsg != "" {
				testutil.AssertAPIError(t, envelope, tt.wantErrorMsg)
				return
			}

			items := decodeComputerListData(t, envelope)

			if len(items) != len(tt.wantItems) {
				t.Fatalf("len(items) = %d, want %d", len(items), len(tt.wantItems))
			}

			for i := range tt.wantItems {
				if items[i] != tt.wantItems[i] {
					t.Fatalf("items[%d] = %+v, want %+v", i, items[i], tt.wantItems[i])
				}
			}
		})
	}
}

func TestComputerHandlerUpdate(t *testing.T) {
	t.Parallel()

	validUuid, updatedComputer, validUpdateBody := existingComputerFixture()

	tests := []struct {
		name         string
		pathID       string
		body         string
		wantStatus   int
		wantComputer *computer.Computer
		wantErrorMsg string
		setupRepo    func(pathID string) *mockComputerRepository
	}{
		{
			name:         "PUT success",
			pathID:       validUuid,
			body:         validUpdateBody,
			wantStatus:   http.StatusOK,
			wantComputer: &updatedComputer,
			wantErrorMsg: "",
			setupRepo: func(pathID string) *mockComputerRepository {
				return updateComputerRepo(pathID, updatedComputer)
			},
		},
		{
			name:         "PUT invalid ID",
			pathID:       "bad id",
			body:         validUpdateBody,
			wantStatus:   http.StatusBadRequest,
			wantComputer: nil,
			wantErrorMsg: "invalid id",
			setupRepo: func(pathID string) *mockComputerRepository {
				return emptyComputerRepo()
			},
		},
		{
			name:         "PUT invalid JSON",
			pathID:       validUuid,
			body:         "not json",
			wantStatus:   http.StatusBadRequest,
			wantComputer: nil,
			wantErrorMsg: "invalid json",
			setupRepo: func(pathID string) *mockComputerRepository {
				return emptyComputerRepo()
			},
		},
		{
			name:         "PUT empty hostname",
			pathID:       validUuid,
			body:         `{"hostname":""}`,
			wantStatus:   http.StatusBadRequest,
			wantComputer: nil,
			wantErrorMsg: "validation failed",
			setupRepo:    func(pathID string) *mockComputerRepository { return emptyComputerRepo() },
		},
		{
			name:         "PUT computer not found",
			pathID:       validUuid,
			wantStatus:   http.StatusNotFound,
			body:         validUpdateBody,
			wantComputer: nil,
			wantErrorMsg: "not found",
			setupRepo: func(pathID string) *mockComputerRepository {
				return &mockComputerRepository{
					updateFn: func(_ context.Context, b computer.Computer) (computer.Computer, error) {
						if b.ID == pathID {
							return computer.Computer{}, domain.ErrNotFound
						}
						return updatedComputer, nil
					},
				}
			},
		},
		{
			name:         "PUT repo failure",
			pathID:       validUuid,
			body:         validUpdateBody,
			wantStatus:   http.StatusInternalServerError,
			wantComputer: nil,
			wantErrorMsg: platform.InternalServerErrorMessage,
			setupRepo: func(pathID string) *mockComputerRepository {
				return &mockComputerRepository{
					updateFn: func(_ context.Context, _ computer.Computer) (computer.Computer, error) {
						return computer.Computer{}, errors.New("db down")
					},
				}
			},
		},
		{
			name:         "PUT whitespace hostname",
			pathID:       validUuid,
			body:         `{"hostname":"   "}`,
			wantStatus:   http.StatusBadRequest,
			wantComputer: nil,
			wantErrorMsg: "validation failed",
			setupRepo:    func(pathID string) *mockComputerRepository { return emptyComputerRepo() },
		},
		{
			name:         "PUT hostname too long",
			pathID:       validUuid,
			body:         fmt.Sprintf(`{"hostname":%q}`, strings.Repeat("a", domain.DefaultMaxStringLength+1)),
			wantStatus:   http.StatusBadRequest,
			wantComputer: nil,
			wantErrorMsg: "validation failed",
			setupRepo:    func(pathID string) *mockComputerRepository { return emptyComputerRepo() },
		},
		{
			name:         "PUT invalid ip",
			pathID:       validUuid,
			body:         fmt.Sprintf(`{"hostname":%q,"ip":%q,"os":%q}`, testutil.TestComputerHostname, testutil.TestComputerInvalidIP, testutil.TestComputerOS),
			wantStatus:   http.StatusBadRequest,
			wantComputer: nil,
			wantErrorMsg: "validation failed",
			setupRepo:    func(pathID string) *mockComputerRepository { return emptyComputerRepo() },
		},
		{
			name:         "PUT empty ip",
			pathID:       validUuid,
			body:         fmt.Sprintf(`{"hostname":%q,"ip":%q,"os":%q}`, testutil.TestComputerHostname, "", testutil.TestComputerOS),
			wantStatus:   http.StatusBadRequest,
			wantComputer: nil,
			wantErrorMsg: "validation failed",
			setupRepo:    func(pathID string) *mockComputerRepository { return emptyComputerRepo() },
		},
		{
			name:         "PUT empty os",
			pathID:       validUuid,
			body:         fmt.Sprintf(`{"hostname":%q,"ip":%q,"os":%q}`, testutil.TestComputerHostname, testutil.TestComputerIP, ""),
			wantStatus:   http.StatusBadRequest,
			wantComputer: nil,
			wantErrorMsg: "validation failed",
			setupRepo:    func(pathID string) *mockComputerRepository { return emptyComputerRepo() },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			h := computer.NewHandler(tt.setupRepo(tt.pathID), platform.NewLogger())

			req := events.APIGatewayProxyRequest{
				HTTPMethod: http.MethodPut,
				Body:       tt.body,
			}

			if tt.pathID != "" {
				req.PathParameters = map[string]string{"id": tt.pathID}
			}

			resp, err := h.Handle(context.Background(), req)
			if err != nil {
				t.Fatalf("handle: %v", err)
			}

			envelope := testutil.RequireStatusAndEnvelope(t, resp, tt.wantStatus)

			if tt.wantErrorMsg != "" {
				testutil.AssertAPIError(t, envelope, tt.wantErrorMsg)
				return
			}

			computer := decodeComputerData(t, envelope)
			assertComputerDataKeys(t, envelope)

			if computer != *tt.wantComputer {
				t.Fatalf("computer = %+v, want %+v", computer, tt.wantComputer)
			}
		})
	}
}

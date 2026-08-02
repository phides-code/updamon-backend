// Unit tests for computer HTTP handling using a mocked repository.
package computer_test

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/phides-code/updamon-backend/internal/computer"
	"github.com/phides-code/updamon-backend/internal/domain"
	"github.com/phides-code/updamon-backend/internal/platform"
	"github.com/phides-code/updamon-backend/internal/testutil"
)

func TestComputerHandlerCreate(t *testing.T) {
	t.Parallel()

	validCreateBody := testutil.ValidComputerBody().JSON(t)

	tests := []struct {
		name         string
		body         string
		setupRepo    func() *mockComputerRepository
		wantStatus   int
		wantErrorMsg string
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
			wantStatus: http.StatusCreated,
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
			wantErrorMsg: domain.ErrAlreadyExists.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			h := computer.NewHandler(tt.setupRepo(), platform.NewLogger())

			resp, err := h.Handle(context.Background(), events.APIGatewayProxyRequest{
				HTTPMethod: http.MethodPost,
				Body:       tt.body,
			})
			envelope := testutil.RequireHandle(t, resp, err, tt.wantStatus)

			if tt.wantErrorMsg != "" {
				testutil.AssertAPIError(t, envelope, tt.wantErrorMsg)
				return
			}

			got := decodeComputerData(t, envelope)
			assertComputerDataKeys(t, envelope)

			want := testutil.ValidComputerBody()
			if got.Hostname != want.Hostname {
				t.Fatalf("hostname = %q, want %q", got.Hostname, want.Hostname)
			}
			if got.IP != want.IP {
				t.Fatalf("ip = %q, want %q", got.IP, want.IP)
			}
			if got.OS != want.OS {
				t.Fatalf("os = %q, want %q", got.OS, want.OS)
			}
			if got.Kernel != want.Kernel {
				t.Fatalf("kernel = %q, want %q", got.Kernel, want.Kernel)
			}
			if got.Model != want.Model {
				t.Fatalf("model = %q, want %q", got.Model, want.Model)
			}
			if got.RAM != want.RAM {
				t.Fatalf("ram = %q, want %q", got.RAM, want.RAM)
			}
			if got.CPU != want.CPU {
				t.Fatalf("cpu = %q, want %q", got.CPU, want.CPU)
			}
			if got.Storage != want.Storage {
				t.Fatalf("storage = %q, want %q", got.Storage, want.Storage)
			}

			if err := domain.ValidateID(got.ID); err != nil {
				t.Fatalf("expected generated uuid: %v", err)
			}
			if got.CreatedOn == 0 {
				t.Fatal("expected createdOn in response")
			}
			now := uint64(time.Now().UnixMilli())
			if got.CreatedOn > now || now-got.CreatedOn > 5000 {
				t.Fatalf("createdOn = %d, expected within 5s of %d", got.CreatedOn, now)
			}
		})
	}
}

func TestComputerHandlerDelete(t *testing.T) {
	t.Parallel()

	validUuid, deletedComputer, _ := existingComputerFixture(t)

	tests := []struct {
		name         string
		pathID       string
		wantStatus   int
		wantComputer   *computer.Computer
		wantErrorMsg string
		setupRepo    func(pathID string) *mockComputerRepository
	}{
		{
			name:         "DELETE success",
			pathID:       validUuid,
			wantStatus:   http.StatusOK,
			wantComputer:   &deletedComputer,
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
			wantComputer:   nil,
			wantErrorMsg: domain.ErrInvalidID.Error(),
			setupRepo:    func(pathID string) *mockComputerRepository { return emptyComputerRepo() },
		},
		{
			name:         "DELETE ID not found",
			pathID:       validUuid,
			wantStatus:   http.StatusNotFound,
			wantComputer:   nil,
			wantErrorMsg: domain.ErrNotFound.Error(),
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
			wantComputer:   nil,
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
				req.PathParameters = map[string]string{computer.AttrID: tt.pathID}
			}

			resp, err := h.Handle(context.Background(), req)
			envelope := testutil.RequireHandle(t, resp, err, tt.wantStatus)

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

	validUuid, validComputer, _ := existingComputerFixture(t)

	tests := []struct {
		name         string
		pathID       string
		wantStatus   int
		wantComputer   *computer.Computer
		wantErrorMsg string
		setupRepo    func(pathID string) *mockComputerRepository
	}{
		{
			name:         "GET by ID success",
			pathID:       validUuid,
			wantStatus:   http.StatusOK,
			wantComputer:   &validComputer,
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
			wantComputer:   nil,
			wantErrorMsg: domain.ErrInvalidID.Error(),
			setupRepo:    func(pathID string) *mockComputerRepository { return emptyComputerRepo() },
		},
		{
			name:         "GET by ID not found",
			pathID:       validUuid,
			wantStatus:   http.StatusNotFound,
			wantComputer:   nil,
			wantErrorMsg: domain.ErrNotFound.Error(),
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
			wantComputer:   nil,
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
				req.PathParameters = map[string]string{computer.AttrID: tt.pathID}
			}

			resp, err := h.Handle(context.Background(), req)
			envelope := testutil.RequireHandle(t, resp, err, tt.wantStatus)

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

	validationBodies := newComputerValidationBodies(t)

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
			method:       http.MethodPost,
			body:         "{not json",
			wantStatus:   http.StatusBadRequest,
			wantErrorMsg: domain.ErrInvalidJSON.Error(),
		},
		{
			name:         "POST empty hostname",
			method:       http.MethodPost,
			body:         validationBodies.computerWithEmptyValue,
			wantStatus:   http.StatusBadRequest,
			wantErrorMsg: domain.ErrValidationFailed.Error(),
			setupRepo:    panicComputerRepo,
		},
		{
			name:         "PATCH unsupported method",
			method:       "PATCH",
			body:         "",
			wantStatus:   http.StatusMethodNotAllowed,
			wantErrorMsg: domain.ErrMethodNotAllowed.Error(),
		},
		{
			name:         "POST whitespace hostname",
			method:       http.MethodPost,
			body:         validationBodies.computerWithWhitespace,
			wantStatus:   http.StatusBadRequest,
			wantErrorMsg: domain.ErrValidationFailed.Error(),
			setupRepo:    panicComputerRepo,
		},
		{
			name:         "POST hostname too long",
			method:       http.MethodPost,
			body:         validationBodies.computerWithValueTooLong,
			wantStatus:   http.StatusBadRequest,
			wantErrorMsg: domain.ErrValidationFailed.Error(),
			setupRepo:    panicComputerRepo,
		},
		{
			name:         "POST invalid ip",
			method:       http.MethodPost,
			body:         validationBodies.computerWithInvalidIP,
			wantStatus:   http.StatusBadRequest,
			wantErrorMsg: domain.ErrValidationFailed.Error(),
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
			testutil.AssertAPIError(t, testutil.RequireHandle(t, resp, err, tt.wantStatus), tt.wantErrorMsg)
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
			envelope := testutil.RequireHandle(t, resp, err, tt.wantStatus)

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

	validUuid, updatedComputer, validUpdateBody := existingComputerFixture(t)
	validationBodies := newComputerValidationBodies(t)

	tests := []struct {
		name         string
		pathID       string
		body         string
		wantStatus   int
		wantComputer   *computer.Computer
		wantErrorMsg string
		setupRepo    func(pathID string) *mockComputerRepository
	}{
		{
			name:         "PUT success",
			pathID:       validUuid,
			body:         validUpdateBody,
			wantStatus:   http.StatusOK,
			wantComputer:   &updatedComputer,
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
			wantComputer:   nil,
			wantErrorMsg: domain.ErrInvalidID.Error(),
			setupRepo: func(pathID string) *mockComputerRepository {
				return emptyComputerRepo()
			},
		},
		{
			name:         "PUT invalid JSON",
			pathID:       validUuid,
			body:         "not json",
			wantStatus:   http.StatusBadRequest,
			wantComputer:   nil,
			wantErrorMsg: domain.ErrInvalidJSON.Error(),
			setupRepo: func(pathID string) *mockComputerRepository {
				return emptyComputerRepo()
			},
		},
		{
			name:         "PUT empty hostname",
			pathID:       validUuid,
			body:         validationBodies.computerWithEmptyValue,
			wantStatus:   http.StatusBadRequest,
			wantComputer:   nil,
			wantErrorMsg: domain.ErrValidationFailed.Error(),
			setupRepo:    func(pathID string) *mockComputerRepository { return emptyComputerRepo() },
		},
		{
			name:         "PUT computer not found",
			pathID:       validUuid,
			wantStatus:   http.StatusNotFound,
			body:         validUpdateBody,
			wantComputer:   nil,
			wantErrorMsg: domain.ErrNotFound.Error(),
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
			wantComputer:   nil,
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
			body:         validationBodies.computerWithWhitespace,
			wantStatus:   http.StatusBadRequest,
			wantComputer:   nil,
			wantErrorMsg: domain.ErrValidationFailed.Error(),
			setupRepo:    func(pathID string) *mockComputerRepository { return emptyComputerRepo() },
		},
		{
			name:         "PUT hostname too long",
			pathID:       validUuid,
			body:         validationBodies.computerWithValueTooLong,
			wantStatus:   http.StatusBadRequest,
			wantComputer: nil,
			wantErrorMsg: domain.ErrValidationFailed.Error(),
			setupRepo:    func(pathID string) *mockComputerRepository { return emptyComputerRepo() },
		},
		{
			name:         "PUT invalid ip",
			pathID:       validUuid,
			body:         validationBodies.computerWithInvalidIP,
			wantStatus:   http.StatusBadRequest,
			wantComputer: nil,
			wantErrorMsg: domain.ErrValidationFailed.Error(),
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
				req.PathParameters = map[string]string{computer.AttrID: tt.pathID}
			}

			resp, err := h.Handle(context.Background(), req)
			envelope := testutil.RequireHandle(t, resp, err, tt.wantStatus)

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

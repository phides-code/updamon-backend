// Unit tests for sitrep HTTP handling using a mocked repository.
package sitrep_test

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/phides-code/updamon-backend/internal/sitrep"
	"github.com/phides-code/updamon-backend/internal/domain"
	"github.com/phides-code/updamon-backend/internal/platform"
	"github.com/phides-code/updamon-backend/internal/testutil"
)

func TestSitrepHandlerCreate(t *testing.T) {
	t.Parallel()

	validCreateBody := testutil.ValidSitrepBody().JSON(t)

	tests := []struct {
		name         string
		body         string
		setupRepo    func() *mockSitrepRepository
		wantStatus   int
		wantErrorMsg string
	}{
		{
			name: "success",
			body: validCreateBody,
			setupRepo: func() *mockSitrepRepository {
				return &mockSitrepRepository{
					createFn: func(_ context.Context, sitrep sitrep.Sitrep) (sitrep.Sitrep, error) {
						return sitrep, nil
					},
				}
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "repo failure",
			body: validCreateBody,
			setupRepo: func() *mockSitrepRepository {
				return &mockSitrepRepository{
					createFn: func(_ context.Context, _ sitrep.Sitrep) (sitrep.Sitrep, error) {
						return sitrep.Sitrep{}, errors.New("db down")
					},
				}
			},
			wantStatus:   http.StatusInternalServerError,
			wantErrorMsg: platform.InternalServerErrorMessage,
		},
		{
			name: "duplicate id",
			body: validCreateBody,
			setupRepo: func() *mockSitrepRepository {
				return &mockSitrepRepository{
					createFn: func(_ context.Context, _ sitrep.Sitrep) (sitrep.Sitrep, error) {
						return sitrep.Sitrep{}, domain.ErrAlreadyExists
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

			h := sitrep.NewHandler(tt.setupRepo(), platform.NewLogger())

			resp, err := h.Handle(context.Background(), events.APIGatewayProxyRequest{
				HTTPMethod: http.MethodPost,
				Body:       tt.body,
			})
			envelope := testutil.RequireHandle(t, resp, err, tt.wantStatus)

			if tt.wantErrorMsg != "" {
				testutil.AssertAPIError(t, envelope, tt.wantErrorMsg)
				return
			}

			got := decodeSitrepData(t, envelope)
			assertSitrepDataKeys(t, envelope)

			want := testutil.ValidSitrepBody()
			if got.Hostname != want.Hostname {
				t.Fatalf("hostname = %q, want %q", got.Hostname, want.Hostname)
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


func TestSitrepHandlerGetByID(t *testing.T) {
	t.Parallel()

	validUuid, validSitrep := existingSitrepFixture(t)

	tests := []struct {
		name         string
		pathID       string
		wantStatus   int
		wantSitrep   *sitrep.Sitrep
		wantErrorMsg string
		setupRepo    func(pathID string) *mockSitrepRepository
	}{
		{
			name:         "GET by ID success",
			pathID:       validUuid,
			wantStatus:   http.StatusOK,
			wantSitrep:   &validSitrep,
			wantErrorMsg: "",
			setupRepo: func(pathID string) *mockSitrepRepository {
				return &mockSitrepRepository{
					getFn: func(_ context.Context, id string) (sitrep.Sitrep, error) {
						if id != pathID {
							return sitrep.Sitrep{}, domain.ErrNotFound
						}
						return validSitrep, nil
					},
				}
			},
		},
		{
			name:         "GET by ID invalid",
			pathID:       "bad id",
			wantStatus:   http.StatusBadRequest,
			wantSitrep:   nil,
			wantErrorMsg: domain.ErrInvalidID.Error(),
			setupRepo:    func(pathID string) *mockSitrepRepository { return emptySitrepRepo() },
		},
		{
			name:         "GET by ID not found",
			pathID:       validUuid,
			wantStatus:   http.StatusNotFound,
			wantSitrep:   nil,
			wantErrorMsg: domain.ErrNotFound.Error(),
			setupRepo: func(pathID string) *mockSitrepRepository {
				return &mockSitrepRepository{
					getFn: func(_ context.Context, id string) (sitrep.Sitrep, error) {
						if id == pathID {
							return sitrep.Sitrep{}, domain.ErrNotFound
						}
						return validSitrep, nil
					},
				}
			},
		},
		{
			name:         "GET by ID repo failure",
			pathID:       validUuid,
			wantStatus:   http.StatusInternalServerError,
			wantSitrep:   nil,
			wantErrorMsg: platform.InternalServerErrorMessage,
			setupRepo: func(pathID string) *mockSitrepRepository {
				return &mockSitrepRepository{
					getFn: func(_ context.Context, _ string) (sitrep.Sitrep, error) {
						return sitrep.Sitrep{}, errors.New("db down")
					},
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			h := sitrep.NewHandler(tt.setupRepo(tt.pathID), platform.NewLogger())

			req := events.APIGatewayProxyRequest{
				HTTPMethod: http.MethodGet,
			}

			if tt.pathID != "" {
				req.PathParameters = map[string]string{sitrep.AttrID: tt.pathID}
			}

			resp, err := h.Handle(context.Background(), req)
			envelope := testutil.RequireHandle(t, resp, err, tt.wantStatus)

			if tt.wantErrorMsg != "" {
				testutil.AssertAPIError(t, envelope, tt.wantErrorMsg)
				return
			}

			sitrep := decodeSitrepData(t, envelope)
			assertSitrepDataKeys(t, envelope)

			if sitrep != *tt.wantSitrep {
				t.Fatalf("sitrep = %+v, want %+v", sitrep, tt.wantSitrep)
			}
		})
	}
}

func TestSitrepHandlerClientErrors(t *testing.T) {
	t.Parallel()

	validationBodies := newSitrepValidationBodies(t)

	tests := []struct {
		name         string
		method       string
		body         string
		wantStatus   int
		wantErrorMsg string
		setupRepo    func() *mockSitrepRepository
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
			body:         validationBodies.sitrepWithEmptyValue,
			wantStatus:   http.StatusBadRequest,
			wantErrorMsg: domain.ErrValidationFailed.Error(),
			setupRepo:    panicSitrepRepo,
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
			body:         validationBodies.sitrepWithWhitespace,
			wantStatus:   http.StatusBadRequest,
			wantErrorMsg: domain.ErrValidationFailed.Error(),
			setupRepo:    panicSitrepRepo,
		},
		{
			name:         "POST hostname too long",
			method:       http.MethodPost,
			body:         validationBodies.sitrepWithValueTooLong,
			wantStatus:   http.StatusBadRequest,
			wantErrorMsg: domain.ErrValidationFailed.Error(),
			setupRepo:    panicSitrepRepo,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := emptySitrepRepo()
			if tt.setupRepo != nil {
				repo = tt.setupRepo()
			}

			h := sitrep.NewHandler(repo, platform.NewLogger())

			req := events.APIGatewayProxyRequest{
				HTTPMethod: tt.method,
				Body:       tt.body,
			}

			resp, err := h.Handle(context.Background(), req)
			testutil.AssertAPIError(t, testutil.RequireHandle(t, resp, err, tt.wantStatus), tt.wantErrorMsg)
		})
	}
}

func TestSitrepHandlerList(t *testing.T) {
	t.Parallel()

	sitrepOne, sitrepTwo, _ := testutil.ListSitreps(false)
	wantItems := []sitrep.Sitrep{sitrepOne, sitrepTwo}

	tests := []struct {
		name         string
		wantStatus   int
		wantItems    []sitrep.Sitrep
		wantErrorMsg string
		setupRepo    func() *mockSitrepRepository
	}{
		{
			name:       "GET list returns items",
			wantStatus: http.StatusOK,
			wantItems:  wantItems,
			setupRepo:  func() *mockSitrepRepository { return listSitrepRepo(wantItems) },
		},
		{
			name:       "GET list empty",
			wantStatus: http.StatusOK,
			wantItems:  []sitrep.Sitrep{},
			setupRepo:  func() *mockSitrepRepository { return listSitrepRepo([]sitrep.Sitrep{}) },
		},
		{
			name:         "GET list repo failure",
			wantStatus:   http.StatusInternalServerError,
			wantErrorMsg: platform.InternalServerErrorMessage,
			setupRepo: func() *mockSitrepRepository {
				return &mockSitrepRepository{
					listFn: func(_ context.Context) ([]sitrep.Sitrep, error) {
						return nil, errors.New("db down")
					},
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			h := sitrep.NewHandler(tt.setupRepo(), platform.NewLogger())

			resp, err := h.Handle(context.Background(), events.APIGatewayProxyRequest{
				HTTPMethod: http.MethodGet,
			})
			envelope := testutil.RequireHandle(t, resp, err, tt.wantStatus)

			if tt.wantErrorMsg != "" {
				testutil.AssertAPIError(t, envelope, tt.wantErrorMsg)
				return
			}

			items := decodeSitrepListData(t, envelope)

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


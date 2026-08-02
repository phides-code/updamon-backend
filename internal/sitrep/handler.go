// HTTP handler for /sitreps: maps API Gateway requests to repository operations.
package sitrep

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/phides-code/updamon-backend/internal/domain"
	"github.com/phides-code/updamon-backend/internal/platform"
)

type Handler struct {
	repo   Repository
	logger *platform.Logger
}

func NewHandler(repo Repository, logger *platform.Logger) *Handler {
	return &Handler{repo: repo, logger: logger}
}

func (h *Handler) Handle(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	const op = "sitrep request"

	id := strings.TrimSpace(req.PathParameters[AttrID])

	switch req.HTTPMethod {
	case http.MethodGet:
		if id == "" {
			return h.list(ctx, req)
		}
		return h.getByID(ctx, id)
	case http.MethodPost:
		return h.create(ctx, req.Body)
	default:
		return h.errorResponse(ctx, domain.ErrMethodNotAllowed, op)
	}
}

func (h *Handler) list(ctx context.Context, _ events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	const op = "list sitreps"

	items, err := h.repo.List(ctx)
	if err != nil {
		return h.errorResponse(ctx, err, op)
	}

	return platform.SuccessResponse(http.StatusOK, items)
}

func (h *Handler) getByID(ctx context.Context, id string) (events.APIGatewayProxyResponse, error) {
	const op = "get sitrep"

	if err := domain.ValidateID(id); err != nil {
		return h.errorResponse(ctx, err, op)
	}

	item, err := h.repo.GetByID(ctx, id)
	if err != nil {
		return h.errorResponse(ctx, err, op)
	}

	return platform.SuccessResponse(http.StatusOK, item)
}

type writePayload struct {
	Hostname string `json:"hostname"`
}

func (h *Handler) create(ctx context.Context, body string) (events.APIGatewayProxyResponse, error) {
	const op = "create sitrep"

	var payload writePayload
	if err := json.Unmarshal([]byte(body), &payload); err != nil {
		return h.errorResponse(ctx, domain.ErrInvalidJSON, op)
	}

	input := CreateInput{
		Hostname: payload.Hostname,
	}
	if err := ValidateCreateInput(input); err != nil {
		return h.errorResponse(ctx, err, op)
	}

	item := Sitrep{
		ID:        domain.NewID(),
		Hostname:  payload.Hostname,
		CreatedOn: uint64(time.Now().UnixMilli()),
	}

	created, err := h.repo.Create(ctx, item)
	if err != nil {
		return h.errorResponse(ctx, err, op)
	}

	return platform.SuccessResponse(http.StatusCreated, created)
}

func (h *Handler) errorResponse(ctx context.Context, err error, operation string) (events.APIGatewayProxyResponse, error) {
	if platform.IsClientError(err) {
		h.logger.InfoContext(ctx, operation+" client error", "error", err.Error())
	} else {
		h.logger.LogError(ctx, operation+" failed", err)
	}

	return platform.ClientErrorResponse(err)
}

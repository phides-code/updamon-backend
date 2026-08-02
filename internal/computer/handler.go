// HTTP handler for /computers: maps API Gateway requests to repository operations.
package computer

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
	const op = "computer request"

	id := strings.TrimSpace(req.PathParameters[AttrID])

	switch req.HTTPMethod {
	case http.MethodGet:
		if id == "" {
			return h.list(ctx, req)
		}
		return h.getByID(ctx, id)
	case http.MethodPost:
		return h.create(ctx, req.Body)
	case http.MethodPut:
		return h.update(ctx, id, req.Body)
	case http.MethodDelete:
		return h.delete(ctx, id)
	default:
		return h.errorResponse(ctx, domain.ErrMethodNotAllowed, op)
	}
}

func (h *Handler) list(ctx context.Context, _ events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	const op = "list computers"

	items, err := h.repo.List(ctx)
	if err != nil {
		return h.errorResponse(ctx, err, op)
	}

	return platform.SuccessResponse(http.StatusOK, items)
}

func (h *Handler) getByID(ctx context.Context, id string) (events.APIGatewayProxyResponse, error) {
	const op = "get computer"

	if err := domain.ValidateID(id); err != nil {
		return h.errorResponse(ctx, err, op)
	}

	computer, err := h.repo.GetByID(ctx, id)
	if err != nil {
		return h.errorResponse(ctx, err, op)
	}

	return platform.SuccessResponse(http.StatusOK, computer)
}

type writePayload struct {
	Hostname string `json:"hostname"`
	IP       string `json:"ip"`
	OS       string `json:"os"`
	Kernel   string `json:"kernel"`
	Model    string `json:"model"`
	RAM      string `json:"ram"`
	CPU      string `json:"cpu"`
	Storage  string `json:"storage"`
}

func (h *Handler) create(ctx context.Context, body string) (events.APIGatewayProxyResponse, error) {
	const op = "create computer"

	var payload writePayload
	if err := json.Unmarshal([]byte(body), &payload); err != nil {
		return h.errorResponse(ctx, domain.ErrInvalidJSON, op)
	}

	input := CreateInput{
		Hostname: payload.Hostname,
		IP:       payload.IP,
		OS:       payload.OS,
		Kernel:   payload.Kernel,
		Model:    payload.Model,
		RAM:      payload.RAM,
		CPU:      payload.CPU,
		Storage:  payload.Storage,
	}
	if err := ValidateCreateInput(input); err != nil {
		return h.errorResponse(ctx, err, op)
	}

	computer := Computer{
		ID:        domain.NewID(),
		Hostname:  payload.Hostname,
		IP:        payload.IP,
		OS:        payload.OS,
		Kernel:    payload.Kernel,
		Model:     payload.Model,
		RAM:       payload.RAM,
		CPU:       payload.CPU,
		Storage:   payload.Storage,
		CreatedOn: uint64(time.Now().UnixMilli()),
	}

	created, err := h.repo.Create(ctx, computer)
	if err != nil {
		return h.errorResponse(ctx, err, op)
	}

	return platform.SuccessResponse(http.StatusCreated, created)
}

func (h *Handler) update(ctx context.Context, id, body string) (events.APIGatewayProxyResponse, error) {
	const op = "update computer"

	if err := domain.ValidateID(id); err != nil {
		return h.errorResponse(ctx, err, op)
	}

	var payload writePayload
	if err := json.Unmarshal([]byte(body), &payload); err != nil {
		return h.errorResponse(ctx, domain.ErrInvalidJSON, op)
	}

	input := UpdateInput{
		ID:       id,
		Hostname: payload.Hostname,
		IP:       payload.IP,
		OS:       payload.OS,
		Kernel:   payload.Kernel,
		Model:    payload.Model,
		RAM:      payload.RAM,
		CPU:      payload.CPU,
		Storage:  payload.Storage,
	}
	if err := ValidateUpdateInput(input); err != nil {
		return h.errorResponse(ctx, err, op)
	}

	updated, err := h.repo.Update(ctx, Computer{
		ID:       id,
		Hostname: payload.Hostname,
		IP:       payload.IP,
		OS:       payload.OS,
		Kernel:   payload.Kernel,
		Model:    payload.Model,
		RAM:      payload.RAM,
		CPU:      payload.CPU,
		Storage:  payload.Storage,
	})
	if err != nil {
		return h.errorResponse(ctx, err, op)
	}

	return platform.SuccessResponse(http.StatusOK, updated)
}

func (h *Handler) delete(ctx context.Context, id string) (events.APIGatewayProxyResponse, error) {
	const op = "delete computer"

	if err := domain.ValidateID(id); err != nil {
		return h.errorResponse(ctx, err, op)
	}

	deleted, err := h.repo.Delete(ctx, id)
	if err != nil {
		return h.errorResponse(ctx, err, op)
	}

	return platform.SuccessResponse(http.StatusOK, deleted)
}

func (h *Handler) errorResponse(ctx context.Context, err error, operation string) (events.APIGatewayProxyResponse, error) {
	if platform.IsClientError(err) {
		h.logger.InfoContext(ctx, operation+" client error", "error", err.Error())
	} else {
		h.logger.LogError(ctx, operation+" failed", err)
	}

	return platform.ClientErrorResponse(err)
}

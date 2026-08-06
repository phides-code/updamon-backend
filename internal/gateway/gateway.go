// API gateway: auth gate and first-path-segment routing to registered resource handlers.
package gateway

import (
	"context"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/phides-code/updamon-backend/internal/domain"
	"github.com/phides-code/updamon-backend/internal/platform"
)

type ResourceHandler interface {
	Handle(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
}

type Gateway struct {
	logger   *platform.Logger
	cfToken  string
	adminKey string
	handlers map[string]ResourceHandler
}

func NewGateway(logger *platform.Logger) *Gateway {
	return NewGatewayWithAuth(logger, platform.ExpectedCFTToken(), platform.ExpectedAdminKey())
}

func NewGatewayWithAuth(logger *platform.Logger, cfToken string, adminKey string) *Gateway {
	return &Gateway{
		logger:   logger,
		cfToken:  cfToken,
		adminKey: adminKey,
		handlers: make(map[string]ResourceHandler),
	}
}

func (g *Gateway) Register(prefix string, handler ResourceHandler) {
	g.handlers[prefix] = handler
}

func (g *Gateway) Handle(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if req.HTTPMethod != http.MethodOptions {
		if !platform.LocalMode() && !platform.ValidCFTToken(g.cfToken, req.Headers) {
			return platform.ClientErrorResponse(domain.ErrUnauthorized)
		}
		if !platform.ValidAdminKey(g.adminKey, req.Headers) {
			return platform.ClientErrorResponse(domain.ErrUnauthorized)
		}
	}

	logger := g.logger.WithRequestID(req.RequestContext.RequestID)
	logger.InfoContext(ctx, "incoming request",
		"method", req.HTTPMethod,
		"path", req.Path,
	)

	segment, ok := firstPathSegment(req.Path)

	if !ok {
		return platform.ClientErrorResponse(domain.ErrNotFound)
	}

	handler, ok := g.handlers[segment]
	if !ok {
		return platform.ClientErrorResponse(domain.ErrNotFound)
	}

	return handler.Handle(ctx, req)
}

func firstPathSegment(path string) (string, bool) {
	trimmed := strings.Trim(path, "/")
	if trimmed == "" {
		return "", false
	}
	return strings.Split(trimmed, "/")[0], true
}

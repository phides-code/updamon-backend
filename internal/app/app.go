// Composition root: loads AWS config, constructs repositories, and registers resource handlers on the gateway.
package app

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/phides-code/updamon-backend/internal/computer"
	"github.com/phides-code/updamon-backend/internal/gateway"
	"github.com/phides-code/updamon-backend/internal/platform"
	"github.com/phides-code/updamon-backend/internal/sitrep"
)

func Build(ctx context.Context, logger *platform.Logger) (*gateway.Gateway, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("load aws config: %w", err)
	}

	client := dynamodb.NewFromConfig(cfg)
	computerRepo := computer.NewRepository(client)
	sitrepRepo := sitrep.NewRepository(client)
	return buildGateway(logger, computerRepo, sitrepRepo), nil
}

func buildGateway(logger *platform.Logger, computerRepo computer.Repository, sitrepRepo sitrep.Repository) *gateway.Gateway {
	g := gateway.NewGateway(logger)
	g.Register(computer.PathPrefix, computer.NewHandler(computerRepo, logger))
	g.Register(sitrep.PathPrefix, sitrep.NewHandler(sitrepRepo, logger))
	return g
}

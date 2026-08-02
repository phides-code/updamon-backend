// No-op sitrep.Repository for composition smoke tests.
package app

import (
	"context"

	"github.com/phides-code/updamon-backend/internal/sitrep"
)

type stubSitrepRepo struct{}

func (stubSitrepRepo) Create(_ context.Context, _ sitrep.Sitrep) (sitrep.Sitrep, error) {
	return sitrep.Sitrep{}, nil
}
func (stubSitrepRepo) GetByID(_ context.Context, _ string) (sitrep.Sitrep, error) {
	return sitrep.Sitrep{}, nil
}
func (stubSitrepRepo) List(_ context.Context) ([]sitrep.Sitrep, error) {
	return nil, nil
}

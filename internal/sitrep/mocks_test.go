// Mock Repository helpers for sitrep handler and router tests.
package sitrep_test

import (
	"context"

	"github.com/phides-code/updamon-backend/internal/sitrep"
	"github.com/phides-code/updamon-backend/internal/testutil"
)

type mockSitrepRepository struct {
	createFn func(ctx context.Context, b sitrep.Sitrep) (sitrep.Sitrep, error)
	getFn    func(ctx context.Context, id string) (sitrep.Sitrep, error)
	listFn   func(ctx context.Context) ([]sitrep.Sitrep, error)
}

func (m *mockSitrepRepository) Create(ctx context.Context, b sitrep.Sitrep) (sitrep.Sitrep, error) {
	return m.createFn(ctx, b)
}

func (m *mockSitrepRepository) GetByID(ctx context.Context, id string) (sitrep.Sitrep, error) {
	return m.getFn(ctx, id)
}

func (m *mockSitrepRepository) List(ctx context.Context) ([]sitrep.Sitrep, error) {
	return m.listFn(ctx)
}

func emptySitrepRepo() *mockSitrepRepository {
	return &mockSitrepRepository{
		createFn: func(_ context.Context, _ sitrep.Sitrep) (sitrep.Sitrep, error) {
			return sitrep.Sitrep{}, nil
		},
		getFn: func(_ context.Context, _ string) (sitrep.Sitrep, error) {
			return sitrep.Sitrep{}, nil
		},
		listFn: func(_ context.Context) ([]sitrep.Sitrep, error) {
			return nil, nil
		},
	}
}

// dispatchSitrepRepo is a permissive mock for router tests (GET by id succeeds).
func dispatchSitrepRepo() *mockSitrepRepository {
	return &mockSitrepRepository{
		getFn: func(_ context.Context, gotID string) (sitrep.Sitrep, error) {
			return sitrep.Sitrep{
				ID:       gotID,
				Hostname: testutil.TestSitrepHostname,
				AptLog:   testutil.TestSitrepAptLog,
				Last:     testutil.TestSitrepLast,
			}, nil
		},
		listFn: func(_ context.Context) ([]sitrep.Sitrep, error) {
			return nil, nil
		},
		createFn: func(_ context.Context, b sitrep.Sitrep) (sitrep.Sitrep, error) {
			return b, nil
		},
	}
}

func listSitrepRepo(items []sitrep.Sitrep) *mockSitrepRepository {
	return &mockSitrepRepository{
		listFn: func(_ context.Context) ([]sitrep.Sitrep, error) {
			return items, nil
		},
	}
}

func panicSitrepRepo() *mockSitrepRepository {
	panicFn := func() {
		panic("repository must not be called")
	}
	return &mockSitrepRepository{
		createFn: func(context.Context, sitrep.Sitrep) (sitrep.Sitrep, error) {
			panicFn()
			return sitrep.Sitrep{}, nil
		},
		getFn: func(context.Context, string) (sitrep.Sitrep, error) {
			panicFn()
			return sitrep.Sitrep{}, nil
		},
		listFn: func(context.Context) ([]sitrep.Sitrep, error) {
			panicFn()
			return nil, nil
		},
	}
}

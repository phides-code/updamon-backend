// Computer mock repository helpers for handler and router tests.
package computer_test

import (
	"context"

	"github.com/phides-code/go-multi-api/internal/computer"
	"github.com/phides-code/go-multi-api/internal/domain"
)

type mockComputerRepository struct {
	createFn func(ctx context.Context, b computer.Computer) (computer.Computer, error)
	getFn    func(ctx context.Context, id string) (computer.Computer, error)
	listFn   func(ctx context.Context) ([]computer.Computer, error)
	updateFn func(ctx context.Context, b computer.Computer) (computer.Computer, error)
	deleteFn func(ctx context.Context, id string) (computer.Computer, error)
}

func (m *mockComputerRepository) Create(ctx context.Context, b computer.Computer) (computer.Computer, error) {
	return m.createFn(ctx, b)
}

func (m *mockComputerRepository) GetByID(ctx context.Context, id string) (computer.Computer, error) {
	return m.getFn(ctx, id)
}

func (m *mockComputerRepository) List(ctx context.Context) ([]computer.Computer, error) {
	return m.listFn(ctx)
}

func (m *mockComputerRepository) Update(ctx context.Context, b computer.Computer) (computer.Computer, error) {
	return m.updateFn(ctx, b)
}

func (m *mockComputerRepository) Delete(ctx context.Context, id string) (computer.Computer, error) {
	return m.deleteFn(ctx, id)
}

func emptyComputerRepo() *mockComputerRepository {
	return &mockComputerRepository{
		createFn: func(_ context.Context, _ computer.Computer) (computer.Computer, error) {
			return computer.Computer{}, nil
		},
		getFn: func(_ context.Context, _ string) (computer.Computer, error) {
			return computer.Computer{}, nil
		},
		listFn: func(_ context.Context) ([]computer.Computer, error) {
			return nil, nil
		},
		updateFn: func(_ context.Context, _ computer.Computer) (computer.Computer, error) {
			return computer.Computer{}, nil
		},
		deleteFn: func(_ context.Context, _ string) (computer.Computer, error) {
			return computer.Computer{}, nil
		},
	}
}

// dispatchComputerRepo returns a permissive mock for router dispatch tests (GET by id succeeds).
func dispatchComputerRepo() *mockComputerRepository {
	return &mockComputerRepository{
		getFn: func(_ context.Context, gotID string) (computer.Computer, error) {
			return computer.Computer{ID: gotID, Content: "found"}, nil
		},
		listFn: func(_ context.Context) ([]computer.Computer, error) {
			return nil, nil
		},
		createFn: func(_ context.Context, b computer.Computer) (computer.Computer, error) {
			return b, nil
		},
		updateFn: func(_ context.Context, b computer.Computer) (computer.Computer, error) {
			return b, nil
		},
		deleteFn: func(_ context.Context, _ string) (computer.Computer, error) {
			return computer.Computer{}, nil
		},
	}
}

func listComputerRepo(items []computer.Computer) *mockComputerRepository {
	return &mockComputerRepository{
		listFn: func(_ context.Context) ([]computer.Computer, error) {
			return items, nil
		},
	}
}

func updateComputerRepo(wantID string, updated computer.Computer) *mockComputerRepository {
	return &mockComputerRepository{
		updateFn: func(_ context.Context, b computer.Computer) (computer.Computer, error) {
			if b.ID != wantID {
				return computer.Computer{}, domain.ErrNotFound
			}
			return updated, nil
		},
	}
}

func panicComputerRepo() *mockComputerRepository {
	panicFn := func() {
		panic("repository must not be called")
	}
	return &mockComputerRepository{
		createFn: func(context.Context, computer.Computer) (computer.Computer, error) {
			panicFn()
			return computer.Computer{}, nil
		},
		getFn: func(context.Context, string) (computer.Computer, error) {
			panicFn()
			return computer.Computer{}, nil
		},
		listFn: func(context.Context) ([]computer.Computer, error) {
			panicFn()
			return nil, nil
		},
		updateFn: func(context.Context, computer.Computer) (computer.Computer, error) {
			panicFn()
			return computer.Computer{}, nil
		},
		deleteFn: func(context.Context, string) (computer.Computer, error) {
			panicFn()
			return computer.Computer{}, nil
		},
	}
}

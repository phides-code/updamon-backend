// No-op computer.Repository for composition smoke tests.
package app

import (
	"context"

	"github.com/phides-code/updamon-backend/internal/computer"
)

type stubComputerRepo struct{}

func (stubComputerRepo) Create(_ context.Context, _ computer.Computer) (computer.Computer, error) {
	return computer.Computer{}, nil
}
func (stubComputerRepo) GetByID(_ context.Context, _ string) (computer.Computer, error) {
	return computer.Computer{}, nil
}
func (stubComputerRepo) List(_ context.Context) ([]computer.Computer, error) {
	return nil, nil
}
func (stubComputerRepo) Update(_ context.Context, _ computer.Computer) (computer.Computer, error) {
	return computer.Computer{}, nil
}
func (stubComputerRepo) Delete(_ context.Context, _ string) (computer.Computer, error) {
	return computer.Computer{}, nil
}

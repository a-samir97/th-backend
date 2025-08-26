package repository

import (
	"context"
	"testing"

	"thamaniyah/internal/domain"
)

// TestMediaRepository is an interface test to ensure all implementations
// satisfy the MediaRepository interface
func TestMediaRepositoryInterface(t *testing.T) {
	// This is a compile-time check to ensure our interface is properly defined
	var _ MediaRepository = (*MockMediaRepository)(nil)
}

// MockMediaRepository can be used in tests
type MockMediaRepository struct{}

func (m *MockMediaRepository) Create(ctx context.Context, media *domain.Media) error {
	return nil
}

func (m *MockMediaRepository) GetByID(ctx context.Context, id string) (*domain.Media, error) {
	return nil, nil
}

func (m *MockMediaRepository) GetAll(ctx context.Context, limit, offset int) ([]*domain.Media, error) {
	return nil, nil
}

func (m *MockMediaRepository) Update(ctx context.Context, media *domain.Media) error {
	return nil
}

func (m *MockMediaRepository) Delete(ctx context.Context, id string) error {
	return nil
}

func (m *MockMediaRepository) GetByStatus(ctx context.Context, status domain.MediaStatus, limit, offset int) ([]*domain.Media, error) {
	return nil, nil
}

func (m *MockMediaRepository) UpdateStatus(ctx context.Context, id string, status domain.MediaStatus) error {
	return nil
}

func (m *MockMediaRepository) GetTotal(ctx context.Context) (int64, error) {
	return 0, nil
}

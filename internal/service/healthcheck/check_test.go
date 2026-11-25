package healthcheck

import (
	repo_mocks "bookmark-management/internal/repository/mocks"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

var errRepo = errors.New("repo error")

func TestService_Check(t *testing.T) {
	t.Parallel()

	const (
		serviceName = "bookmark_service"
		instanceId  = "instance_1"
	)

	testcases := []struct {
		name         string
		setupMock    func(t *testing.T, ctx context.Context) *repo_mocks.HealthcheckRepository
		expectedName string
		expectedId   string
		expectedErr  error
	}{
		{
			name: "success",
			setupMock: func(t *testing.T, ctx context.Context) *repo_mocks.HealthcheckRepository {
				mockRepo := repo_mocks.NewHealthcheckRepository(t)
				mockRepo.On("Check", ctx).Return(nil)
				return mockRepo
			},
			expectedName: serviceName,
			expectedId:   instanceId,
		},
		{
			name: "repository error",
			setupMock: func(t *testing.T, ctx context.Context) *repo_mocks.HealthcheckRepository {
				mockRepo := repo_mocks.NewHealthcheckRepository(t)
				mockRepo.On("Check", ctx).Return(errRepo)
				return mockRepo
			},
			expectedErr: errRepo,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			svc := NewService(serviceName, instanceId, tc.setupMock(t, ctx))

			name, id, err := svc.Check(ctx)

			assert.Equal(t, tc.expectedName, name)
			assert.Equal(t, tc.expectedId, id)
			if tc.expectedErr != nil {
				assert.ErrorIs(t, err, tc.expectedErr)
				assert.ErrorContains(t, err, "healthcheck")
				return
			}
			assert.NoError(t, err)
		})
	}
}

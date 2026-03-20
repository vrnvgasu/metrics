package metric

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	models "github.com/vrnvgasu/metrics/internal/model"
	"github.com/vrnvgasu/metrics/internal/repository"
	mockrepository "github.com/vrnvgasu/metrics/internal/repository/mocks"
)

func Test_createOrUpdate(t *testing.T) {
	t.Parallel()

	controller := gomock.NewController(t)
	defer controller.Finish()

	tests := []struct {
		name        string
		m           models.Metrics
		storage     func() repository.Storage
		expectedErr require.ErrorAssertionFunc
	}{
		{
			name: "no metrics; add gauge metric",
			m:    models.Metrics{ID: "1", MType: models.Gauge},
			storage: func() repository.Storage {
				storageMock := mockrepository.NewMockStorage(controller)
				storageMock.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil)

				return storageMock
			},
			expectedErr: require.NoError,
		},
		{
			name: "no metrics; add counter metric",
			m:    models.Metrics{ID: "1", MType: models.Counter},
			storage: func() repository.Storage {
				storageMock := mockrepository.NewMockStorage(controller)
				storageMock.EXPECT().GetByTypeAndID(gomock.Any(), models.Counter, "1").Return(nil, repository.ErrNotFound)
				storageMock.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil)

				return storageMock
			},
			expectedErr: require.NoError,
		},
		{
			name: "no metrics; add wrong metric",
			m:    models.Metrics{ID: "1", MType: "dummy"},
			storage: func() repository.Storage {
				storageMock := mockrepository.NewMockStorage(controller)

				return storageMock
			},
			expectedErr: require.Error,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := &Service{
				storage: tt.storage(),
			}

			err := s.createOrUpdate(t.Context(), []*models.Metrics{&tt.m})
			tt.expectedErr(t, err)
		})
	}
}

package cresource_usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Uikola/distributedCalculator2/orchestrator/internal/entity"
	"github.com/Uikola/distributedCalculator2/orchestrator/internal/usecase/cresource_usecase"
	"github.com/Uikola/distributedCalculator2/orchestrator/internal/usecase/cresource_usecase/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestListCResources(t *testing.T) {
	cases := []struct {
		name string

		mockErr  error
		mockResp []entity.CResource

		want    map[string]string
		wantErr error
	}{
		{
			name: "success",

			mockResp: []entity.CResource{
				{ID: 1, Name: "test1", Address: "localhost:30451", Expression: "", Occupied: false, OrchestratorAlive: true},
				{ID: 2, Name: "test2", Address: "localhost:30452", Expression: "2 + 2", Occupied: true, OrchestratorAlive: true},
			},

			want: map[string]string{"localhost:30451": "", "localhost:30452": "2 + 2"},
		},
		{
			name: "repository error",

			mockErr: errors.New("mock err"),

			wantErr: errors.New("mock err"),
		},
	}

	for _, tCase := range cases {
		t.Run(tCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepository := mocks.NewMockcResourceRepository(ctrl)
			mockRepository.EXPECT().ListCResources(context.Background()).Return(tCase.mockResp, tCase.mockErr)

			useCase := cresource_usecase.NewUseCaseImpl(mockRepository)

			got, gotErr := useCase.ListCResources(context.Background())

			require.Equal(t, tCase.want, got)
			require.Equal(t, tCase.wantErr, gotErr)
		})
	}
}

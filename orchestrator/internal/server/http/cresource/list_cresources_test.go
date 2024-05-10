package cresource_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Uikola/distributedCalculator2/orchestrator/internal/server/http/cresource"
	"github.com/Uikola/distributedCalculator2/orchestrator/internal/server/http/cresource/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestListCResources(t *testing.T) {
	cases := []struct {
		name    string
		expCode int

		mockErr  error
		mockResp map[string]string

		want    map[string]string
		respErr bool
	}{
		{
			name:    "success",
			expCode: http.StatusOK,

			mockResp: map[string]string{"test": "localhost:30005", "test2": "localhost:30060"},

			want: map[string]string{"test": "localhost:30005", "test2": "localhost:30060"},
		},
		{
			name:    "use case error",
			expCode: http.StatusInternalServerError,

			mockErr: errors.New("mock err"),

			want:    map[string]string{"reason": "internal error"},
			respErr: true,
		},
	}

	for _, tCase := range cases {
		t.Run(tCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUseCase := mocks.NewMockcResourceUseCase(ctrl)

			if !tCase.respErr || tCase.mockErr != nil {
				mockUseCase.EXPECT().ListCResources(context.Background()).Return(tCase.mockResp, tCase.mockErr)
			}

			handler := cresource.NewHandler(mockUseCase)

			req, err := http.NewRequest(http.MethodPost, "/cresources", nil)
			require.NoError(t, err)
			rec := httptest.NewRecorder()

			handler.ListCResources(rec, req)
			defer rec.Result().Body.Close()

			var got map[string]string
			err = json.NewDecoder(rec.Result().Body).Decode(&got)
			require.NoError(t, err)

			require.Equal(t, tCase.expCode, rec.Result().StatusCode)
			require.Equal(t, tCase.want, got)
		})
	}
}

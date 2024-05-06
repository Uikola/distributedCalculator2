package user_test

import (
	"github.com/Uikola/distributedCalculator2/orchestrator/internal/entity"
	"github.com/Uikola/distributedCalculator2/orchestrator/internal/errorz"
	"github.com/Uikola/distributedCalculator2/orchestrator/internal/server/http/user"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestValidateUpdateOperation(t *testing.T) {
	cases := []struct {
		name string
		req  entity.UpdateOperationRequest
		err  error
	}{
		{
			name: "valid operation and time",
			req:  entity.UpdateOperationRequest{Operation: "+", Time: 10},
		},
		{
			name: "invalid operation",
			req:  entity.UpdateOperationRequest{Operation: "test", Time: 10},
			err:  errorz.ErrInvalidOperation,
		},
	}

	for _, tCase := range cases {
		tCase := tCase
		t.Run(tCase.name, func(t *testing.T) {

			err := user.ValidateUpdateOperation(tCase.req)
			require.Equal(t, tCase.err, err)
		})
	}
}

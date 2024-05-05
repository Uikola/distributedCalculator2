package expression_test

import (
	"github.com/Uikola/distributedCalculator2/orchestrator/internal/errorz"
	"github.com/Uikola/distributedCalculator2/orchestrator/internal/server/http/expression"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestValidateExpression(t *testing.T) {
	cases := []struct {
		name string
		expr string
		err  error
	}{
		{
			name: "valid expression",
			expr: "2 + 2",
		},
		{
			name: "invalid expression syntax",
			expr: "1 + 2 + test",
			err:  errorz.ErrInvalidExpression,
		},
		{
			name: "invalid expression",
			expr: "1 + 2()3(",
			err:  errorz.ErrInvalidExpression,
		},
	}

	for _, tCase := range cases {
		tCase := tCase
		t.Run(tCase.name, func(t *testing.T) {

			err := expression.ValidateExpression(tCase.expr)
			require.Equal(t, tCase.err, err)
		})
	}
}

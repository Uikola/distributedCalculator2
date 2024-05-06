package user_usecase_test

import (
	"github.com/Uikola/distributedCalculator2/orchestrator/internal/usecase/user_usecase"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGeneratePasswordHash(t *testing.T) {
	want := "82f8809f42d911d1bd5199021d69d15ea91d1fad"

	password := "testPassword"
	got := user_usecase.GeneratePasswordHash(password)

	require.Equal(t, want, got)
}

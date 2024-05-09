package user_usecase_test

import (
	"context"
	"errors"
	"github.com/Uikola/distributedCalculator2/orchestrator/internal/entity"
	"github.com/Uikola/distributedCalculator2/orchestrator/internal/errorz"
	"github.com/Uikola/distributedCalculator2/orchestrator/internal/usecase/user_usecase"
	"github.com/Uikola/distributedCalculator2/orchestrator/internal/usecase/user_usecase/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestGeneratePasswordHash(t *testing.T) {
	want := "82f8809f42d911d1bd5199021d69d15ea91d1fad"

	password := "testPassword"
	got := user_usecase.GeneratePasswordHash(password)

	require.Equal(t, want, got)
}

func TestCreate(t *testing.T) {
	cases := []struct {
		name  string
		input entity.RegisterRequest

		existsMockInput string
		createMockInput entity.User
		existsMockResp  bool
		existsMockErr   error
		createMockErr   error

		wantErr error
	}{
		{
			name:  "success",
			input: entity.RegisterRequest{Login: "TestUser", Password: "TestUserPass123"},

			existsMockInput: "TestUser",
			createMockInput: entity.User{Login: "TestUser", Password: user_usecase.GeneratePasswordHash("TestUserPass123")},
			existsMockResp:  false,
		},
		{
			name:  "exists mock error",
			input: entity.RegisterRequest{Login: "TestUser", Password: "TestUserPass123"},

			existsMockInput: "TestUser",
			existsMockErr:   errors.New("exists mock err"),

			wantErr: errors.New("exists mock err"),
		},
		{
			name:  "user already exists",
			input: entity.RegisterRequest{Login: "TestUser", Password: "TestUserPass123"},

			existsMockInput: "TestUser",
			existsMockResp:  true,

			wantErr: errorz.ErrUserAlreadyExists,
		},
		{
			name:  "create mock error",
			input: entity.RegisterRequest{Login: "TestUser", Password: "TestUserPass123"},

			existsMockInput: "TestUser",
			createMockInput: entity.User{Login: "TestUser", Password: user_usecase.GeneratePasswordHash("TestUserPass123")},
			existsMockResp:  false,
			createMockErr:   errors.New("create mock error"),

			wantErr: errors.New("create mock error"),
		},
	}

	for _, tCase := range cases {
		t.Run(tCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepository := mocks.NewMockuserRepository(ctrl)

			mockRepository.EXPECT().Exists(context.Background(), tCase.existsMockInput).Return(tCase.existsMockResp, tCase.existsMockErr)
			if tCase.existsMockErr == nil && !tCase.existsMockResp {
				mockRepository.EXPECT().Create(context.Background(), tCase.createMockInput).Return(tCase.createMockErr)
			}

			usecase := user_usecase.NewUseCaseImpl(mockRepository)

			gotErr := usecase.Create(context.Background(), tCase.input)

			require.Equal(t, tCase.wantErr, gotErr)
		})
	}
}

func TestLogin(t *testing.T) {
	cases := []struct {
		name  string
		input entity.LoginRequest

		mockInput entity.LoginRequest
		mockResp  entity.User
		mockErr   error

		wantErr error
	}{
		{
			name:  "success",
			input: entity.LoginRequest{Login: "TestUser", Password: "TestPass123"},

			mockInput: entity.LoginRequest{Login: "TestUser", Password: user_usecase.GeneratePasswordHash("TestPass123")},
			mockResp:  entity.User{ID: 1, Login: "TestUser", Password: user_usecase.GeneratePasswordHash("TestPass123")},
		},
		{
			name:  "mock error",
			input: entity.LoginRequest{Login: "TestUser", Password: "TestPass123"},

			mockInput: entity.LoginRequest{Login: "TestUser", Password: user_usecase.GeneratePasswordHash("TestPass123")},
			mockErr:   errors.New("mock err"),

			wantErr: errors.New("mock err"),
		},
	}

	for _, tCase := range cases {
		t.Run(tCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepository := mocks.NewMockuserRepository(ctrl)

			mockRepository.EXPECT().GetUser(context.Background(), tCase.mockInput.Login, tCase.mockInput.Password).Return(tCase.mockResp, tCase.mockErr)

			usecase := user_usecase.NewUseCaseImpl(mockRepository)

			_, gotErr := usecase.Login(context.Background(), tCase.input)

			require.Equal(t, tCase.wantErr, gotErr)
		})
	}
}

func TestListOperations(t *testing.T) {
	cases := []struct {
		name string

		userID   uint
		mockResp entity.User
		mockErr  error

		want    map[string]uint
		wantErr error
	}{
		{
			name: "success",

			userID:   1,
			mockResp: entity.User{ID: 1, Login: "TestUser", Password: "TestPass123", Addition: 10, Subtraction: 10, Multiplication: 10, Division: 10},

			want: map[string]uint{"+": 10, "-": 10, "*": 10, "/": 10},
		},
		{
			name: "mock err",

			userID:  1,
			mockErr: errors.New("mock err"),

			wantErr: errors.New("mock err"),
		},
	}

	for _, tCase := range cases {
		t.Run(tCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepository := mocks.NewMockuserRepository(ctrl)

			ctx := context.WithValue(context.Background(), "userID", tCase.userID)

			mockRepository.EXPECT().GetByID(ctx, tCase.userID).Return(tCase.mockResp, tCase.mockErr)

			usecase := user_usecase.NewUseCaseImpl(mockRepository)

			got, gotErr := usecase.ListOperations(ctx)

			require.Equal(t, tCase.wantErr, gotErr)
			require.Equal(t, tCase.want, got)
		})
	}
}

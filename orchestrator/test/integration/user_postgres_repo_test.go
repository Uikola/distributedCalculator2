package integration

import (
	"context"
	"github.com/Uikola/distributedCalculator2/orchestrator/internal/db"
	"github.com/Uikola/distributedCalculator2/orchestrator/internal/db/repository/postgres"
	"github.com/Uikola/distributedCalculator2/orchestrator/internal/entity"
	"github.com/Uikola/distributedCalculator2/orchestrator/internal/errorz"
	"github.com/Uikola/distributedCalculator2/orchestrator/pkg/zlog"
	"github.com/Uikola/distributedCalculator2/orchestrator/test/testhelper"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type UserRepoTestSuite struct {
	suite.Suite
	pgContainer *testhelper.PostgresContainer
	repository  *postgres.UserRepository
	ctx         context.Context
	userID      uint
}

func (suite *UserRepoTestSuite) SetupSuite() {
	log.Logger = zlog.Default(true, "dev", zerolog.InfoLevel)
	suite.userID = 1
	suite.ctx = context.Background()
	pgContainer, err := testhelper.CreatePostgresContainer(suite.ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create test postgres container")
	}
	suite.pgContainer = pgContainer

	database := db.InitDB(suite.pgContainer.ConnStr)

	repository := postgres.NewUserRepository(database)
	suite.repository = repository
}

func (suite *UserRepoTestSuite) TearDownSuite() {
	if err := suite.pgContainer.Terminate(suite.ctx); err != nil {
		log.Fatal().Err(err).Msg("failed to terminate test postgres container")
	}
}

func (suite *UserRepoTestSuite) SetupTest() {
	if err := suite.repository.Create(suite.ctx, entity.User{Login: "TestUser", Password: "HashedTestPass"}); err != nil {
		log.Fatal().Err(err).Msg("failed to create test user data")
	}
	suite.userID++
}

func (suite *UserRepoTestSuite) TearDownTest() {
	if err := suite.repository.CleanUp(suite.ctx); err != nil {
		log.Fatal().Err(err).Msg("failed to clean up users table")
	}
}

func (suite *UserRepoTestSuite) TestCreate() {
	t := suite.T()

	err := suite.repository.Create(suite.ctx, entity.User{
		Login: "TestTest", Password: "HashedTestPass",
	})
	require.NoError(t, err)
}

func (suite *UserRepoTestSuite) TestExists() {
	t := suite.T()

	cases := []struct {
		name string

		login string

		want bool
	}{
		{
			name: "exists",

			login: "TestUser",

			want: true,
		},
		{
			name: "no exists",

			login: "Test",

			want: false,
		},
	}

	for _, tCase := range cases {
		t.Run(tCase.name, func(t *testing.T) {
			got, err := suite.repository.Exists(suite.ctx, tCase.login)
			require.NoError(t, err)
			require.Equal(t, tCase.want, got)
		})
	}
}

func (suite *UserRepoTestSuite) TestGetUser() {
	t := suite.T()

	cases := []struct {
		name string

		login    string
		password string

		want    entity.User
		wantErr error
	}{
		{
			name: "success",

			login:    "TestUser",
			password: "HashedTestPass",

			want: entity.User{ID: suite.userID, Login: "TestUser", Password: "HashedTestPass"},
		},
		{
			name: "user not found",

			login:    "TestUnknownUser",
			password: "TestPass123",

			wantErr: errorz.ErrUserNotFound,
		},
	}

	for _, tCase := range cases {
		t.Run(tCase.name, func(t *testing.T) {
			got, gotErr := suite.repository.GetUser(suite.ctx, tCase.login, tCase.password)
			require.Equal(t, tCase.want, got)
			require.Equal(t, tCase.wantErr, gotErr)
		})
	}
}

func (suite *UserRepoTestSuite) TestGetByID() {
	t := suite.T()

	want := entity.User{ID: suite.userID, Login: "TestUser", Password: "HashedTestPass", Addition: 5, Subtraction: 5, Multiplication: 5, Division: 5}
	got, err := suite.repository.GetByID(suite.ctx, suite.userID)
	require.NoError(t, err)
	require.Equal(t, want, got)
}

func (suite *UserRepoTestSuite) TestUpdateOperation() {
	t := suite.T()

	err := suite.repository.UpdateOperation(suite.ctx, suite.userID, "addition", 10)
	require.NoError(t, err)
}

func TestUserRepoTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepoTestSuite))
}

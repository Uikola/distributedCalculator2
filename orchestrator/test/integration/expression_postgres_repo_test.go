package integration

import (
	"context"
	"testing"

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
)

type ExpressionRepoTestSuite struct {
	suite.Suite
	pgContainer   *testhelper.PostgresContainer
	repository    *postgres.ExpressionRepository
	ctx           context.Context
	expressionID1 uint
	expressionID2 uint
}

func (suite *ExpressionRepoTestSuite) SetupSuite() {
	log.Logger = zlog.Default(true, "dev", zerolog.InfoLevel)
	suite.expressionID1 = 1
	suite.expressionID2 = 2
	suite.ctx = context.Background()
	pgContainer, err := testhelper.CreatePostgresContainer(suite.ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create test postgres container")
	}
	suite.pgContainer = pgContainer

	database := db.InitPostgres(suite.pgContainer.ConnStr)

	repository := postgres.NewExpressionRepository(database)
	suite.repository = repository
}

func (suite *ExpressionRepoTestSuite) TearDownSuite() {
	if err := suite.pgContainer.Terminate(suite.ctx); err != nil {
		log.Fatal().Err(err).Msg("failed to terminate test postgres contaioner")
	}
}

func (suite *ExpressionRepoTestSuite) SetupTest() {
	_, err := suite.repository.AddExpression(suite.ctx, entity.Expression{Expression: "1 + 1", Status: entity.InProgress, OwnerID: 1})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create test computing resource data")
	}
	suite.expressionID1++

	_, err = suite.repository.AddExpression(suite.ctx, entity.Expression{Expression: "2 + 2", Status: entity.InProgress, OwnerID: 1})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create test computing resource data")
	}
	suite.expressionID2++
}

func (suite *ExpressionRepoTestSuite) TearDownTest() {
	if err := suite.repository.CleanUp(suite.ctx); err != nil {
		log.Fatal().Err(err).Msg("failed to clean up cresources table")
	}
}

func (suite *ExpressionRepoTestSuite) TestAddExpression() {
	t := suite.T()

	want := entity.Expression{ID: suite.expressionID2, Expression: "3 + 3", Status: entity.InProgress, OwnerID: 1}

	got, err := suite.repository.AddExpression(suite.ctx, entity.Expression{Expression: "3 + 3", Status: entity.InProgress, OwnerID: 1})
	require.NoError(t, err)
	require.Equal(t, want, got)
}

func (suite *ExpressionRepoTestSuite) TestSetErrorStatus() {
	t := suite.T()

	err := suite.repository.SetErrorStatus(suite.ctx, suite.expressionID1)
	require.NoError(t, err)
}

func (suite *ExpressionRepoTestSuite) TestUpdateCResource() {
	t := suite.T()

	err := suite.repository.UpdateCResource(suite.ctx, suite.expressionID1, 1)
	require.NoError(t, err)
}

func (suite *ExpressionRepoTestSuite) TestGetExpressionByID() {
	t := suite.T()

	cases := []struct {
		name string

		expressionID uint

		want    entity.Expression
		wantErr error
	}{
		{
			name: "success",

			expressionID: suite.expressionID1 + 1,

			want: entity.Expression{ID: suite.expressionID1 + 1, Expression: "1 + 1", Status: entity.InProgress, OwnerID: 1},
		},
		{
			name: "expression not found",

			expressionID: 0,

			wantErr: errorz.ErrExpressionNotFound,
		},
	}

	for _, tCase := range cases {
		t.Run(tCase.name, func(t *testing.T) {
			got, err := suite.repository.GetExpressionByID(suite.ctx, tCase.expressionID)
			require.Equal(t, tCase.wantErr, err)
			require.Equal(t, tCase.want, got)
		})
	}
}

func (suite *ExpressionRepoTestSuite) TestListExpressions() {
	t := suite.T()

	want := []entity.Expression{
		{ID: suite.expressionID1 + 2, Expression: "1 + 1", Status: entity.InProgress, OwnerID: 1},
		{ID: suite.expressionID2 + 2, Expression: "2 + 2", Status: entity.InProgress, OwnerID: 1},
	}

	got, err := suite.repository.ListExpressions(suite.ctx, 1)
	require.NoError(t, err)
	require.Equal(t, want, got)
}

func (suite *ExpressionRepoTestSuite) TestUpdateResult() {
	t := suite.T()

	err := suite.repository.UpdateResult(suite.ctx, 1, "2")
	require.NoError(t, err)
}

func (suite *ExpressionRepoTestSuite) TestSetSuccessStatus() {
	t := suite.T()

	err := suite.repository.SetSuccessStatus(suite.ctx, 1)
	require.NoError(t, err)
}

func TestExpressionRepoTestSuite(t *testing.T) {
	suite.Run(t, new(ExpressionRepoTestSuite))
}

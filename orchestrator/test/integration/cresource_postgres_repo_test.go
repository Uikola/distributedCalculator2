package integration

import (
	"context"
	"testing"

	"github.com/Uikola/distributedCalculator2/orchestrator/internal/db"
	"github.com/Uikola/distributedCalculator2/orchestrator/internal/db/repository/postgres"
	"github.com/Uikola/distributedCalculator2/orchestrator/internal/entity"
	"github.com/Uikola/distributedCalculator2/orchestrator/pkg/zlog"
	"github.com/Uikola/distributedCalculator2/orchestrator/test/testhelper"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type CResourceRepoTestSuite struct {
	suite.Suite
	pgContainer *testhelper.PostgresContainer
	repository  *postgres.CResourceRepository
	ctx         context.Context
	cResourceID uint
}

func (suite *CResourceRepoTestSuite) SetupSuite() {
	log.Logger = zlog.Default(true, "dev", zerolog.InfoLevel)
	suite.cResourceID = 1
	suite.ctx = context.Background()
	pgContainer, err := testhelper.CreatePostgresContainer(suite.ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create test postgres container")
	}
	suite.pgContainer = pgContainer

	database := db.InitPostgres(suite.pgContainer.ConnStr)

	repository := postgres.NewCResourceRepository(database)
	suite.repository = repository
}

func (suite *CResourceRepoTestSuite) TearDownSuite() {
	if err := suite.pgContainer.Terminate(suite.ctx); err != nil {
		log.Fatal().Err(err).Msg("failed to terminate test postgres contaioner")
	}
}

func (suite *CResourceRepoTestSuite) SetupTest() {
	if err := suite.repository.Create(suite.ctx, entity.CResource{Name: "Test", Address: "localhost:35678"}); err != nil {
		log.Fatal().Err(err).Msg("failed to create test computing resource data")
	}
	suite.cResourceID++
}

func (suite *CResourceRepoTestSuite) TearDownTest() {
	if err := suite.repository.CleanUp(suite.ctx); err != nil {
		log.Fatal().Err(err).Msg("failed to clean up cresources table")
	}
}

func (suite *CResourceRepoTestSuite) TestCreate() {
	t := suite.T()

	err := suite.repository.Create(suite.ctx, entity.CResource{
		Name:    "TestCResource",
		Address: "localhost:35679",
	})
	require.NoError(t, err)
}

func (suite *CResourceRepoTestSuite) TestExists() {
	t := suite.T()

	cases := []struct {
		name string

		cResourceName string
		address       string

		want bool
	}{
		{
			name: "exists",

			cResourceName: "Test",
			address:       "localhost:35678",

			want: true,
		},
		{
			name: "no exists",

			cResourceName: "NewTest",
			address:       "localhost:34550",

			want: false,
		},
	}

	for _, tCase := range cases {
		t.Run(tCase.name, func(t *testing.T) {
			got, err := suite.repository.Exists(suite.ctx, tCase.cResourceName, tCase.address)
			require.NoError(t, err)
			require.Equal(t, tCase.want, got)
		})
	}
}

func (suite *CResourceRepoTestSuite) TestSetOrchestatorHealth() {
	t := suite.T()

	err := suite.repository.SetOrchestatorHealth(suite.ctx, "Test", true)
	require.NoError(t, err)
}

func (suite *CResourceRepoTestSuite) TestDelete() {
	t := suite.T()

	err := suite.repository.Delete(suite.ctx, "Test")
	require.NoError(t, err)
}

func (suite *CResourceRepoTestSuite) TestAssignExpressionToCResource() {
	t := suite.T()

	cresource, err := suite.repository.AssignExpressionToCResource(suite.ctx, entity.Expression{ID: 1, Expression: "1 + 1", Status: entity.InProgress, OwnerID: 1})
	require.NoError(t, err)
	require.Equal(t, true, cresource.Occupied)
	require.Equal(t, "1 + 1", cresource.Expression)
}

func (suite *CResourceRepoTestSuite) TestUnlinkExpressionFromCResource() {
	t := suite.T()

	err := suite.repository.UnlinkExpressionFromCResource(suite.ctx, entity.Expression{ID: 1, CalculatedBy: suite.cResourceID, OwnerID: 1})
	require.NoError(t, err)
}

func (suite *CResourceRepoTestSuite) TestListCResources() {
	t := suite.T()

	want := []entity.CResource{
		{ID: suite.cResourceID, Name: "Test", Address: "localhost:35678"},
	}

	got, err := suite.repository.ListCResources(suite.ctx)
	require.NoError(t, err)
	require.Equal(t, want, got)
}

func (suite *CResourceRepoTestSuite) TestGetByName() {
	t := suite.T()

	want := entity.CResource{ID: suite.cResourceID, Name: "Test", Address: "localhost:35678"}

	got, err := suite.repository.GetByName(suite.ctx, "Test")
	require.NoError(t, err)
	require.Equal(t, want, got)
}

func TestCResourceRepoTestSuite(t *testing.T) {
	suite.Run(t, new(CResourceRepoTestSuite))
}

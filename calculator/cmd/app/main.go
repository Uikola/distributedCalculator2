package main

import (
	"net"
	"os"

	"github.com/Uikola/distributedCalculator2/calculator/internal/db"
	"github.com/Uikola/distributedCalculator2/calculator/internal/db/repository/postgres"
	"github.com/Uikola/distributedCalculator2/calculator/internal/server/grpc/expression"
	"github.com/Uikola/distributedCalculator2/calculator/internal/server/grpc/heartbeat"
	"github.com/Uikola/distributedCalculator2/calculator/internal/usecase/expression_usecase"
	heartbeatClient "github.com/Uikola/distributedCalculator2/orchestrator/pkg/grpc/client/heartbeat"
	"github.com/Uikola/distributedCalculator2/orchestrator/pkg/zlog"
	expressionpb "github.com/Uikola/distributedCalculator2/proto/expression"
	heartbeatpb "github.com/Uikola/distributedCalculator2/proto/heartbeat"
	"github.com/forscht/namegen"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	log.Logger = zlog.Default(true, "dev", zerolog.InfoLevel)

	database := db.InitDB(os.Getenv("POSTGRES_CONN"))
	defer database.Close()

	expressionRepository := postgres.NewExpressionRepository(database)
	cResourceRepository := postgres.NewCResourceRepository(database)
	userRepository := postgres.NewUserRepository(database)
	expressionUseCase := expression_usecase.NewUseCaseImpl(expressionRepository, cResourceRepository, userRepository)

	name := namegen.New().WithNumberOfWords(1).WithStyle(namegen.Lowercase).Generate()

	orchConn, err := grpc.NewClient(os.Getenv("ORCH_ADDR"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error().Err(err).Msg("failed to connect to the new client")
		os.Exit(1)
	}
	defer orchConn.Close()

	heartbeatOrchClient := heartbeatClient.NewClient(orchConn, database)

	isSuccess, addr, err := heartbeatOrchClient.Register(name)
	if err != nil || !isSuccess {
		log.Error().Err(err).Msg("failed to register service")
		os.Exit(1)
	}

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Error().Err(err).Msg("error starting tcp listener")
		os.Exit(1)
	}

	grpcServer := grpc.NewServer()
	heartbeatCalcServiceServer := heartbeat.NewServer()
	expressionServiceServer := expression.NewServer(expressionUseCase)

	heartbeatpb.RegisterHeartbeatServiceServer(grpcServer, heartbeatCalcServiceServer)
	expressionpb.RegisterExpressionServiceServer(grpcServer, expressionServiceServer)

	go heartbeatOrchClient.SendHeartbeatToOrch(name)

	if err = grpcServer.Serve(lis); err != nil {
		log.Error().Err(err).Msg("error serving grpc")
		os.Exit(1)
	}
}

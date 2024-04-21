package main

import (
	"context"
	"errors"
	"net/http"
	"sync"

	"github.com/Uikola/distributedCalculator2/orchestrator/internal/db"
	"github.com/Uikola/distributedCalculator2/orchestrator/internal/db/repository/postgres"
	"github.com/Uikola/distributedCalculator2/orchestrator/internal/server/grpc/server/heartbeat"
	server "github.com/Uikola/distributedCalculator2/orchestrator/internal/server/http"
	"github.com/Uikola/distributedCalculator2/orchestrator/internal/server/http/cresource"
	"github.com/Uikola/distributedCalculator2/orchestrator/internal/server/http/expression"
	"github.com/Uikola/distributedCalculator2/orchestrator/internal/server/http/user"
	"github.com/Uikola/distributedCalculator2/orchestrator/internal/usecase/cresource_usecase"
	"github.com/Uikola/distributedCalculator2/orchestrator/internal/usecase/expression_usecasse"
	"github.com/Uikola/distributedCalculator2/orchestrator/internal/usecase/user_usecase"
	"github.com/Uikola/distributedCalculator2/orchestrator/pkg/zlog"
	pb "github.com/Uikola/distributedCalculator2/proto/heartbeat"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"

	"net"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	log.Logger = zlog.Default(true, "dev", zerolog.InfoLevel)

	database := db.InitDB(os.Getenv("POSTGRES_CONN"))

	m := CreateMigrate()
	if err = m.Up(); err != nil {
		log.Error().Err(err).Msg("failed to create users table")
	}

	cResourceRepo := postgres.NewCResourceRepository(database)
	expressionRepo := postgres.NewExpressionRepository(database)
	userRepo := postgres.NewUserRepository(database)

	cResourceUseCase := cresource_usecase.NewUseCaseImpl(cResourceRepo)
	expressionUseCase := expression_usecasse.NewUseCaseImpl(expressionRepo, cResourceRepo)
	userUseCase := user_usecase.NewUseCaseImpl(userRepo)

	cResourceHandler := cresource.NewHandler(cResourceUseCase)
	userHandler := user.NewHandler(userUseCase)
	expressionHandler := expression.NewHandler(expressionUseCase)

	srv := server.NewServer(userHandler, expressionHandler, cResourceHandler)

	httpServer := &http.Server{
		Addr:    os.Getenv("HTTP_PORT"),
		Handler: srv,
	}

	lis, err := net.Listen("tcp", os.Getenv("ORCH_ADDR"))

	if err != nil {
		log.Error().Msg(err.Error())
		os.Exit(1)
	}

	grpcServer := grpc.NewServer()

	heartbeatServiceServer := heartbeat.NewServer(cResourceUseCase)

	pb.RegisterHeartbeatServiceServer(grpcServer, heartbeatServiceServer)

	go func() {
		log.Info().Msg("serving grpc")
		if err = grpcServer.Serve(lis); err != nil {
			log.Error().Msg(err.Error())
			os.Exit(1)
		}
	}()

	go func() {
		log.Info().Msg("listening http")
		if err = httpServer.ListenAndServe(); err != nil && !errors.Is(http.ErrServerClosed, err) {
			log.Error().Msg(err.Error())
			os.Exit(1)
		}
	}()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		if err = httpServer.Shutdown(ctx); err != nil {
			log.Error().Msg(err.Error())
			os.Exit(1)
		}
	}()
	wg.Wait()
}

func CreateMigrate() *migrate.Migrate {
	m, err := migrate.New(
		"file://migrations/",
		os.Getenv("PGX_URL"),
	)

	if err != nil {
		log.Error().Msg(err.Error())
	}

	return m
}

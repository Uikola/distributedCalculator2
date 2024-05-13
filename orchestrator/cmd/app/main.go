package main

import (
	"context"
	"errors"
	"github.com/Uikola/distributedCalculator2/orchestrator/pkg/trace"
	"github.com/redis/go-redis/v9"
	"net/http"
	"os/signal"
	"sync"

	"github.com/Uikola/distributedCalculator2/orchestrator/internal/db"
	"github.com/Uikola/distributedCalculator2/orchestrator/internal/db/repository/postgres"
	r "github.com/Uikola/distributedCalculator2/orchestrator/internal/db/repository/redis"
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
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	if err := godotenv.Load(); err != nil {
		log.Error().Msg(err.Error())
		os.Exit(1)
	}

	log.Logger = zlog.Default(true, "dev", zerolog.InfoLevel)
	if err := run(ctx); err != nil {
		log.Error().Msg(err.Error())
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	otelShutdown, err := trace.SetupOTelSDK(ctx, "http://localhost:14268/api/traces", "Orchestrator Service")
	if err != nil {
		return err
	}

	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
	}()

	database, err := db.OtelInitPostgres(os.Getenv("POSTGRES_CONN"))
	if err != nil {
		return err
	}

	m := createMigrate()
	if err = m.Up(); err != nil {
		log.Error().Err(err).Msg("failed to create users table")
	}

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	cache := r.NewCache(client)

	cResourceRepo := postgres.NewCResourceRepository(database)
	expressionRepo := postgres.NewExpressionRepository(database)
	userRepo := postgres.NewUserRepository(database)

	cResourceUseCase := cresource_usecase.NewUseCaseImpl(cResourceRepo)
	expressionUseCase := expression_usecasse.NewUseCaseImpl(expressionRepo, cResourceRepo, cache)
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
		return err
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

	return nil
}

func createMigrate() *migrate.Migrate {
	m, err := migrate.New(
		"file://migrations/",
		os.Getenv("PGX_URL"),
	)

	if err != nil {
		log.Error().Msg(err.Error())
	}

	return m
}

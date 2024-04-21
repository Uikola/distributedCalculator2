package heartbeat

import (
	"context"
	"fmt"
	"math/rand/v2"
	"os"

	"github.com/Uikola/distributedCalculator2/orchestrator/internal/db"
	"github.com/Uikola/distributedCalculator2/orchestrator/internal/entity"
	"github.com/Uikola/distributedCalculator2/orchestrator/pkg/grpc/client/heartbeat"
	pb "github.com/Uikola/distributedCalculator2/proto/heartbeat"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type cResourceUseCase interface {
	Create(ctx context.Context, resource entity.CResource) error
	Exists(ctx context.Context, name, address string) (bool, error)
}

type Server struct {
	pb.HeartbeatServiceServer
	cResourceUseCase cResourceUseCase
}

func NewServer(cResourceUseCase cResourceUseCase) *Server {
	return &Server{cResourceUseCase: cResourceUseCase}
}

func (s *Server) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	log.Info().Msg(fmt.Sprintf("service %s sent register request", in.Name))

	database := db.InitDB("postgres://postgres:fgaSHFRdgkA4@localhost:5432/calcDB")

	port := rand.IntN(20000) + 30000
	addr := fmt.Sprintf("localhost:%d", port)

	calcConn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error().Err(err).Msg("failed to connect to the new client")
		os.Exit(1)
	}

	heartbeatCalcClient := heartbeat.NewClient(calcConn, database)
	log.Info().Msg(fmt.Sprintf("%s:%s", in.Name, addr))
	exists, err := s.cResourceUseCase.Exists(ctx, in.Name, addr)
	if err != nil {
		log.Error().Err(err).Msg("failed to check computing resource existion")
		return &pb.RegisterResponse{Success: false, Address: ""}, err
	}

	if exists {
		go func(conn *grpc.ClientConn, db *sqlx.DB) {
			defer conn.Close()
			defer db.Close()
			heartbeatCalcClient.SendHeartbeatToCalc(in.Name)
		}(calcConn, database)

		return &pb.RegisterResponse{Success: true}, nil
	}

	cResource := entity.CResource{
		Name:    in.Name,
		Address: addr,
	}
	err = s.cResourceUseCase.Create(ctx, cResource)
	if err != nil {
		log.Error().Err(err).Msg("failed to create new computing resource")
		return &pb.RegisterResponse{Success: false, Address: ""}, err
	}

	go func(conn *grpc.ClientConn, db *sqlx.DB) {
		defer conn.Close()
		defer db.Close()
		heartbeatCalcClient.SendHeartbeatToCalc(in.Name)
	}(calcConn, database)

	return &pb.RegisterResponse{Success: true, Address: addr}, nil
}

func (s *Server) Heartbeat(ctx context.Context, in *pb.HeartbeatRequest) (*pb.HeartbeatResponse, error) {
	log.Info().Msg(fmt.Sprintf("get heartbeat from %s", in.Heartbeat))
	return &pb.HeartbeatResponse{Heartbeat: true, Name: in.Heartbeat}, nil
}

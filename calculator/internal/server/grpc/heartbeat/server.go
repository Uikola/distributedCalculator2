package heartbeat

import (
	"context"
	"github.com/Uikola/distributedCalculator2/proto/heartbeat"
	"github.com/rs/zerolog/log"
)

type Server struct {
	heartbeat.HeartbeatServiceServer
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Heartbeat(ctx context.Context, in *heartbeat.HeartbeatRequest) (*heartbeat.HeartbeatResponse, error) {
	log.Info().Msg("get heartbeat from orchestrator")
	return &heartbeat.HeartbeatResponse{Heartbeat: true, Name: in.Heartbeat}, nil
}

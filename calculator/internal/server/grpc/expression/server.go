package expression

import (
	"context"
	"fmt"
	pb "github.com/Uikola/distributedCalculator2/proto/expression"
	"github.com/rs/zerolog/log"
)

type expressionUseCase interface {
	Calculate(ctx context.Context, expression string, expressionID uint) error
}

type Server struct {
	pb.ExpressionServiceServer
	expressionUseCase expressionUseCase
}

func NewServer(expressionUseCase expressionUseCase) *Server {
	return &Server{expressionUseCase: expressionUseCase}
}

func (s *Server) Calculate(ctx context.Context, in *pb.CalculateRequest) (*pb.CalculateResponse, error) {
	log.Info().Msg(fmt.Sprintf("get expression: %s", in.Expression))
	err := s.expressionUseCase.Calculate(ctx, in.Expression, uint(in.Id))
	if err != nil {
		log.Error().Err(err).Msg("failed to calculate expression")
		return &pb.CalculateResponse{}, err
	}

	return &pb.CalculateResponse{}, nil
}

package expression

import (
	"context"

	"github.com/Uikola/distributedCalculator2/orchestrator/internal/entity"
	pb "github.com/Uikola/distributedCalculator2/proto/expression"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type cResourceRepository interface {
	UnlinkExpressionFromCResource(ctx context.Context, expression entity.Expression) error
}

type expressionRepository interface {
	SetErrorStatus(ctx context.Context, id uint) error
}

type Client struct {
	client               pb.ExpressionServiceClient
	expressionRepository expressionRepository
	cResourceRepository  cResourceRepository
}

func NewClient(conn *grpc.ClientConn, expressionRepository expressionRepository, cResourceRepository cResourceRepository) *Client {
	client := pb.NewExpressionServiceClient(conn)
	return &Client{
		client:               client,
		expressionRepository: expressionRepository,
		cResourceRepository:  cResourceRepository,
	}
}

func (ec Client) Calculate(ctx context.Context, expression entity.Expression, conn *grpc.ClientConn) {
	go func() {
		defer conn.Close()
		_, err := ec.client.Calculate(context.Background(), &pb.CalculateRequest{Id: uint64(expression.ID), Expression: expression.Expression})
		if err != nil {
			log.Error().Err(err).Msg("failed to calculate expression")
			if unlinkErr := ec.cResourceRepository.UnlinkExpressionFromCResource(ctx, expression); unlinkErr != nil {
				log.Error().Msg(err.Error())
				return
			}
			if setErr := ec.expressionRepository.SetErrorStatus(ctx, expression.ID); setErr != nil {
				log.Error().Msg(err.Error())
				return
			}

			return
		}
	}()
}

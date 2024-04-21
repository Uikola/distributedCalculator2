package heartbeat

import (
	"context"
	"fmt"
	"time"

	"github.com/Uikola/distributedCalculator2/orchestrator/internal/db/repository/postgres"
	"github.com/Uikola/distributedCalculator2/orchestrator/internal/entity"
	"github.com/Uikola/distributedCalculator2/orchestrator/pkg/grpc/client/expression"
	pb "github.com/Uikola/distributedCalculator2/proto/heartbeat"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type cResourceRepository interface {
	SetOrchestatorHealth(ctx context.Context, name string, isAlive bool) error
	Delete(ctx context.Context, name string) error
	GetByName(ctx context.Context, name string) (entity.CResource, error)
	AssignExpressionToCResource(ctx context.Context, expression entity.Expression) (entity.CResource, error)
	UnlinkExpressionFromCResource(ctx context.Context, expression entity.Expression) error
}

type expressionRepository interface {
	GetByCResourceID(ctx context.Context, cResourceID uint) (entity.Expression, error)
	SetErrorStatus(ctx context.Context, id uint) error
	UpdateCResource(ctx context.Context, expressionID, cResourceID uint) error
}

type Client struct {
	client               pb.HeartbeatServiceClient
	cResourceRepository  cResourceRepository
	expressionRepository expressionRepository
}

func NewClient(conn *grpc.ClientConn, database *sqlx.DB) *Client {
	cResourceRepo := postgres.NewCResourceRepository(database)
	expressionRepo := postgres.NewExpressionRepository(database)
	client := pb.NewHeartbeatServiceClient(conn)

	return &Client{
		client:               client,
		cResourceRepository:  cResourceRepo,
		expressionRepository: expressionRepo,
	}
}

func (hc Client) Register(serviceName string) (bool, string, error) {
	response, err := hc.client.Register(context.Background(), &pb.RegisterRequest{
		Name: serviceName,
	})
	if err != nil {
		return false, "", err
	}

	return response.Success, response.Address, nil
}

func (hc Client) SendHeartbeatToOrch(serviceName string) {
	timer := time.NewTimer(30 * time.Second)
	for {
		log.Info().Msg(fmt.Sprintf("service %s sent heartbeat", serviceName))
		heartbeatTimer := time.NewTimer(5 * time.Second)
		isAlive, err := hc.client.Heartbeat(context.Background(), &pb.HeartbeatRequest{
			Heartbeat: serviceName,
		})
		if err != nil {
			log.Info().Err(err).Msg("failed to send heartbeat")
		}

		if isAlive != nil && isAlive.Heartbeat {
			err = hc.cResourceRepository.SetOrchestatorHealth(context.Background(), serviceName, true)
			if err != nil {
				log.Error().Err(err).Msg("failed to set orchestrator health")
			}
			timer = time.NewTimer(30 * time.Second)
		}

		select {
		case <-timer.C:
			err = hc.cResourceRepository.SetOrchestatorHealth(context.Background(), serviceName, false)
			if err != nil {
				log.Error().Err(err).Msg("failed to set orchestrator health")
			}
			timer = time.NewTimer(30 * time.Second)
			continue
		case <-heartbeatTimer.C:
			continue
		}
	}
}

func (hc Client) SendHeartbeatToCalc(serviceName string) {
	ctx := context.Background()

	timer := time.NewTimer(40 * time.Second)
	for {

		heartbeatTimer := time.NewTimer(5 * time.Second)
		_, err := hc.client.Heartbeat(context.Background(), &pb.HeartbeatRequest{
			Heartbeat: serviceName,
		})

		if err != nil {
			log.Info().Err(err).Msg("failed to send heartbeat")
		} else {
			timer = time.NewTimer(40 * time.Second)
		}

		select {
		case <-timer.C:
			cResource, err := hc.cResourceRepository.GetByName(ctx, serviceName)
			if err != nil {
				log.Error().Msg(err.Error())
			}

			if cResource.Occupied {
				expr, err := hc.expressionRepository.GetByCResourceID(ctx, cResource.ID)
				if err != nil {
					log.Error().Msg(err.Error())
				}
				log.Info().Msg(fmt.Sprintf("выражение %s передаётся с упавшей машины %s", expr.Expression, cResource.Name))

				cResource, err = hc.cResourceRepository.AssignExpressionToCResource(ctx, expr)
				if err != nil {
					log.Error().Msg(err.Error())
				}
				log.Info().Msg(fmt.Sprintf("выражение %s передано на машину %s", expr.Expression, cResource.Name))

				err = hc.expressionRepository.UpdateCResource(ctx, expr.ID, cResource.ID)
				if err != nil {
					log.Error().Msg(err.Error())
				}

				calcConn, err := grpc.NewClient(cResource.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
				if err != nil {

				}

				log.Info().Msg(fmt.Sprintf("машина %s начала вычислять выражение %s", cResource.Name, expr.Expression))
				expressionClient := expression.NewClient(calcConn, hc.expressionRepository, hc.cResourceRepository)
				expressionClient.Calculate(ctx, expr, calcConn)
			}

			err = hc.cResourceRepository.Delete(ctx, serviceName)
			if err != nil {
				log.Error().Err(err).Msg("failed to delete calculator service")
			}
			return
		case <-heartbeatTimer.C:
			continue
		}
	}
}

package postgres

import (
	"context"
	"fmt"
	"github.com/Uikola/distributedCalculator2/calculator/internal/entity"
	"github.com/jmoiron/sqlx"
)

type CResourceRepository struct {
	db *sqlx.DB
}

func NewCResourceRepository(db *sqlx.DB) *CResourceRepository {
	return &CResourceRepository{db: db}
}

func (r CResourceRepository) UnlinkExpressionFromCResource(ctx context.Context, expression entity.Expression) error {
	const op = "CResourceRepository.UnlinkExpressionFromCResource"

	stmt, err := r.db.PreparexContext(ctx, "UPDATE cresources SET occupied = false, expression = '' WHERE id = $1")
	if err != nil {
		return fmt.Errorf("%s:%v", op, err)
	}

	_, err = stmt.ExecContext(ctx, expression.CalculatedBy)
	if err != nil {
		return fmt.Errorf("%s:%v", op, err)
	}

	return nil
}

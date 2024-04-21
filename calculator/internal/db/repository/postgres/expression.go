package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/Uikola/distributedCalculator2/calculator/internal/entity"
	"github.com/jmoiron/sqlx"
)

type ExpressionRepository struct {
	db *sqlx.DB
}

func NewExpressionRepository(db *sqlx.DB) *ExpressionRepository {
	return &ExpressionRepository{db: db}
}

func (r ExpressionRepository) UpdateResult(ctx context.Context, expressionID uint, result string) error {
	const op = "ExpressionRepository.UpdateResult"

	stmt, err := r.db.PreparexContext(ctx, "UPDATE expressions SET result = $1 WHERE id = $2")
	if err != nil {
		return fmt.Errorf("%s:%v", op, err)
	}

	_, err = stmt.ExecContext(ctx, result, expressionID)
	if err != nil {
		return fmt.Errorf("%s:%v", op, err)
	}

	return nil
}

func (r ExpressionRepository) GetExpressionByID(ctx context.Context, id uint) (entity.Expression, error) {
	const op = "ExpressionRepository.GetExpressionByID"

	stmt, err := r.db.PreparexContext(ctx, "SELECT id, expression, calculated_by, owner_id FROM expressions WHERE id = $1")
	if err != nil {
		return entity.Expression{}, fmt.Errorf("%s:%v", op, err)
	}

	row := stmt.QueryRowxContext(ctx, id)

	var expression string
	var calculatedBy, ownerID uint
	err = row.Scan(&id, &expression, &calculatedBy, &ownerID)
	if err != nil {
		return entity.Expression{}, fmt.Errorf("%s:%v", op, err)
	}

	return entity.Expression{
		ID:           id,
		Expression:   expression,
		CalculatedBy: calculatedBy,
		OwnerID:      ownerID,
	}, nil
}

func (r ExpressionRepository) SetErrorStatus(ctx context.Context, id uint) error {
	const op = "ExpressionRepository.SetErrorStatus"

	stmt, err := r.db.PreparexContext(ctx, "UPDATE expressions SET status = $1, calculated_at = $2 WHERE id = $3")
	if err != nil {
		return fmt.Errorf("%s:%v", op, err)
	}

	_, err = stmt.ExecContext(ctx, entity.Error, time.Now(), id)
	if err != nil {
		return fmt.Errorf("%s:%v", op, err)
	}

	return nil
}

func (r ExpressionRepository) SetSuccessStatus(ctx context.Context, id uint) error {
	const op = "ExpressionRepository.SetSuccessStatus"

	stmt, err := r.db.PreparexContext(ctx, "UPDATE expressions SET status = $1, calculated_at = $2 WHERE id = $3")
	if err != nil {
		return fmt.Errorf("%s:%v", op, err)
	}

	_, err = stmt.ExecContext(ctx, entity.OK, time.Now(), id)
	if err != nil {
		return fmt.Errorf("%s:%v", op, err)
	}

	return nil
}

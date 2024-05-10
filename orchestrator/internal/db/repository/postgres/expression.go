package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Uikola/distributedCalculator2/orchestrator/internal/entity"
	"github.com/Uikola/distributedCalculator2/orchestrator/internal/errorz"
	"github.com/jmoiron/sqlx"
)

type ExpressionRepository struct {
	db *sqlx.DB
}

func NewExpressionRepository(db *sqlx.DB) *ExpressionRepository {
	return &ExpressionRepository{db: db}
}

func (r ExpressionRepository) AddExpression(ctx context.Context, expression entity.Expression) (entity.Expression, error) {
	const op = "ExpressionRepository.AddExpression"

	stmt, err := r.db.PreparexContext(ctx, "INSERT INTO expressions(expression, status, created_at, calculated_by, owner_id) VALUES ($1, $2, $3, $4, $5) RETURNING id, expression, status, created_at, owner_id")
	if err != nil {
		return entity.Expression{}, fmt.Errorf("%s:%v", op, err)
	}

	row := stmt.QueryRowxContext(ctx, expression.Expression, expression.Status, expression.CreatedAt, expression.CalculatedBy, expression.OwnerID)

	var id, ownerID uint
	var expr, status string
	var createdAt time.Time
	err = row.Scan(&id, &expr, &status, &createdAt, &ownerID)
	if err != nil {
		return entity.Expression{}, fmt.Errorf("%s:%v", op, err)
	}

	return entity.Expression{
		ID:         id,
		Expression: expr,
		Status:     entity.ExpressionStatus(status),
		CreatedAt:  createdAt,
		OwnerID:    ownerID,
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

func (r ExpressionRepository) UpdateCResource(ctx context.Context, expressionID, cResourceID uint) error {
	const op = "ExpressionRepository.UpdateCResource"

	stmt, err := r.db.PreparexContext(ctx, "UPDATE expressions SET calculated_by = $1 WHERE id = $2")
	if err != nil {
		return fmt.Errorf("%s:%v", op, err)
	}

	_, err = stmt.ExecContext(ctx, cResourceID, expressionID)
	if err != nil {
		return fmt.Errorf("%s:%v", op, err)
	}

	return nil
}

func (r ExpressionRepository) GetExpressionByID(ctx context.Context, id uint) (entity.Expression, error) {
	const op = "ExpressionRepository.GetExpressionByID"

	stmt, err := r.db.PreparexContext(ctx, "SELECT id, expression, result, status, created_at, calculated_at, calculated_by, owner_id FROM expressions WHERE id = $1")
	if err != nil {
		return entity.Expression{}, fmt.Errorf("%s:%v", op, err)
	}

	row := stmt.QueryRowxContext(ctx, id)

	var calculatedBy, ownerID uint
	var expression, result, status string
	var createdAt time.Time
	var calculatedAt sql.NullTime
	err = row.Scan(&id, &expression, &result, &status, &createdAt, &calculatedAt, &calculatedBy, &ownerID)
	switch {
	case errors.Is(sql.ErrNoRows, err):
		return entity.Expression{}, errorz.ErrExpressionNotFound
	case err != nil:
		return entity.Expression{}, fmt.Errorf("%s:%v", op, err)
	}

	expr := entity.Expression{
		ID:           id,
		Expression:   expression,
		Result:       result,
		Status:       entity.ExpressionStatus(status),
		CreatedAt:    createdAt,
		CalculatedBy: calculatedBy,
		OwnerID:      ownerID,
	}

	if calculatedAt.Valid {
		expr.CalculatedAt = calculatedAt.Time
	}

	return expr, nil
}

func (r ExpressionRepository) ListExpressions(ctx context.Context, userID uint) ([]entity.Expression, error) {
	const op = "ExpressionRepository.ListExpressions"

	stmt, err := r.db.PreparexContext(ctx, "SELECT id, expression, result, status, created_at, calculated_at, calculated_by, owner_id FROM expressions WHERE owner_id = $1")
	if err != nil {
		return []entity.Expression{}, fmt.Errorf("%s:%v", op, err)
	}

	rows, err := stmt.QueryxContext(ctx, userID)
	if err != nil {
		return []entity.Expression{}, fmt.Errorf("%s:%v", op, err)
	}

	var id, calculatedBy, ownerID uint
	var expr, result, status string
	var createdAt time.Time
	var calculatedAt sql.NullTime
	var expressions []entity.Expression

	for rows.Next() {
		err = rows.Scan(&id, &expr, &result, &status, &createdAt, &calculatedAt, &calculatedBy, &ownerID)
		if err != nil {
			return []entity.Expression{}, fmt.Errorf("%s:%v", op, err)
		}
		expression := entity.Expression{
			ID:           id,
			Expression:   expr,
			Result:       result,
			Status:       entity.ExpressionStatus(status),
			CreatedAt:    createdAt,
			CalculatedBy: calculatedBy,
			OwnerID:      ownerID,
		}
		if calculatedAt.Valid {
			expression.CalculatedAt = calculatedAt.Time
		}

		expressions = append(expressions, expression)
	}

	return expressions, nil
}

func (r ExpressionRepository) GetByCResourceID(ctx context.Context, cResourceID uint) (entity.Expression, error) {
	const op = "ExpressionRepository.GetByCResourceID"

	stmt, err := r.db.PreparexContext(ctx, "SELECT id, expression, result, status, created_at, calculated_at, calculated_by, owner_id FROM expressions WHERE calculated_by = $1")
	if err != nil {
		return entity.Expression{}, fmt.Errorf("%s:%v", op, err)
	}

	row := stmt.QueryRowxContext(ctx, cResourceID)

	var id, calculatedBy, ownerID uint
	var expr, result, status string
	var createdAt time.Time
	var calculatedAt sql.NullTime
	err = row.Scan(&id, &expr, &result, &status, &createdAt, &calculatedAt, &calculatedBy, &ownerID)
	if err != nil {
		return entity.Expression{}, fmt.Errorf("%s:%v", op, err)
	}

	expression := entity.Expression{
		ID:           id,
		Expression:   expr,
		Result:       result,
		Status:       entity.ExpressionStatus(status),
		CreatedAt:    createdAt,
		CalculatedBy: calculatedBy,
		OwnerID:      ownerID,
	}
	if calculatedAt.Valid {
		expression.CalculatedAt = calculatedAt.Time
	}

	return expression, nil
}

func (r ExpressionRepository) CleanUp(ctx context.Context) error {
	const op = "ExpressionRepository.CleanUp"

	stmt, err := r.db.PreparexContext(ctx, "DELETE FROM expressions")
	if err != nil {
		return fmt.Errorf("%s:%v", op, err)
	}

	_, err = stmt.ExecContext(ctx)
	if err != nil {
		return fmt.Errorf("%s:%v", op, err)
	}

	return nil
}

package postgres

import (
	"context"
	"fmt"

	"github.com/Uikola/distributedCalculator2/calculator/internal/entity"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r UserRepository) GetByID(ctx context.Context, id uint) (entity.User, error) {
	const op = "UserRepository.GetByID"

	stmt, err := r.db.PreparexContext(ctx, "SELECT id, login, password, addition, subtraction, multiplication, division FROM users WHERE id = $1")
	if err != nil {
		return entity.User{}, fmt.Errorf("%s:%v", op, err)
	}

	row := stmt.QueryRowxContext(ctx, id)

	var login, password string
	var addition, subtraction, multiplication, division uint
	err = row.Scan(&id, &login, &password, &addition, &subtraction, &multiplication, &division)
	if err != nil {
		return entity.User{}, fmt.Errorf("%s:%v", op, err)
	}

	return entity.User{
		ID:             id,
		Login:          login,
		Password:       password,
		Addition:       addition,
		Subtraction:    subtraction,
		Multiplication: multiplication,
		Division:       division,
	}, nil
}

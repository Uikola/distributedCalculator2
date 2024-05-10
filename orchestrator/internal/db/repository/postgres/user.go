package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Uikola/distributedCalculator2/orchestrator/internal/entity"
	"github.com/Uikola/distributedCalculator2/orchestrator/internal/errorz"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r UserRepository) Create(ctx context.Context, user entity.User) error {
	const op = "UserRepository.Create"

	stmt, err := r.db.PreparexContext(ctx, "INSERT INTO users(login, password) VALUES ($1, $2)")
	if err != nil {
		return fmt.Errorf("%s:%v", op, err)
	}

	_, err = stmt.ExecContext(ctx, user.Login, user.Password)
	if err != nil {
		return fmt.Errorf("%s:%v", op, err)
	}

	return nil
}

func (r UserRepository) Exists(ctx context.Context, login string) (bool, error) {
	const op = "UserRepository.Exists"

	stmt, err := r.db.PreparexContext(ctx, "SELECT login FROM users WHERE login = $1")
	if err != nil {
		return false, fmt.Errorf("%s:%v", op, err)
	}

	row := stmt.QueryRowxContext(ctx, login)
	err = row.Scan(&login)
	switch {
	case errors.Is(sql.ErrNoRows, err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("%s:%v", op, err)
	}

	return true, nil
}

func (r UserRepository) GetUser(ctx context.Context, login, password string) (entity.User, error) {
	const op = "UserRepository.GetUser"

	stmt, err := r.db.PreparexContext(ctx, "SELECT id, login, password FROM users WHERE login = $1 AND password = $2")
	if err != nil {
		return entity.User{}, fmt.Errorf("%s:%v", op, err)
	}

	row := stmt.QueryRowxContext(ctx, login, password)

	var id uint
	err = row.Scan(&id, &login, &password)
	switch {
	case errors.Is(sql.ErrNoRows, err):
		return entity.User{}, errorz.ErrUserNotFound
	case err != nil:
		return entity.User{}, fmt.Errorf("%s:%v", op, err)
	}

	return entity.User{
		ID:       id,
		Login:    login,
		Password: password,
	}, nil
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

func (r UserRepository) UpdateOperation(ctx context.Context, userID uint, operation string, operationTime uint) error {
	const op = "UserRepository.UpdateOperation"

	stmt, err := r.db.PreparexContext(ctx, fmt.Sprintf("UPDATE users SET %s = $1 WHERE id = $2", operation))
	if err != nil {
		return fmt.Errorf("%s:%v", op, err)
	}

	_, err = stmt.ExecContext(ctx, operationTime, userID)
	if err != nil {
		return fmt.Errorf("%s:%v", op, err)
	}

	return nil
}

func (r UserRepository) CleanUp(ctx context.Context) error {
	const op = "UserRepository.CleanUp"

	stmt, err := r.db.PreparexContext(ctx, "DELETE FROM users")
	if err != nil {
		return fmt.Errorf("%s:%v", op, err)
	}

	_, err = stmt.ExecContext(ctx)
	if err != nil {
		return fmt.Errorf("%s:%v", op, err)
	}

	return nil
}

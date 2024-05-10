package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/rand"

	"github.com/Uikola/distributedCalculator2/orchestrator/internal/entity"
	"github.com/Uikola/distributedCalculator2/orchestrator/internal/errorz"
	"github.com/jmoiron/sqlx"
)

type CResourceRepository struct {
	db *sqlx.DB
}

func NewCResourceRepository(db *sqlx.DB) *CResourceRepository {
	return &CResourceRepository{db: db}
}

func (r CResourceRepository) Create(ctx context.Context, resource entity.CResource) error {
	const op = "CResourceRepository.Create"

	stmt, err := r.db.PreparexContext(ctx, "INSERT INTO cresources(name, address) VALUES ($1, $2)")
	if err != nil {
		return fmt.Errorf("%s:%v", op, err)
	}

	_, err = stmt.ExecContext(ctx, resource.Name, resource.Address)
	if err != nil {
		return fmt.Errorf("%s:%v", op, err)
	}

	return nil
}

//func (r CResourceRepository) GetCResourceByID(ctx context.Context, id uint) (entity.CResource, error) {
//	const op = "CResourceRepository.GetCResourceByID"
//
//	stmt, err := r.db.PreparexContext(ctx, "SELECT id, name, address, expression, occupied, orchestrator_alive FROM cresources WHERE id = $1")
//	if err != nil {
//		return entity.CResource{}, fmt.Errorf("%s:%v",op, err)
//	}
//
//	row := stmt.QueryRowxContext(ctx, id)
//
//	var name, address, expression string
//	var occupied, orchestratorAlive bool
//	err = row.Scan(&id, &name, &address, &expression, &occupied, &orchestratorAlive)
//	if err != nil {
//		return entity.CResource{}, err
//	}
//
//	return entity.CResource{
//		ID:                id,
//		Name:              name,
//		Address:           address,
//		Expression:        expression,
//		Occupied:          occupied,
//		OrchestratorAlive: orchestratorAlive,
//	}, nil
//}

func (r CResourceRepository) Exists(ctx context.Context, name, address string) (bool, error) {
	const op = "CResourceRepository.Exists"

	stmt, err := r.db.PreparexContext(ctx, "SELECT name FROM cresources WHERE name = $1 OR address = $2")
	if err != nil {
		return false, fmt.Errorf("%s:%v", op, err)
	}

	row := stmt.QueryRowxContext(ctx, name, address)

	err = row.Scan(&name)
	switch {
	case errors.Is(sql.ErrNoRows, err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("%s:%v", op, err)
	}

	return true, nil
}

func (r CResourceRepository) SetOrchestatorHealth(ctx context.Context, name string, isAlive bool) error {
	const op = "CResourceRepository.SetOrchestatorHealth"

	stmt, err := r.db.PreparexContext(ctx, "UPDATE cresources SET orchestrator_alive = $1 WHERE name = $2")
	if err != nil {
		return fmt.Errorf("%s:%v", op, err)
	}

	_, err = stmt.ExecContext(ctx, isAlive, name)
	if err != nil {
		return fmt.Errorf("%s:%v", op, err)
	}

	return nil
}

func (r CResourceRepository) Delete(ctx context.Context, name string) error {
	const op = "CResourceRepository.Delete"

	stmt, err := r.db.PreparexContext(ctx, "DELETE FROM cresources WHERE name = $1")
	if err != nil {
		return fmt.Errorf("%s:%v", op, err)
	}

	_, err = stmt.ExecContext(ctx, name)
	if err != nil {
		return fmt.Errorf("%s:%v", op, err)
	}

	return nil
}

func (r CResourceRepository) AssignExpressionToCResource(ctx context.Context, expression entity.Expression) (entity.CResource, error) {
	const op = "CResourceRepository.AssignExpressionToCResource"

	tx, err := r.db.BeginTxx(ctx, &sql.TxOptions{
		Isolation: sql.LevelDefault,
	})
	if err != nil {
		tx.Rollback()
		return entity.CResource{}, fmt.Errorf("%s:%v", op, err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	stmt, err := tx.PreparexContext(ctx, "SELECT id, name, address, expression, occupied, orchestrator_alive FROM cresources WHERE occupied = false")
	if err != nil {
		return entity.CResource{}, fmt.Errorf("%s:%v", op, err)
	}

	rows, err := stmt.QueryxContext(ctx)
	if err != nil {
		return entity.CResource{}, fmt.Errorf("%s:%v", op, err)
	}

	var cResources []entity.CResource
	for rows.Next() {
		var id uint
		var name, address, expr string
		var occupied, orchestratorAlive bool

		err = rows.Scan(&id, &name, &address, &expr, &occupied, &orchestratorAlive)
		if err != nil {
			return entity.CResource{}, fmt.Errorf("%s:%v", op, err)
		}
		cResource := entity.CResource{
			ID:                id,
			Name:              name,
			Address:           address,
			Expression:        expr,
			Occupied:          occupied,
			OrchestratorAlive: orchestratorAlive,
		}
		cResources = append(cResources, cResource)
	}

	if len(cResources) == 0 {
		return entity.CResource{}, errorz.ErrNoAvailableResources
	}
	cResource := cResources[rand.Intn(len(cResources))]

	stmt, err = tx.PreparexContext(ctx, "UPDATE cresources SET occupied = true, expression = $1 WHERE id = $2")
	if err != nil {
		return entity.CResource{}, fmt.Errorf("%s:%v", op, err)
	}

	_, err = stmt.ExecContext(ctx, expression.Expression, cResource.ID)
	if err != nil {
		return entity.CResource{}, fmt.Errorf("%s:%v", op, err)
	}
	tx.Commit()
	cResource.Occupied = true
	cResource.Expression = expression.Expression

	return cResource, nil
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

func (r CResourceRepository) ListCResources(ctx context.Context) ([]entity.CResource, error) {
	const op = "CResourceRepository.ListCResources"

	stmt, err := r.db.PreparexContext(ctx, "SELECT id, name, address, expression, occupied, orchestrator_alive FROM cresources")
	if err != nil {
		return []entity.CResource{}, fmt.Errorf("%s:%v", op, err)
	}

	rows, err := stmt.QueryxContext(ctx)
	if err != nil {
		return []entity.CResource{}, fmt.Errorf("%s:%v", op, err)
	}

	var id uint
	var name, address, expression string
	var occupied, orchestratorAlive bool
	var cResources []entity.CResource

	for rows.Next() {
		err = rows.Scan(&id, &name, &address, &expression, &occupied, &orchestratorAlive)
		if err != nil {
			return []entity.CResource{}, fmt.Errorf("%s:%v", op, err)
		}
		cResource := entity.CResource{
			ID:                id,
			Name:              name,
			Address:           address,
			Expression:        expression,
			Occupied:          occupied,
			OrchestratorAlive: orchestratorAlive,
		}
		cResources = append(cResources, cResource)
	}

	return cResources, nil
}

func (r CResourceRepository) GetByName(ctx context.Context, name string) (entity.CResource, error) {
	const op = "CResourceRepository.GetByName"

	stmt, err := r.db.PreparexContext(ctx, "SELECT id, name, address, expression, occupied, orchestrator_alive FROM cresources WHERE name = $1")
	if err != nil {
		return entity.CResource{}, fmt.Errorf("%s:%v", op, err)
	}

	row := stmt.QueryRowxContext(ctx, name)

	var id uint
	var address, expression string
	var occupied, orchestratorAlive bool
	err = row.Scan(&id, &name, &address, &expression, &occupied, &orchestratorAlive)
	if err != nil {
		return entity.CResource{}, fmt.Errorf("%s:%v", op, err)
	}

	return entity.CResource{
		ID:                id,
		Name:              name,
		Address:           address,
		Expression:        expression,
		Occupied:          occupied,
		OrchestratorAlive: orchestratorAlive,
	}, nil
}

func (r CResourceRepository) CleanUp(ctx context.Context) error {
	const op = "CResourceRepository.CleanUp"

	stmt, err := r.db.PreparexContext(ctx, "DELETE FROM cresources")
	if err != nil {
		return fmt.Errorf("%s:%v", op, err)
	}

	_, err = stmt.ExecContext(ctx)
	if err != nil {
		return fmt.Errorf("%s:%v", op, err)
	}

	return nil
}

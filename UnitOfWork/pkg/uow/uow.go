package uow

import (
	"context"
	"database/sql"
	"fmt"
)

type RepositoryFactory func(tx *sql.Tx) any

type UnitOfWork interface {
	Register(name string, factory RepositoryFactory)
	Get(ctx context.Context, name string) (any, error)
	Do(ctx context.Context, fn func(uow UnitOfWork) error) error
	CommitOrRollback() error
	Rollback() error
	UnRegister(name string)
}

type UnitOfWorkImpl struct {
	db           *sql.DB
	repositories map[string]RepositoryFactory
	tx           *sql.Tx
}

func NewUnitOfWork(ctx context.Context, db *sql.DB) UnitOfWork {
	return &UnitOfWorkImpl{
		db:           db,
		repositories: make(map[string]RepositoryFactory),
	}
}

func (u *UnitOfWorkImpl) Register(name string, factory RepositoryFactory) {
	u.repositories[name] = factory
}

func (u *UnitOfWorkImpl) UnRegister(name string) {
	delete(u.repositories, name)
}

func (u *UnitOfWorkImpl) Do(ctx context.Context, fn func(uow UnitOfWork) error) error {
	if u.tx != nil {
		return fmt.Errorf("unit of work already started")
	}
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	u.tx = tx
	err = fn(u)
	if err != nil {
		if errRb := u.Rollback(); errRb != nil {
			return fmt.Errorf("rollback failed: %v, original error: %v", errRb, err)
		}
		return err
	}
	return u.CommitOrRollback()
}

func (u *UnitOfWorkImpl) Get(ctx context.Context, name string) (any, error) {
	if u.tx == nil {
		tx, err := u.db.BeginTx(ctx, nil)
		if err != nil {
			return nil, err
		}
		u.tx = tx
	}
	return u.repositories[name](u.tx), nil
}

func (u *UnitOfWorkImpl) CommitOrRollback() error {
	if u.tx == nil {
		return fmt.Errorf("unit of work not started")
	}
	err := u.tx.Commit()
	if err != nil {
		if errRb := u.Rollback(); errRb != nil {
			return fmt.Errorf("rollback failed: %v, original error: %v", errRb, err)
		}
		return err
	}
	u.tx = nil
	return nil
}

func (u *UnitOfWorkImpl) Rollback() error {
	if u.tx == nil {
		return fmt.Errorf("unit of work not started")
	}
	err := u.tx.Rollback()
	if err != nil {
		return err
	}
	u.tx = nil
	return nil
}

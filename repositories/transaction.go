package repositories

import (
	"context"
	"github.com/hotrungnhan/surl/utils/injects"
	"errors"

	"github.com/samber/do/v2"
	"gorm.io/gorm"
)

/*
Best use case

err = d.txRepo.Begin(context) -> inject db into context

	defer func() {
		err = d.txRepo.Exec(context, err)  -> try commit if no issue found, reject if there any issue
	}()
*/

var UninitializeTransactionError = errors.New("Transaction isn't initialize yet")

type TransactionRepository interface {
	Begin(context.Context) (context.Context, error)
	Exec(context.Context, error) error
	Commit(context.Context) error
	Rollback(context.Context) error
}

type transactionRepositoryImpl struct {
	db *injects.DB
}

func (r *transactionRepositoryImpl) Begin(ctx context.Context) (context.Context, error) {
	db := r.db.Master.Begin()

	if db.Error != nil {
		return nil, db.Error
	}

	ctx = context.WithValue(ctx, "db", db)
	return ctx, nil
}

func (r *transactionRepositoryImpl) Commit(ctx context.Context) error {
	if ctx.Value("db") == nil {
		return UninitializeTransactionError
	}
	err := ctx.Value("db").(*gorm.DB).Commit().Error

	return err
}

func (r *transactionRepositoryImpl) Rollback(ctx context.Context) error {
	if ctx.Value("db") == nil {
		return UninitializeTransactionError
	}
	err := ctx.Value("db").(*gorm.DB).Commit().Error

	return err
}

func (r *transactionRepositoryImpl) Exec(ctx context.Context, trackerError error) error {
	if ctx.Value("db") == nil {
		return UninitializeTransactionError
	}
	if trackerError != nil {
		txErr := r.Rollback(ctx)
		if txErr != nil {
			trackerError = errors.Join(txErr, trackerError)
		}
	}

	trackerError = r.Commit(ctx)

	return trackerError
}

func NewTransactionRepository(i do.Injector) (TransactionRepository, error) {
	return &transactionRepositoryImpl{
		db: do.MustInvoke[*injects.DB](i),
	}, nil
}

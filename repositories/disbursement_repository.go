package repositories

import (
	"context"

	"github.com/peang/amartha-loan-service/models"
	"github.com/uptrace/bun"
)

type DisbursementRepositoryInterface interface {
	Save(ctx context.Context, loan *models.Disbursment) (*models.Disbursment, error)
}

type disbursementRepository struct {
	db *bun.DB
}

func NewDisbursementRepository(db *bun.DB) DisbursementRepositoryInterface {
	return &disbursementRepository{
		db: db,
	}
}

func (r *disbursementRepository) Save(ctx context.Context, disbursement *models.Disbursment) (*models.Disbursment, error) {
	_, err := r.db.NewInsert().Model(disbursement).On("CONFLICT (id) DO UPDATE").Returning("id").Exec(ctx)
	if err != nil {
		return nil, err
	}

	return disbursement, nil
}

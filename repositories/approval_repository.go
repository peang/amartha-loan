package repositories

import (
	"context"

	"github.com/peang/amartha-loan-service/models"
	"github.com/uptrace/bun"
)

type ApprovalRepositoryInterface interface {
	Save(ctx context.Context, loan *models.Approval) (*models.Approval, error)
}

type approvalRepository struct {
	db *bun.DB
}

func NewApprovalRepository(db *bun.DB) ApprovalRepositoryInterface {
	return &approvalRepository{
		db: db,
	}
}

func (r *approvalRepository) Save(ctx context.Context, approval *models.Approval) (*models.Approval, error) {
	_, err := r.db.NewInsert().Model(approval).On("CONFLICT (id) DO UPDATE").Returning("id").Exec(ctx)
	if err != nil {
		return nil, err
	}

	return approval, nil
}

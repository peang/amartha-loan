package repositories

import (
	"context"
	"database/sql"

	"github.com/gotidy/ptr"
	"github.com/peang/amartha-loan-service/models"
	"github.com/peang/amartha-loan-service/utils"
	"github.com/uptrace/bun"
)

type LoanRepositoryInterface interface {
	List(ctx context.Context, page int, perPage int, sort string, filter LoanRepositoryFilter) (*[]models.Loan, int, error)
	Save(tx *bun.Tx, ctx context.Context, loan *models.Loan) (*models.Loan, error)
	Detail(ctx context.Context, id string) (loan *models.Loan, err error)
}

type LoanRepositoryFilter struct {
	Status *models.LoanStatus
}

type loanRepository struct {
	approvalRepository ApprovalRepositoryInterface
	db                 *bun.DB
}

func NewLoanRepository(db *bun.DB, approvalRepository ApprovalRepositoryInterface) LoanRepositoryInterface {
	return &loanRepository{
		approvalRepository: approvalRepository,
		db:                 db,
	}
}

func (r *loanRepository) Save(tx *bun.Tx, ctx context.Context, loan *models.Loan) (*models.Loan, error) {
	if tx == nil {
		trx, _ := r.db.Begin()
		tx = ptr.Of(trx)
	}

	if loan.Status == models.LoanStatusApproved && loan.ApprovalID == nil {
		approval := loan.Approval

		_, err := r.approvalRepository.Save(ctx, approval)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		loan.ApprovalID = &approval.ID
	}

	_, err := r.db.NewInsert().Model(loan).On("CONFLICT (id) DO UPDATE").Returning("id").Exec(ctx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if tx != nil {
		tx.Commit()
	}

	return loan, nil
}

func (r *loanRepository) Detail(ctx context.Context, uuid string) (*models.Loan, error) {
	var loan models.Loan
	err := r.db.NewSelect().Model(&loan).Where("uuid = ?", uuid).Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &loan, nil
}

func (r *loanRepository) List(ctx context.Context, page int, perPage int, sort string, filter LoanRepositoryFilter) (*[]models.Loan, int, error) {
	sorts := utils.GenerateSort(sort)
	offset, limit := utils.GenerateOffsetLimit(page, perPage)

	var loans []models.Loan
	sl := r.db.NewSelect().Model(&loans)
	if filter.Status != nil {
		sl.Where("? = ?", bun.Ident("status"), filter.Status)
	}

	count, err := sl.Limit(limit).Offset(offset).OrderExpr(sorts).ScanAndCount(context.TODO())
	if err != nil {
		return nil, 0, err
	}

	if len(loans) == 0 {
		return &[]models.Loan{}, count, nil
	}

	return &loans, count, nil
}

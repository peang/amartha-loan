package repositories

import (
	"context"

	"github.com/peang/amartha-loan-service/models"
	"github.com/peang/amartha-loan-service/utils"
	"github.com/uptrace/bun"
)

type InvestmentRepositoryInterface interface {
	List(ctx context.Context, page int, perPage int, sort string, filter InvestmentRepositoryFilter) (*[]models.Investment, int, error)
	Save(ctx context.Context, loan *models.Investment) (*models.Investment, error)
	UpdateMany(ctx context.Context, filter InvestmentRepositoryFilter, value InvestmentRepositoryValues) error
}

type InvestmentRepositoryFilter struct {
	LoanID *uint
}
type InvestmentRepositoryValues struct {
	SendAggreementEmail *bool
}

type investmentRepository struct {
	loanRepository LoanRepositoryInterface
	db             *bun.DB
}

func NewInvestmentRepository(db *bun.DB, loanRepository LoanRepositoryInterface) InvestmentRepositoryInterface {
	return &investmentRepository{
		loanRepository: loanRepository,
		db:             db,
	}
}

func (r *investmentRepository) Save(ctx context.Context, investment *models.Investment) (*models.Investment, error) {
	tx, _ := r.db.Begin()

	_, err := r.db.NewInsert().Model(investment).On("CONFLICT (id) DO UPDATE").Returning("id").Exec(ctx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	_, err = r.loanRepository.Save(&tx, ctx, investment.Loan)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return investment, nil
}

func (r *investmentRepository) List(ctx context.Context, page int, perPage int, sort string, filter InvestmentRepositoryFilter) (*[]models.Investment, int, error) {
	offset, limit := utils.GenerateOffsetLimit(page, perPage)

	var investments []models.Investment
	sl := r.db.NewSelect().Model(&investments)
	sl.Relation("Investor")
	if filter.LoanID != nil {
		sl.Where("? = ?", bun.Ident("loan_id"), filter.LoanID)
	}

	count, err := sl.Group("investment.investor_id", "investor.id").Column("investment.investor_id").Limit(limit).Offset(offset).ScanAndCount(context.TODO())
	if err != nil {
		return nil, 0, err
	}

	if len(investments) == 0 {
		return &[]models.Investment{}, count, nil
	}

	return &investments, count, nil
}

func (r *investmentRepository) UpdateMany(ctx context.Context, filter InvestmentRepositoryFilter, value InvestmentRepositoryValues) error {
	investments := models.Investment{}

	sl := r.db.NewUpdate().Model(&investments)
	if filter.LoanID != nil {
		sl.Where("? = ?", bun.Ident("loan_id"), filter.LoanID)
	}

	if value.SendAggreementEmail != nil {
		sl.Set("send_aggreement_email = ?", value.SendAggreementEmail)
	}

	_, err := sl.Exec(context.TODO())
	if err != nil {
		return err
	}
	return nil
}

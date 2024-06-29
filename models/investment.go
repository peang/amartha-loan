package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Investment struct {
	bun.BaseModel `bun:"table:investments"`

	ID                  uint       `bun:"id,pk,nullzero"`
	LoanID              uint       `bun:"loan_id"`
	InvestorID          uint       `bun:"investor_id"`
	Amount              float64    `bun:"amount"`
	ROI                 float64    `bun:"roi"`
	SendAggreementEmail bool       `bun:"send_aggreement_email"`
	CreatedAt           time.Time  `bun:"created_at"`
	UpdatedAt           *time.Time `bun:"updated_at,nullzero"`

	Loan     *Loan `bun:"rel:has-one,join:loan_id=id"`
	Investor *User `bun:"rel:has-one,join:investor_id=id"`
}

func NewInvestment(investorId uint, amount float64, loan *Loan) (*Investment, error) {
	err := loan.Invest(amount)
	if err != nil {
		return nil, err
	}

	return &Investment{
		LoanID:     loan.ID,
		InvestorID: investorId,
		Amount:     amount,
		ROI:        (amount / loan.ProposedAmount) * loan.ROI,
		Loan:       loan,
	}, nil
}

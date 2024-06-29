package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type LoanStatus int

const (
	LoanStatusProposed LoanStatus = iota
	LoanStatusApproved
	LoanStatusInvested
	LoanStatusDisbursed
)

func (s LoanStatus) String() string {
	switch s {
	case LoanStatusProposed:
		return "proposed"
	case LoanStatusApproved:
		return "approved"
	case LoanStatusInvested:
		return "invested"
	case LoanStatusDisbursed:
		return "disbursed"
	default:
		return "unknown"
	}
}

type Loan struct {
	bun.BaseModel `bun:"table:loans"`

	ID               uint       `bun:"id,pk,nullzero"`
	UUID             uuid.UUID  `bun:"uuid"`
	BorrowerID       uint       `bun:"borrower_id"`
	ApprovalID       *uint      `bun:"approval_id"`
	DisbursmentID    *uint      `bun:"disbursement_id"`
	ProposedAmount   float64    `bun:"proposed_amount"`
	PrincipalAmount  float64    `bun:"principal_amount"`
	Rate             float64    `bun:"rate"`
	ROI              float64    `bun:"roi"`
	Status           LoanStatus `bun:"status"`
	AgreementFileURL string     `bun:"aggreement_file_url"`
	CreatedAt        time.Time  `bun:"created_at"`
	UpdatedAt        *time.Time `bun:"updated_at,nullzero"`

	Approval    *Approval    `bun:"rel:has-one,join:approval_id=id"`
	Disbursment *Disbursment `bun:"rel:has-one,join:disbursement_id=id"`
}

func NewPropose(
	borowerID uint,
	amount float64,
) *Loan {
	return &Loan{
		UUID:           uuid.New(),
		BorrowerID:     borowerID,
		ProposedAmount: amount,
		Rate:           5,
		ROI:            amount * 5 / 100,
		Status:         LoanStatusProposed,
		CreatedAt:      time.Now(),
	}
}

func (l *Loan) Approve(fieldValidatorId uint, approvalFileUrl string) {
	l.Status = LoanStatusApproved

	l.Approval = &Approval{
		FieldValidatorID: fieldValidatorId,
		ApprovalFileURL:  approvalFileUrl,
		CreatedAt:        time.Now(),
	}
}

func (l *Loan) Invest(amount float64) error {
	if l.PrincipalAmount+amount > l.ProposedAmount {
		return errors.New("loan_invested_amount_exceeds_proposed_amount")
	}

	l.PrincipalAmount += float64(amount)
	if l.PrincipalAmount == l.ProposedAmount {
		l.Status = LoanStatusInvested
	}

	return nil
}

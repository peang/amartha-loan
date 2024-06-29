package dto_request

import "mime/multipart"

type ProposeLoanDTO struct {
	BorowwerID uint    `validate:"required"`
	Amount     float64 `validate:"required" json:"amount"`
}

type ApproveLoanDTO struct {
	LoanID           string                `validate:"required"`
	FieldValidatorID uint                  `validate:"required"`
	ProveImage       *multipart.FileHeader `validate:"required"`
}

type ApprovedLoanListDTO struct {
	Page    string
	PerPage string
}

type InvestLoanDTO struct {
	LoanID     string  `validate:"required"`
	InvestorID uint    `validate:"required"`
	Amount     float64 `validate:"required" json:"amount"`
}

type DisburseLoanDTO struct {
	LoanID           string                `validate:"required"`
	FieldOfficerID   uint                  `validate:"required"`
	AggreementLetter *multipart.FileHeader `validate:"required"`
}

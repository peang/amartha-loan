package dto_response

import (
	"time"

	"github.com/peang/amartha-loan-service/models"
)

type loanDetail struct {
	ID              string    `json:"id"`
	BorowwerID      uint      `json:"borowwer_id"`
	ProposedAmount  float64   `json:"proposed_amount"`
	PrincipalAmount float64   `json:"principal_amount"`
	Rate            float64   `json:"rate"`
	ROI             float64   `json:"roi"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
}

func LoanDetailResponse(loan *models.Loan) loanDetail {
	return loanDetail{
		ID:              loan.UUID.String(),
		BorowwerID:      loan.BorrowerID,
		ProposedAmount:  loan.ProposedAmount,
		PrincipalAmount: loan.PrincipalAmount,
		Rate:            loan.Rate,
		ROI:             loan.ROI,
		Status:          loan.Status.String(),
		CreatedAt:       loan.CreatedAt,
	}
}

type loanList struct {
	ID              string  `json:"id"`
	BorowwerID      uint    `json:"borowwer_id"`
	ProposedAmount  float64 `json:"proposed_amount"`
	PrincipalAmount float64 `json:"principal_amount"`
	Rate            float64 `json:"rate"`
	ROI             float64 `json:"roi"`
	Status          string  `json:"status"`
}

func LoanListResponse(loans *[]models.Loan) []loanList {
	var responses = make([]loanList, 0)
	for _, loan := range *loans {
		response := loanList{
			ID:              loan.UUID.String(),
			BorowwerID:      loan.BorrowerID,
			ProposedAmount:  loan.ProposedAmount,
			PrincipalAmount: loan.PrincipalAmount,
			Rate:            loan.Rate,
			ROI:             loan.ROI,
			Status:          loan.Status.String(),
		}
		responses = append(responses, response)
	}
	return responses
}

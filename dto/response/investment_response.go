package dto_response

import (
	"time"

	"github.com/peang/amartha-loan-service/models"
)

type invstmentDetail struct {
	Loan      loanDetail `json:"loan"`
	Amount    float64    `json:"amount"`
	ROI       float64    `json:"roi"`
	CreatedAt time.Time  `json:"created_at"`
}

func InvestmentDetailResponse(investment *models.Investment) invstmentDetail {
	return invstmentDetail{
		Loan:      LoanDetailResponse(investment.Loan),
		Amount:    investment.Amount,
		ROI:       investment.ROI,
		CreatedAt: investment.CreatedAt,
	}
}

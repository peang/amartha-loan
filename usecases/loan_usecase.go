package usecases

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync"

	"github.com/gotidy/ptr"
	dto_request "github.com/peang/amartha-loan-service/dto/request"
	"github.com/peang/amartha-loan-service/models"
	"github.com/peang/amartha-loan-service/repositories"
	"github.com/peang/amartha-loan-service/services"
	"github.com/peang/amartha-loan-service/services/file_services"
)

type LoanUsecaseInterface interface {
	Propose(ctx context.Context, dto *dto_request.ProposeLoanDTO) (*models.Loan, error)
	Approve(ctx context.Context, dto *dto_request.ApproveLoanDTO) (*models.Loan, error)
	GetAvailableLoans(ctx context.Context, dto *dto_request.ApprovedLoanListDTO) (*[]models.Loan, int, error)
	Invest(ctx context.Context, dto *dto_request.InvestLoanDTO) (*models.Investment, error)
}

type loanUsecase struct {
	loanRepository       repositories.LoanRepositoryInterface
	investmentRepository repositories.InvestmentRepositoryInterface
	fileService          file_services.FileServiceInterface
}

func NewLoanUsecase(
	loanRepository repositories.LoanRepositoryInterface,
	investmentRepository repositories.InvestmentRepositoryInterface,
	fileService file_services.FileServiceInterface,
) LoanUsecaseInterface {
	return &loanUsecase{
		loanRepository:       loanRepository,
		investmentRepository: investmentRepository,
		fileService:          fileService,
	}
}

func (u *loanUsecase) Propose(ctx context.Context, dto *dto_request.ProposeLoanDTO) (*models.Loan, error) {
	loan := models.NewPropose(dto.BorowwerID, dto.Amount)

	aggreementPdfUrl, err := services.GenerateAgreementPDF(loan.UUID.String())
	if err != nil {
		return nil, err
	}

	loan.AgreementFileURL = *aggreementPdfUrl

	loan, err = u.loanRepository.Save(nil, ctx, loan)
	if err != nil {
		return nil, err
	}

	return loan, nil
}

func (u *loanUsecase) Approve(ctx context.Context, dto *dto_request.ApproveLoanDTO) (*models.Loan, error) {
	loan, err := u.loanRepository.Detail(ctx, dto.LoanID)
	if err != nil {
		return nil, err
	}

	if loan == nil {
		return nil, errors.New("loan_not_found")
	}

	if loan.Status != models.LoanStatusProposed {
		return nil, errors.New("only_proposed_loan_allowed")
	}

	approvalFileUrl, err := u.fileService.Upload(dto.ProveImage)
	if err != nil {
		return nil, err
	}

	loan.Approve(dto.FieldValidatorID, approvalFileUrl)

	loan, err = u.loanRepository.Save(nil, ctx, loan)
	if err != nil {
		return nil, err
	}

	return loan, nil
}

func (u *loanUsecase) GetAvailableLoans(ctx context.Context, dto *dto_request.ApprovedLoanListDTO) (*[]models.Loan, int, error) {
	filter := repositories.LoanRepositoryFilter{
		Status: ptr.Of(models.LoanStatusApproved),
	}

	page, err := strconv.Atoi(dto.Page)
	if err != nil || page < 1 {
		page = 1
	}

	perPage, err := strconv.Atoi(dto.PerPage)
	if err != nil || perPage < 1 || perPage > 100 {
		perPage = 10
	}

	loans, count, err := u.loanRepository.List(ctx, page, perPage, "created_at", filter)
	if err != nil {
		return nil, 0, err
	}

	return loans, count, nil
}

func (u *loanUsecase) Invest(ctx context.Context, dto *dto_request.InvestLoanDTO) (*models.Investment, error) {
	loan, err := u.loanRepository.Detail(ctx, dto.LoanID)
	if err != nil {
		return nil, err
	}

	if loan == nil {
		return nil, errors.New("loan_not_found")
	}

	if loan.Status != models.LoanStatusApproved {
		return nil, errors.New("only_approved_loan_allowed")
	}

	investment, err := models.NewInvestment(dto.InvestorID, dto.Amount, loan)
	if err != nil {
		return nil, err
	}

	investment, err = u.investmentRepository.Save(ctx, investment)
	if err != nil {
		return nil, err
	}

	if loan.Status == models.LoanStatusInvested {
		go u.SendEmailToInvestors(ctx, loan)
	}

	return investment, nil
}

func (u *loanUsecase) SendEmailToInvestors(ctx context.Context, loan *models.Loan) error {
	page := 1

	investorChan := make(chan models.Investment, 10)
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go SendEmailWorker(&wg, investorChan)
	}

	for {
		investments, _, err := u.investmentRepository.List(ctx, page, 10, "created_at", repositories.InvestmentRepositoryFilter{
			LoanID: &loan.ID,
		})
		if err != nil {
			fmt.Println(err)
		}

		if investments != nil && len(*investments) == 0 {
			break
		}

		for _, investment := range *investments {
			investorChan <- investment
		}

		page++
	}

	close(investorChan)
	wg.Wait()

	u.investmentRepository.UpdateMany(ctx, repositories.InvestmentRepositoryFilter{
		LoanID: &loan.ID,
	}, repositories.InvestmentRepositoryValues{
		SendAggreementEmail: ptr.Of(true),
	})

	return nil
}

func SendEmailWorker(wg *sync.WaitGroup, ch <-chan models.Investment) {
	defer wg.Done()

	for investment := range ch {
		fmt.Printf("Sending email to %s\n", investment.Investor.Email)
	}
}

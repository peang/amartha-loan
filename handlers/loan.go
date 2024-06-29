package handlers

import (
	"encoding/json"
	"net/http"
	"path/filepath"

	"github.com/labstack/echo/v4"
	dto_request "github.com/peang/amartha-loan-service/dto/request"
	dto_response "github.com/peang/amartha-loan-service/dto/response"
	middleware "github.com/peang/amartha-loan-service/middlewares"
	"github.com/peang/amartha-loan-service/usecases"
	"github.com/peang/amartha-loan-service/utils"
)

type loanHandler struct {
	loanUseCase usecases.LoanUsecaseInterface
}

func NewLoanHandler(
	e *echo.Echo,
	middleware *middleware.Middleware,
	loanUseCase usecases.LoanUsecaseInterface,
) {
	handler := &loanHandler{
		loanUseCase: loanUseCase,
	}

	loanGroup := e.Group("/loans", middleware.JWTAuth(), middleware.RBACMiddleware())

	// For Borowwer User
	loanGroup.POST("/propose", handler.propose)

	// For Field Validator User
	loanGroup.POST("/:id/approve", handler.approve)

	// For Investor user
	loanGroup.GET("/available", handler.getListAvailable)
	loanGroup.POST("/:id/invest", handler.invest)
}

func (h *loanHandler) propose(ctx echo.Context) error {
	context := ctx.Get("payload").(utils.Payload)
	var payload struct {
		Amount float64 `validate:"required"`
	}

	// This also could use Validator v10 to validate
	err := json.NewDecoder(ctx.Request().Body).Decode(&payload)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.Error{
			Code:  http.StatusBadRequest,
			Error: "Invalid Payload",
		})
	}

	dto := dto_request.ProposeLoanDTO{
		BorowwerID: context.ID,
		Amount:     payload.Amount,
	}

	loan, err := h.loanUseCase.Propose(ctx.Request().Context(), &dto)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, utils.Error{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})
	}

	return ctx.JSON(http.StatusCreated, utils.Response{
		Message: "Loan Created",
		Data:    dto_response.LoanDetailResponse(loan),
	})
}

func (h *loanHandler) approve(ctx echo.Context) error {
	context := ctx.Get("payload").(utils.Payload)

	dto := dto_request.ApproveLoanDTO{
		LoanID:           ctx.Param("id"),
		FieldValidatorID: context.ID,
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		return err
	}

	files := form.File["file"]
	if len(files) == 0 {
		return ctx.JSON(http.StatusBadRequest, utils.Error{
			Code:  http.StatusBadRequest,
			Error: "No Prove Uploaded",
		})
	}
	dto.ProveImage = files[0]

	allowedExtensions := map[string]bool{
		".pdf":  true,
		".jpeg": true,
	}
	ext := filepath.Ext(dto.ProveImage.Filename)
	if !allowedExtensions[ext] {
		return ctx.JSON(http.StatusBadRequest, utils.Error{
			Code:  http.StatusBadRequest,
			Error: "FILE TYPE NOT ALLOWED",
		})
	}

	loan, err := h.loanUseCase.Approve(ctx.Request().Context(), &dto)
	if err != nil {
		return ctx.JSON(utils.GetErrorCode(err.Error()), utils.Error{
			Code:  utils.GetErrorCode(err.Error()),
			Error: err.Error(),
		})
	}

	return ctx.JSON(http.StatusCreated, utils.Response{
		Message: "Loan Approved",
		Data:    dto_response.LoanDetailResponse(loan),
	})
}

func (h *loanHandler) getListAvailable(ctx echo.Context) error {
	dto := dto_request.ApprovedLoanListDTO{
		Page:    ctx.QueryParam("page"),
		PerPage: ctx.QueryParam("per_page"),
	}

	loans, count, err := h.loanUseCase.GetAvailableLoans(ctx.Request().Context(), &dto)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, utils.Error{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, utils.Response{
		Message: "Loan List",
		Data:    dto_response.LoanListResponse(loans),
		Meta:    utils.GenerateMeta(dto.Page, dto.PerPage, count),
	})
}

func (h *loanHandler) invest(ctx echo.Context) error {
	context := ctx.Get("payload").(utils.Payload)

	var payload struct {
		Amount float64 `validate:"required"`
	}

	// This also could use Validator v10 to validate
	err := json.NewDecoder(ctx.Request().Body).Decode(&payload)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.Error{
			Code:  http.StatusBadRequest,
			Error: "Invalid Payload",
		})
	}

	dto := dto_request.InvestLoanDTO{
		LoanID:     ctx.Param("id"),
		InvestorID: context.ID,
		Amount:     payload.Amount,
	}

	investment, err := h.loanUseCase.Invest(ctx.Request().Context(), &dto)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, utils.Error{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})
	}

	return ctx.JSON(http.StatusCreated, utils.Response{
		Message: "Invest Success",
		Data:    dto_response.InvestmentDetailResponse(investment),
	})
}

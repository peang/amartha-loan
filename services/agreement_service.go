package services

import (
	"fmt"
	"io/ioutil"

	"github.com/jung-kurt/gofpdf"
)

func GenerateAgreementPDF(loanId string) (*string, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")

	filepath := fmt.Sprintf("file_uploads/loan_%s.pdf", loanId)

	pdf.SetTitle("Loan Aggreement", true)
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Hello, World!")

	tmpFile, err := ioutil.TempFile("", "*.pdf")
	if err != nil {
		return nil, err
	}
	defer tmpFile.Close()

	err = pdf.Output(tmpFile)
	if err != nil {
		return nil, err
	}

	return &filepath, nil
}

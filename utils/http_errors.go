package utils

var httpErrors = map[string]int{
	// Loans Error
	"loan_not_found":             404,
	"only_proposed_loan_allowed": 400,
}

func GetErrorCode(err string) int {
	if httpErrors[err] == 0 {
		return 500
	}

	return httpErrors[err]
}

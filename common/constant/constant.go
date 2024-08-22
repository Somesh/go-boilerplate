package constant

import "errors"

var (
	//LIST OF ERROR
	ErrorDBConn         = errors.New("Unable to find DB Connection")
	ErrorNoResult       = errors.New("Query returns no result")
	ErrNotFound         = errors.New("Not Found")
	ErrInvalidInput     = errors.New("Invalid Input Parameters")
	ErrLocation         = errors.New("Please try againg after sometime")

	ErrUnexpected   = errors.New("Unexpected error has occured")
	ErrConnectivity = errors.New("Please try again after sometime.")
	ErrMaxPageSize  = errors.New("Exceeded maximum page size limit of 20")
	ErrPromoFailure = errors.New("Please try again after sometime.")

	ErrInValidPriceRange = errors.New("Invalid Min or Max Price")
	ErrInValidDateRange  = errors.New("Invalid Date Range")
)

const (

	// Environment
	ENV_DEVELOPMENT = "development"
	ENV_STAGING     = "staging"
	ENV_PRODUCTION  = "production"

	//Flags
	ACTIVE          = 1
	DEACTIVE        = 0
	BLANK           = ""
	DateFormat      = "2006-01-02"
	DateTimeFormat  = "2006-01-02 15:04:05"
	Seconds24HRS    = 86400
	VAR_UPPER_LIMIT = 255


	LOC_DEFAULT      = "UTC"
	LOC_LOCAL        = "Local"
	LOC_JAKARTA      = "Asia/Jakarta"


)
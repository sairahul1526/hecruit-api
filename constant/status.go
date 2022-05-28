package constant

// server status codes
const (
	StatusCodeOk             = "200"
	StatusCodeCreated        = "201"
	StatusCodeBadRequest     = "400"
	StatusCodeForbidden      = "403"
	StatusCodeSessionExpired = "440"
	StatusCodeServerError    = "500"
	StatusCodeDuplicateEntry = "1000"
)

// type of alerts for frontend to show
const (
	NoDialog   = "0"
	ShowDialog = "1"
	ShowToast  = "2"
)

// user status
const (
	UserActive   = "1"
	UserArchived = "2"
	UserBlocked  = "3"
)

// company status
const (
	CompanyActive   = "1"
	CompanyArchived = "2"
)

// location status
const (
	LocationActive   = "1"
	LocationArchived = "2"
)

// teams status
const (
	TeamActive   = "1"
	TeamArchived = "2"
)

// job status
const (
	JobDraft    = "0"
	JobActive   = "1"
	JobArchived = "2"
)

// job status status
const (
	JobStatusActive   = "1"
	JobStatusArchived = "2"
)

// job status type
const (
	JobStatusApplied  = "1"
	JobStatusHired    = "2"
	JobStatusRejected = "3"
)

package constant

// required fields for api endpoints
var (
	ApplicationUpdateRequiredFields = []string{"job_id", "application_id", "status"}
	JobAddRequiredFields            = []string{"company_id", "team_id", "name", "description", "employment_type", "country", "remote_option"}
	TeamAddRequiredFields           = []string{"company_id", "name"}
	UserLoginRequiredFields         = []string{"email", "password"}
	UserSignUpRequiredFields        = []string{"name", "email", "password", "company_name"}
)

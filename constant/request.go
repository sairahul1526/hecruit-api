package constant

// required fields for api endpoints
var (
	ApplicationAddRequiredFields = []string{"job_id", "company_id", "name", "email", "resume"}
	JobAddRequiredFields         = []string{"company_id", "team_id", "name", "description", "employment_type", "location_id", "remote_option"}
	LocationAddRequiredFields    = []string{"company_id", "name"}
	TeamAddRequiredFields        = []string{"company_id", "name"}
	UserLoginRequiredFields      = []string{"email", "password"}
	UserInviteRequiredFields     = []string{"email"}
	UserSignUpRequiredFields     = []string{"name", "email", "password", "company_name", "company_jobs_link"}
)

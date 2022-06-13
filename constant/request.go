package constant

// required fields for api endpoints
var (
	ApplicationAddRequiredFields     = []string{"job_id", "company_id", "name", "email", "resume"}
	InterviewAddRequiredFields       = []string{"application_id", "job_id", "title", "organizer", "attendees", "start_at", "end_at"}
	JobAddRequiredFields             = []string{"company_id", "team_id", "name", "description", "employment_type", "location_id", "remote_option"}
	LocationAddRequiredFields        = []string{"company_id", "name"}
	NoteAddRequiredFields            = []string{"application_id", "job_id", "note"}
	TeamAddRequiredFields            = []string{"company_id", "name"}
	UserLoginRequiredFields          = []string{"email", "password"}
	UserInviteRequiredFields         = []string{"email"}
	UserForgotPasswordRequiredFields = []string{"email"}
	UserSignUpRequiredFields         = []string{"name", "email", "password", "company_name", "company_jobs_link"}
)

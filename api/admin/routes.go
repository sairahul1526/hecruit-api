package admin

import "github.com/gorilla/mux"

// LoadAdminRoutes - load all admin routes with admin prefix
func LoadAdminRoutes(router *mux.Router) {
	adminRoutes := router.PathPrefix("/admin").Subrouter()

	// middlewares
	adminRoutes.Use(APIKeyMiddleware)
	adminRoutes.Use(CheckAuthToken)
	adminRoutes.Use(CheckUserValid)

	// application
	adminRoutes.HandleFunc("/application", ApplicationGet).Queries(
		"application_id", "{application_id}",
	).Methods("GET")
	adminRoutes.HandleFunc("/application", ApplicationsGet).Methods("GET")
	adminRoutes.HandleFunc("/application", ApplicationUpdate).Queries(
		"application_id", "{application_id}",
	).Methods("PUT")
	adminRoutes.HandleFunc("/application-move", ApplicationMove).Queries(
		"application_id", "{application_id}",
	).Methods("PUT")

	// company
	adminRoutes.HandleFunc("/company", CompanyGet).Methods("GET")
	adminRoutes.HandleFunc("/company", CompanyUpdate).Methods("PUT")

	// interview
	adminRoutes.HandleFunc("/interview", InterviewsGet).Methods("GET")
	adminRoutes.HandleFunc("/interview", InterviewAdd).Methods("POST")
	adminRoutes.HandleFunc("/interview", InterviewCancel).Queries(
		"interview_id", "{interview_id}",
	).Methods("DELETE")

	// invoice
	adminRoutes.HandleFunc("/invoice", InvoicesGet).Methods("GET")

	// home
	adminRoutes.HandleFunc("/home", Home).Methods("GET")

	// job
	adminRoutes.HandleFunc("/job", JobGet).Queries(
		"job_id", "{job_id}",
	).Methods("GET")
	adminRoutes.HandleFunc("/job", JobsGet).Queries(
		"team_id", "{team_id}",
		"status", "{status}",
	).Methods("GET")
	adminRoutes.HandleFunc("/job", JobAdd).Methods("POST")
	adminRoutes.HandleFunc("/job", JobUpdate).Queries(
		"job_id", "{job_id}",
	).Methods("PUT")

	// job
	adminRoutes.HandleFunc("/job-status", JobStatusAdd).Queries(
		"job_id", "{job_id}",
	).Methods("POST")
	adminRoutes.HandleFunc("/job-status", JobStatusUpdate).Queries(
		"job_id", "{job_id}",
	).Methods("PUT")

	// user
	adminRoutes.HandleFunc("/user", UsersGet).Queries(
		"company_id", "{company_id}",
	).Methods("GET")
	adminRoutes.HandleFunc("/user", UserGet).Methods("GET")
	adminRoutes.HandleFunc("/user", UserUpdate).Methods("PUT")
	adminRoutes.HandleFunc("/user-invite", UserInvite).Methods("POST")
	adminRoutes.HandleFunc("/user-maintain", UserMaintain).Queries(
		"user_id", "{user_id}",
		"company_id", "{company_id}",
	).Methods("PUT")
	adminRoutes.HandleFunc("/login", UserLogin).Methods("POST")
	adminRoutes.HandleFunc("/signup", UserSignUp).Methods("POST")
	adminRoutes.HandleFunc("/refresh-token", UserRefreshToken).Methods("GET")
	adminRoutes.HandleFunc("/email-verify", UserEmailVerify).Queries(
		"token", "{token}",
	).Methods("GET")
	adminRoutes.HandleFunc("/reset-password", UserResetPassword).Methods("POST")

	// note
	adminRoutes.HandleFunc("/note", NoteAdd).Methods("POST")

	// location
	adminRoutes.HandleFunc("/location", LocationsGet).Methods("GET")
	adminRoutes.HandleFunc("/location", LocationAdd).Methods("POST")
	adminRoutes.HandleFunc("/location", LocationUpdate).Queries(
		"location_id", "{location_id}",
	).Methods("PUT")

	// team
	adminRoutes.HandleFunc("/team", TeamsGet).Methods("GET")
	adminRoutes.HandleFunc("/team", TeamAdd).Methods("POST")
	adminRoutes.HandleFunc("/team", TeamUpdate).Queries(
		"team_id", "{team_id}",
	).Methods("PUT")

	// upload
	adminRoutes.HandleFunc("/upload-signed-url", UploadSignedURL).Queries(
		"file_type", "{file_type}",
		"path_type", "{path_type}",
	).Methods("GET")
}

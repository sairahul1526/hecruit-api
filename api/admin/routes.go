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
	adminRoutes.HandleFunc("/application-move", ApplicationMove).Queries(
		"application_id", "{application_id}",
	).Methods("PUT")

	// company
	adminRoutes.HandleFunc("/company", CompanyGet).Methods("GET")
	adminRoutes.HandleFunc("/company", CompanyUpdate).Methods("PUT")

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
	adminRoutes.HandleFunc("/job-status", JobStatusGet).Queries(
		"job_id", "{job_id}",
	).Methods("GET")

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

}

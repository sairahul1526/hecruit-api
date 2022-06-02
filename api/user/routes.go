package user

import "github.com/gorilla/mux"

// LoadUserRoutes - load all user routes with user prefix
func LoadUserRoutes(router *mux.Router) {
	adminRoutes := router.PathPrefix("/user").Subrouter()

	// application
	adminRoutes.HandleFunc("/application", ApplicationAdd).Methods("POST")

	// company
	adminRoutes.HandleFunc("/company", CompanyGet).Queries(
		"company_url", "{company_url}",
	).Methods("GET")

	// location
	adminRoutes.HandleFunc("/location", LocationsGet).Queries(
		"company_id", "{company_id}",
	).Methods("GET")

	// meta
	adminRoutes.HandleFunc("/job", JobsGet).Queries(
		"company_id", "{company_id}",
	).Methods("GET")
	adminRoutes.HandleFunc("/job", JobGet).Queries(
		"job_id", "{job_id}",
	).Methods("GET")

	// team
	adminRoutes.HandleFunc("/team", TeamsGet).Queries(
		"company_id", "{company_id}",
	).Methods("GET")

	// upload
	adminRoutes.HandleFunc("/upload-signed-url", UploadSignedURL).Queries(
		"file_type", "{file_type}",
	).Methods("GET")

}

package miscellaneous

import "github.com/gorilla/mux"

// LoadMiscellaneousRoutes - load all miscellaneous routes with miscellaneous prefix
func LoadMiscellaneousRoutes(router *mux.Router) {
	adminRoutes := router.PathPrefix("/miscellaneous").Subrouter()

	// meta
	adminRoutes.HandleFunc("/meta", MetaGet).Methods("GET")

}

package api

import (
	"encoding/json"
	"net/http"

	AdminAPI "hecruit-backend/api/admin"
	MiscellaneousAPI "hecruit-backend/api/miscellaneous"
	UserAPI "hecruit-backend/api/user"

	"github.com/gorilla/mux"
)

// HealthCheck .
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	// for load balancer/beanstalk to know whether server/ec2 is healthy
	json.NewEncoder(w).Encode("ok")
}

// LoaderIO .
func LoaderIO(w http.ResponseWriter, r *http.Request) {
	// for loader io verification
	w.Write([]byte("loaderio-e20e1e5221282763725d65fd70174666"))
	// json.NewEncoder(w).Encode("loaderio-e20e1e5221282763725d65fd70174666")
}

// LoadRouter - get mux router with all the routes
func LoadRouter() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/loaderio-e20e1e5221282763725d65fd70174666/", LoaderIO).Methods("GET")

	AdminAPI.LoadAdminRoutes(router)
	MiscellaneousAPI.LoadMiscellaneousRoutes(router)
	UserAPI.LoadUserRoutes(router)

	router.Path("/").HandlerFunc(HealthCheck).Methods("GET")

	return router
}

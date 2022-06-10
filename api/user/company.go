package user

import (
	CONFIG "hecruit-backend/config"
	CONSTANT "hecruit-backend/constant"
	DB "hecruit-backend/database"
	"net/http"
	"strings"

	UTIL "hecruit-backend/util"
)

func CompanyGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// get company details
	company, err := DB.SelectSQL(CONSTANT.CompaniesTable, []string{"id", "name", "description", "logo", "website", "banner", "status"}, map[string]string{"jobs_link": r.FormValue("company_url")})
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}
	if len(company) == 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.CompanyNotFoundMessage, CONSTANT.ShowDialog, response)
		return
	}
	if !strings.EqualFold(company[0]["status"], CONSTANT.CompanyActive) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.CompanyNotActiveMessage, CONSTANT.ShowDialog, response)
		return
	}

	// get locations in a company
	locations, err := DB.SelectSQL(CONSTANT.LocationsTable, []string{"id", "name"}, map[string]string{"company_id": company[0]["id"]})
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}

	// get teams in a company
	teams, err := DB.SelectSQL(CONSTANT.TeamsTable, []string{"id", "name"}, map[string]string{"company_id": company[0]["id"]})
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}

	// get active jobs in a company
	jobs, err := DB.SelectProcess("select id, name, employment_type, team_id, location_id, remote_option, created_at from "+CONSTANT.JobsTable+" where company_id = $1 and status = '"+CONSTANT.JobActive+"' order by name asc", company[0]["id"])
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}

	jobsMap := UTIL.ConvertMapToKeyMapArray(jobs, "team_id")

	response["jobs_count"] = len(jobs)
	response["jobs"] = jobsMap
	response["teams"] = teams
	response["locations"] = locations
	response["company"] = company[0]
	response["employment_types"] = CONSTANT.EmploymentTypes
	response["remote_options"] = CONSTANT.RemoteOptions
	response["media_url"] = CONFIG.S3MediaURL
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

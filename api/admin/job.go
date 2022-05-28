package admin

import (
	CONSTANT "hecruit-backend/constant"
	DB "hecruit-backend/database"
	"net/http"
	"strings"

	UTIL "hecruit-backend/util"
)

func JobGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// get job details
	job, err := DB.SelectSQL(CONSTANT.JobsTable, []string{"id", "company_id", "team_id", "name", "description", "employment_type", "salary", "location_id", "remote_option", "status", "created_at"}, map[string]string{"id": r.FormValue("job_id"), "company_id": r.Header.Get("company_id")})
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}
	if len(job) == 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.JobNotFoundMessage, CONSTANT.ShowDialog, response)
		return
	}

	job[0]["team_name"], _ = DB.QueryRowSQL("select name from " + CONSTANT.TeamsTable + " where id = '" + job[0]["team_id"] + "'")
	job[0]["location_name"], _ = DB.QueryRowSQL("select name from " + CONSTANT.LocationsTable + " where id = '" + job[0]["location_id"] + "'")
	job[0]["employment_type_name"] = CONSTANT.EmploymentTypes[job[0]["employment_type"]]
	job[0]["remote_option_name"] = CONSTANT.RemoteOptions[job[0]["remote_option"]]

	response["job"] = job[0]
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func JobsGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// get jobs in a team
	jobs, err := DB.SelectSQL(CONSTANT.JobsTable, []string{"id", "name", "employment_type", "salary", "location_id", "remote_option"}, map[string]string{"team_id": r.FormValue("team_id"), "location_id": r.FormValue("location_id"), "status": r.FormValue("status")})
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}
	// get job, location ids to get details
	jobIDs := UTIL.ExtractValuesFromArrayMap(jobs, "id")

	// get number of applications for each job
	applicationsCount, err := DB.SelectProcess("select job_id, count(*) as applications from " + CONSTANT.ApplicationsTable + " where job_id in ('" + strings.Join(jobIDs, "','") + "') group by job_id")
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}

	applicationsCountMap := UTIL.ConvertMapToKeyMap(applicationsCount, "job_id")

	for _, job := range jobs {
		job["applications"] = applicationsCountMap[job["id"]]["applications"]
	}

	response["jobs"] = jobs
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func JobAdd(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// read request body
	body, err := UTIL.ReadRequestBodyToMap(r)
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	body["company_id"] = r.Header.Get("company_id")

	// check for required fields
	fieldCheck := UTIL.RequiredFiledsCheck(body, CONSTANT.JobAddRequiredFields)
	if len(fieldCheck) > 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, fieldCheck+" required", CONSTANT.ShowDialog, response)
		return
	}

	// create job
	jobID, _, err := DB.InsertWithUniqueID(CONSTANT.JobsTable, map[string]string{
		"company_id":      body["company_id"],
		"team_id":         body["team_id"],
		"name":            body["name"],
		"description":     body["description"],
		"employment_type": body["employment_type"],
		"location_id":     body["location_id"],
		"remote_option":   body["remote_option"],
		"salary":          body["salary"],
		"status":          body["status"],
	}, "id")
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}

	response["job_id"] = jobID

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func JobUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// read request body
	body, err := UTIL.ReadRequestBodyToMap(r)
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	body["company_id"] = r.Header.Get("company_id")

	_, err = DB.UpdateSQL(CONSTANT.JobsTable, map[string]string{
		"id": r.FormValue("job_id"),
	}, map[string]string{
		"company_id":      body["company_id"],
		"team_id":         body["team_id"],
		"name":            body["name"],
		"description":     body["description"],
		"employment_type": body["employment_type"],
		"location_id":     body["location_id"],
		"remote_option":   body["remote_option"],
		"salary":          body["salary"],
		"status":          body["status"],
	})
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

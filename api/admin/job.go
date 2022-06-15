package admin

import (
	CONSTANT "hecruit-backend/constant"
	DB "hecruit-backend/database"
	"net/http"
	"strconv"
	"strings"

	UTIL "hecruit-backend/util"
)

func JobGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// get job details
	job, err := DB.SelectSQL(CONSTANT.JobsTable, []string{"id", "company_id", "team_id", "name", "description", "employment_type", "salary", "location_id", "remote_option", "status", "created_at", "number_of_positions", "hiring_manager", "responsibilities", "requirements"}, map[string]string{"id": r.FormValue("job_id"), "company_id": r.Header.Get("company_id")})
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}
	if len(job) == 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.JobNotFoundMessage, CONSTANT.ShowDialog, response)
		return
	}

	// get job status
	jobStatus, err := DB.SelectProcess("select id, name, type from "+CONSTANT.JobStatusTable+" where job_id = $1 and status = '"+CONSTANT.JobStatusActive+`' order by "order"`, r.FormValue("job_id"))
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}

	// get company details
	company, err := DB.SelectProcess("select name, jobs_link from " + CONSTANT.CompaniesTable + " where id = '" + job[0]["company_id"] + "'")
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}

	job[0]["company_jobs_link"] = company[0]["jobs_link"]
	job[0]["team_name"], _ = DB.QueryRowSQL("select name from " + CONSTANT.TeamsTable + " where id = '" + job[0]["team_id"] + "'")
	job[0]["location_name"], _ = DB.QueryRowSQL("select name from " + CONSTANT.LocationsTable + " where id = '" + job[0]["location_id"] + "'")
	job[0]["employment_type_name"] = CONSTANT.EmploymentTypes[job[0]["employment_type"]]
	job[0]["remote_option_name"] = CONSTANT.RemoteOptions[job[0]["remote_option"]]

	response["job"] = job[0]
	response["job_status"] = jobStatus
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func JobsGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// build query
	wheres := []string{
		" company_id = $1 ",
	}
	queryArgs := []interface{}{
		r.Header.Get("company_id"),
	}
	i := 2
	for key, val := range r.URL.Query() {
		switch key {
		case "team_id":
			if len(val[0]) > 0 {
				wheres = append(wheres, " team_id = $"+strconv.Itoa(i))
				queryArgs = append(queryArgs, val[0])
				i++
			}
		case "location_id":
			if len(val[0]) > 0 {
				wheres = append(wheres, " location_id = $"+strconv.Itoa(i))
				queryArgs = append(queryArgs, val[0])
				i++
			}
		case "status":
			if len(val[0]) > 0 {
				wheres = append(wheres, " status = $"+strconv.Itoa(i))
				queryArgs = append(queryArgs, val[0])
				i++
			}
		}
	}

	where := ""
	if len(wheres) > 0 {
		where = " where " + strings.Join(wheres, " and ")
	}

	// get jobs
	jobs, err := DB.SelectProcess("select id, name, employment_type, salary, location_id, team_id, remote_option from "+CONSTANT.JobsTable+where+" order by name asc", queryArgs...)
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}

	if strings.EqualFold(r.FormValue("applications"), "true") {
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
		"company_id":          body["company_id"],
		"team_id":             body["team_id"],
		"name":                body["name"],
		"description":         body["description"],
		"employment_type":     body["employment_type"],
		"location_id":         body["location_id"],
		"remote_option":       body["remote_option"],
		"salary":              body["salary"],
		"status":              body["status"],
		"number_of_positions": body["number_of_positions"],
		"hiring_manager":      body["hiring_manager"],
		"responsibilities":    body["responsibilities"],
		"requirements":        body["requirements"],
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

	_, err = DB.UpdateSQL(CONSTANT.JobsTable, map[string]string{
		"id":         r.FormValue("job_id"),
		"company_id": r.Header.Get("company_id"),
	}, map[string]string{
		"team_id":             body["team_id"],
		"name":                body["name"],
		"description":         body["description"],
		"employment_type":     body["employment_type"],
		"location_id":         body["location_id"],
		"remote_option":       body["remote_option"],
		"salary":              body["salary"],
		"status":              body["status"],
		"number_of_positions": body["number_of_positions"],
		"hiring_manager":      body["hiring_manager"],
		"responsibilities":    body["responsibilities"],
		"requirements":        body["requirements"],
	})
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

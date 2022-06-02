package user

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
	job, err := DB.SelectSQL(CONSTANT.JobsTable, []string{"id", "company_id", "team_id", "name", "description", "employment_type", "salary", "location_id", "remote_option", "status", "created_at", "responsibilities", "requirements"}, map[string]string{"id": r.FormValue("job_id")})
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

	// build query
	wheres := []string{}
	queryArgs := []interface{}{}
	i := 1
	for key, val := range r.URL.Query() {
		switch key {
		case "team_id":
			if len(val[0]) > 0 {
				wheres = append(wheres, " team_id = $"+strconv.Itoa(i))
				queryArgs = append(queryArgs, val[0])
				i++
			}
		case "name":
			if len(val[0]) > 0 {
				wheres = append(wheres, " name ilike '%%"+val[0]+"%%' ")
			}
		case "employment_type":
			if len(val[0]) > 0 {
				wheres = append(wheres, " employment_type = $"+strconv.Itoa(i))
				queryArgs = append(queryArgs, val[0])
				i++
			}
		case "remote_option":
			if len(val[0]) > 0 {
				wheres = append(wheres, " remote_option = $"+strconv.Itoa(i))
				queryArgs = append(queryArgs, val[0])
				i++
			}
		case "location_id":
			if len(val[0]) > 0 {
				wheres = append(wheres, " location_id = $"+strconv.Itoa(i))
				queryArgs = append(queryArgs, val[0])
				i++
			}
		case "company_id":
			if len(val[0]) > 0 {
				wheres = append(wheres, " company_id = $"+strconv.Itoa(i))
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

	// get active jobs in a company
	jobs, err := DB.SelectProcess("select id, name, employment_type, team_id, location_id, remote_option, created_at from "+CONSTANT.JobsTable+where+" order by name asc", queryArgs...)
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}

	jobsMap := UTIL.ConvertMapToKeyMapArray(jobs, "team_id")

	response["jobs_count"] = len(jobs)
	response["jobs"] = jobsMap
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

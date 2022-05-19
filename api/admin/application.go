package admin

import (
	CONSTANT "hecruit-backend/constant"
	DB "hecruit-backend/database"
	"net/http"
	"strconv"
	"strings"

	UTIL "hecruit-backend/util"
)

func ApplicationGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// get application details
	application, err := DB.SelectSQL(CONSTANT.ApplicationsTable, []string{"id", "job_id", "name", "email", "created_at", "phone", "current_company", "current_salary", "expected_salary", "notice_period", "total_experience", "location", "linkedin_link", "twitter_link", "github_link", "portfolio_link", "other_link", "cover", "additional_information", "resume", "rating"}, map[string]string{"id": r.FormValue("application_id")})
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}
	if len(application) == 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.ApplicationNotFoundMessage, CONSTANT.ShowDialog, response)
		return
	}

	// get job, team details
	job, err := DB.SelectSQL(CONSTANT.JobsTable, []string{"team_id", "name"}, map[string]string{"id": application[0]["job_id"]})
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}
	if len(job) == 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.JobNotFoundMessage, CONSTANT.ShowDialog, response)
		return
	}

	team, err := DB.SelectSQL(CONSTANT.TeamsTable, []string{"name"}, map[string]string{"id": job[0]["team_id"]})
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}
	if len(team) == 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.TeamNotFoundMessage, CONSTANT.ShowDialog, response)
		return
	}

	// get all activities done on this application
	activities, err := DB.SelectProcess("select status, created_at from "+CONSTANT.ApplicationActivitiesTable+" where application_id = $1", r.FormValue("application_id"))
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}

	application[0]["team_name"] = team[0]["name"]
	application[0]["job_name"] = job[0]["name"]
	response["application"] = application[0]
	response["activities"] = activities
	response["media_url"] = CONSTANT.MediaURL
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func ApplicationsGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// build query
	wheres := []string{}
	queryArgs := []interface{}{}
	i := 1
	for key, val := range r.URL.Query() {
		switch key {
		case "name":
			if len(val[0]) > 0 {
				wheres = append(wheres, " name ilike '%%"+val[0]+"%%' ")
			}
		case "phone":
			if len(val[0]) > 0 {
				wheres = append(wheres, " phone ilike '%%"+val[0]+"%%' ")
			}
		case "email":
			if len(val[0]) > 0 {
				wheres = append(wheres, " email ilike '%%"+val[0]+"%%' ")
			}
		case "job_id":
			if len(val[0]) > 0 {
				wheres = append(wheres, " job_id = $"+strconv.Itoa(i))
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
	// get applications
	applications, err := DB.SelectProcess("select id, name, email, location, status, rating, created_at, updated_at from "+CONSTANT.ApplicationsTable+where+" order by created_at desc limit "+strconv.Itoa(CONSTANT.ResultsPerPageAdmin)+" offset "+strconv.Itoa((UTIL.GetPageNumber(r.FormValue("page"))-1)*CONSTANT.ResultsPerPageAdmin), queryArgs...)
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}
	// get status ids to get details
	status := UTIL.ExtractValuesFromArrayMap(applications, "status")

	// get status names
	statusNames, err := DB.SelectProcess("select id, name from " + CONSTANT.JobStatusTable + " where id in ('" + strings.Join(status, "','") + "')")
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}

	// get total number of applications
	applicationsCount, err := DB.SelectProcess("select count(*) as ctn from "+CONSTANT.ApplicationsTable+where, queryArgs...)
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}

	if len(r.FormValue("job_id")) > 0 {
		// get number of applications for each status
		applicationsByStatus, err := DB.SelectProcess("select status, count(*) as ctn from "+CONSTANT.ApplicationsTable+" where job_id = $1 group by status order by status", r.FormValue("job_id"))
		if err != nil {
			UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
			return
		}

		response["applications_by_status"] = UTIL.ConvertMapToKeyMap(applicationsByStatus, "status")
	}

	statusNamesMap := UTIL.ConvertMapToKeyMap(statusNames, "id")

	for _, application := range applications {
		application["status_name"] = statusNamesMap[application["status"]]["name"]
	}

	response["applications"] = applications
	response["applications_count"] = applicationsCount[0]["ctn"]
	response["no_pages"] = strconv.Itoa(UTIL.GetNumberOfPages(applicationsCount[0]["ctn"], CONSTANT.ResultsPerPageAdmin))

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func ApplicationMove(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// read request body
	body, err := UTIL.ReadRequestBodyToMap(r)
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	body["company_id"] = r.Header.Get("company_id")

	// TODO get application status and get next status to move to that
	// and also check for rejection query param and move to rejected status
	// check for required fields
	fieldCheck := UTIL.RequiredFiledsCheck(body, CONSTANT.ApplicationUpdateRequiredFields)
	if len(fieldCheck) > 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, fieldCheck+" required", CONSTANT.ShowDialog, response)
		return
	}

	_, err = DB.InsertSQL(CONSTANT.ApplicationActivitiesTable, map[string]string{
		"job_id":         body["job_id"],
		"application_id": body["application_id"],
		"status":         body["status"],
	})
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

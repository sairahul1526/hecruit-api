package admin

import (
	CONFIG "hecruit-backend/config"
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
	application, err := DB.SelectSQL(CONSTANT.ApplicationsTable, []string{"id", "job_id", "name", "email", "created_at", "phone", "current_company", "current_salary", "expected_salary", "notice_period", "total_experience", "location", "linkedin_link", "twitter_link", "github_link", "portfolio_link", "other_link", "cover", "additional_information", "resume", "rating", "status"}, map[string]string{"id": r.FormValue("application_id")})
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
	activities, err := DB.SelectProcess("select status, created_at, created_by from "+CONSTANT.ApplicationActivitiesTable+" where application_id = $1 order by created_at", r.FormValue("application_id"))
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}
	// get status, user ids to get details
	status := UTIL.ExtractValuesFromArrayMap(activities, "status")
	userIDs := UTIL.ExtractValuesFromArrayMap(activities, "created_by")

	// get status names
	statusNames, err := DB.SelectProcess("select id, name, type from " + CONSTANT.JobStatusTable + " where id in ('" + strings.Join(status, "','") + "')")
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}

	// get users
	users, err := DB.SelectProcess("select id, name from " + CONSTANT.UsersTable + " where id in ('" + strings.Join(userIDs, "','") + "')")
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}

	statusNamesMap := UTIL.ConvertMapToKeyMap(statusNames, "id")
	usersMap := UTIL.ConvertMapToKeyMap(users, "id")

	for _, activity := range activities {
		activity["status_name"] = statusNamesMap[activity["status"]]["name"]
		activity["status_type"] = statusNamesMap[activity["status"]]["type"]
		activity["created_by_name"] = usersMap[activity["created_by"]]["name"]
	}

	application[0]["team_name"] = team[0]["name"]
	application[0]["job_name"] = job[0]["name"]
	response["application"] = application[0]
	response["activities"] = activities
	response["media_url"] = CONFIG.S3MediaURL
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
	statusNames, err := DB.SelectProcess("select id, name, type from " + CONSTANT.JobStatusTable + " where id in ('" + strings.Join(status, "','") + "')")
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
		application["status_type"] = statusNamesMap[application["status"]]["type"]
	}

	response["applications"] = applications
	response["applications_count"] = applicationsCount[0]["ctn"]
	response["no_pages"] = strconv.Itoa(UTIL.GetNumberOfPages(applicationsCount[0]["ctn"], CONSTANT.ResultsPerPageAdmin))

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func ApplicationUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// read request body
	body, err := UTIL.ReadRequestBodyToMap(r)
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	_, err = DB.UpdateSQL(CONSTANT.ApplicationsTable, map[string]string{
		"id":         r.FormValue("application_id"),
		"company_id": r.Header.Get("company_id"),
	}, map[string]string{
		"rating": body["rating"],
	})
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func ApplicationMove(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	application, err := DB.SelectSQL(CONSTANT.ApplicationsTable, []string{"job_id", "company_id"}, map[string]string{"id": r.FormValue("application_id")})
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}
	if len(application) == 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.ApplicationNotFoundMessage, CONSTANT.ShowDialog, response)
		return
	}

	lastestApplicationActivity, err := DB.QueryRowSQL("select status from "+CONSTANT.ApplicationActivitiesTable+" where application_id = $1 order by created_at desc limit 1", r.FormValue("application_id"))
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}
	if len(lastestApplicationActivity) == 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.ApplicationStatusNotFoundMessage, CONSTANT.ShowDialog, response)
		return
	}

	// get status
	var nextApplicationStatus string
	if len(r.FormValue("reject")) > 0 {
		// get reject status
		nextApplicationStatus, err = DB.QueryRowSQL("select id from " + CONSTANT.JobStatusTable + " where job_id = '" + application[0]["job_id"] + "' and type = " + CONSTANT.JobStatusRejected + " limit 1")
		if err != nil {
			UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
			return
		}
	} else {
		// go to next state, except reject
		nextApplicationStatus, err = DB.QueryRowSQL("select id from " + CONSTANT.JobStatusTable + " where job_id = '" + application[0]["job_id"] + "' and type != " + CONSTANT.JobStatusRejected + ` and "order" > (select "order" from ` + CONSTANT.JobStatusTable + " where id = '" + lastestApplicationActivity + `' limit 1) order by "order" asc limit 1`)
		if err != nil {
			UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
			return
		}
	}
	if len(nextApplicationStatus) == 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.ApplicationStatusNotFoundMessage, CONSTANT.ShowDialog, response)
		return
	}

	_, _, err = DB.InsertWithUniqueID(CONSTANT.ApplicationActivitiesTable, map[string]string{
		"job_id":         application[0]["job_id"],
		"application_id": r.FormValue("application_id"),
		"status":         nextApplicationStatus,
		"company_id":     r.Header.Get("company_id"),
		"created_by":     r.Header.Get("user_id"),
	}, "id")
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}

	DB.UpdateSQL(CONSTANT.ApplicationsTable, map[string]string{
		"id": r.FormValue("application_id"),
	}, map[string]string{
		"status": nextApplicationStatus,
	})

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

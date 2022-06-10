package admin

import (
	CONSTANT "hecruit-backend/constant"
	DB "hecruit-backend/database"
	"net/http"

	UTIL "hecruit-backend/util"
)

func Home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// since user already checked in middleware, no need to check for zero size
	user, err := DB.SelectSQL(CONSTANT.UsersTable, []string{"name"}, map[string]string{"id": r.Header.Get("user_id")})
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}

	jobsDetails, err := DB.SelectProcess("select count(id) as active_jobs, sum(number_of_positions) as openings_to_fill  from "+CONSTANT.JobsTable+" where company_id = $1 and status = '"+CONSTANT.JobActive+"'", r.Header.Get("company_id"))
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}

	applicationsCount, err := DB.SelectProcess("select count(*) as total_applications from "+CONSTANT.ApplicationsTable+" where company_id = $1", r.Header.Get("company_id"))
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}

	// applicationsDetails, err := DB.SelectProcess("select j.type, count(*) as ctn from "+CONSTANT.ApplicationsTable+" a left join "+CONSTANT.JobStatusTable+" j on a.status = j.id where a.company_id = $1 group by j.type;", r.Header.Get("company_id"))
	// if err != nil {
	// 	UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
	// 	return
	// }

	// applicationsDetailsMap := UTIL.ConvertMapToKeyMap(applicationsDetails, "type")

	response["active_jobs"] = "0"
	response["openings_to_fill"] = "0"
	if len(jobsDetails) > 0 {
		response["active_jobs"] = jobsDetails[0]["active_jobs"]
		response["openings_to_fill"] = jobsDetails[0]["openings_to_fill"]
	}
	response["total_applications"] = "0"
	if len(applicationsCount) > 0 {
		response["total_applications"] = applicationsCount[0]["total_applications"]
	}

	response["user_name"] = user[0]["name"]
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

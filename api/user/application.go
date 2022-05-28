package user

import (
	CONSTANT "hecruit-backend/constant"
	DB "hecruit-backend/database"
	"net/http"

	UTIL "hecruit-backend/util"
)

func ApplicationAdd(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// read request body
	body, err := UTIL.ReadRequestBodyToMap(r)
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	// check for required fields
	fieldCheck := UTIL.RequiredFiledsCheck(body, CONSTANT.ApplicationAddRequiredFields)
	if len(fieldCheck) > 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, fieldCheck+" required", CONSTANT.ShowDialog, response)
		return
	}

	// get first status
	status, err := DB.QueryRowSQL("select id from "+CONSTANT.JobStatusTable+` where job_id = $1 order by "order" limit 1`, body["job_id"])
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}

	// create application
	applicationID, _, err := DB.InsertWithUniqueID(CONSTANT.ApplicationsTable, map[string]string{
		"job_id":                 body["job_id"],
		"company_id":             body["company_id"],
		"name":                   body["name"],
		"email":                  body["email"],
		"phone":                  body["phone"],
		"current_company":        body["current_company"],
		"current_salary":         body["current_salary"],
		"expected_salary":        body["expected_salary"],
		"notice_period":          body["notice_period"],
		"total_experience":       body["total_experience"],
		"location":               body["location"],
		"linkedin_link":          body["linkedin_link"],
		"twitter_link":           body["twitter_link"],
		"github_link":            body["github_link"],
		"portfolio_link":         body["portfolio_link"],
		"other_link":             body["other_link"],
		"cover":                  body["cover"],
		"additional_information": body["additional_information"],
		"resume":                 body["resume"],
		"status":                 status,
	}, "id")
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}

	DB.InsertWithUniqueID(CONSTANT.ApplicationActivitiesTable, map[string]string{
		"application_id": applicationID,
		"job_id":         body["job_id"],
		"company_id":     body["company_id"],
		"status":         status,
	}, "id")

	response["application_id"] = applicationID

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

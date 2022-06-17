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

	// get company name
	companyName, _ := DB.QueryRowSQL("select name from "+CONSTANT.CompaniesTable+" where id = $1", body["company_id"])
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}
	// get job name
	jobName, _ := DB.QueryRowSQL("select name from "+CONSTANT.JobsTable+" where id = $1", body["job_id"])
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}
	if len(companyName) > 0 && len(jobName) > 0 {

		// send email to applicant
		DB.InsertWithUniqueID(CONSTANT.EmailsTable, map[string]string{
			"from":     companyName + " <" + CONSTANT.NoReplyEmail + ">",
			"to":       body["email"],
			"reply_to": CONSTANT.SupportEmail,
			"title":    "Thank you for your application to " + companyName,
			"body": `Hi ` + body["name"] + `,

			<br><br>

			Thank you for your interest in ` + companyName + `! We wanted to let you know we received your application for ` + jobName + `, and we are delighted that you would consider joining our team.
			
			<br><br>

			Our team will review your application and will be in touch if your qualifications match our needs for the role. If you are not selected for this position, keep an eye on our jobs page as we're growing and adding openings.
			
			<br><br>
			
			Best,

			<br>

			` + companyName + ` Team`,
			"company_id":     body["company_id"],
			"job_id":         body["job_id"],
			"application_id": applicationID,
			"status":         CONSTANT.EmailTobeSent,
		}, "id")
	}

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

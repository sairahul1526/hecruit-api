package admin

import (
	CONSTANT "hecruit-backend/constant"
	DB "hecruit-backend/database"
	"net/http"

	UTIL "hecruit-backend/util"
)

func InterviewAdd(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// read request body
	body, err := UTIL.ReadRequestBodyToMap(r)
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	// check for required fields
	fieldCheck := UTIL.RequiredFiledsCheck(body, CONSTANT.InterviewAddRequiredFields)
	if len(fieldCheck) > 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, fieldCheck+" required", CONSTANT.ShowDialog, response)
		return
	}

	// create interview
	interviewID, _, err := DB.InsertWithUniqueID(CONSTANT.InterviewsTable, map[string]string{
		"company_id":     r.Header.Get("company_id"),
		"job_id":         body["job_id"],
		"application_id": body["application_id"],
		"title":          body["title"],
		"organizer":      body["organizer"],
		"attendees":      body["attendees"],
		"start_at":       body["start_at"],
		"end_at":         body["end_at"],
		"created_by":     r.Header.Get("user_id"),
		"status":         CONSTANT.InterviewActive,
	}, "id")
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}

	response["interview_id"] = interviewID

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

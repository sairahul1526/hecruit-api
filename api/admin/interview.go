package admin

import (
	CONSTANT "hecruit-backend/constant"
	DB "hecruit-backend/database"
	"net/http"
	"os"
	"strconv"
	"strings"

	UTIL "hecruit-backend/util"
)

func InterviewsGet(w http.ResponseWriter, r *http.Request) {
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
		case "date":
			if len(val[0]) > 0 {
				temp := strings.Split(val[0], ",")
				if len(temp[0]) > 0 {
					wheres = append(wheres, "$"+strconv.Itoa(i)+" <= start_at ")
					queryArgs = append(queryArgs, temp[0])
					i++
				}
				if len(temp[1]) > 0 {
					wheres = append(wheres, " start_at <= $"+strconv.Itoa(i))
					queryArgs = append(queryArgs, temp[1])
					i++
				}
			}
		case "organizer":
			if len(val[0]) > 0 {
				wheres = append(wheres, " organizer ilike '%%"+val[0]+"%%' ")
			}
		case "attendees":
			if len(val[0]) > 0 {
				wheres = append(wheres, " attendees ilike '%%"+val[0]+"%%' ")
			}
		case "job_id":
			if len(val[0]) > 0 {
				wheres = append(wheres, " job_id = $"+strconv.Itoa(i))
				queryArgs = append(queryArgs, val[0])
				i++
			}
		case "status":
			if len(val[0]) > 0 {
				switch val[0] {
				case "active":
					wheres = append(wheres, " status = '"+CONSTANT.InterviewActive+"' and end_at > now() ")
				case "completed":
					wheres = append(wheres, " status = '"+CONSTANT.InterviewActive+"' and end_at <= now() ")
				case "cancelled":
					wheres = append(wheres, " status = '"+CONSTANT.InterviewCancelled+"' ")
				}
			}
		}
	}

	where := ""
	if len(wheres) > 0 {
		where = " where " + strings.Join(wheres, " and ")
	}
	// get interviews
	interviews, err := DB.SelectProcess("select * from "+CONSTANT.InterviewsTable+where+" order by created_at desc limit "+strconv.Itoa(CONSTANT.ResultsPerPageAdmin)+" offset "+strconv.Itoa((UTIL.GetPageNumber(r.FormValue("page"))-1)*CONSTANT.ResultsPerPageAdmin), queryArgs...)
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}
	// get jobs ids to get details
	jobIDs := UTIL.ExtractValuesFromArrayMap(interviews, "job_id")

	// get job details
	jobs, err := DB.SelectProcess("select id, name from " + CONSTANT.JobsTable + " where id in ('" + strings.Join(jobIDs, "','") + "')")
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}

	// get total number of interviews
	interviewsCount, err := DB.SelectProcess("select count(*) as ctn from "+CONSTANT.InterviewsTable+where, queryArgs...)
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}

	jobsMap := UTIL.ConvertMapToKeyMap(jobs, "id")

	for _, interview := range interviews {
		interview["job_name"] = jobsMap[interview["job_id"]]["name"]
	}

	response["interviews"] = interviews
	response["interviews_count"] = interviewsCount[0]["ctn"]
	response["no_pages"] = strconv.Itoa(UTIL.GetNumberOfPages(interviewsCount[0]["ctn"], CONSTANT.ResultsPerPageAdmin))

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

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

	// check if email is valid, based on regex
	if !UTIL.IsEmailValid(body["organizer"]) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.UseValidEmailMessage, CONSTANT.ShowDialog, response)
		return
	}
	for _, email := range strings.Split(body["attendees"], ",") {
		if !UTIL.IsEmailValid(email) {
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.UseValidEmailMessage, CONSTANT.ShowDialog, response)
			return
		}
	}

	// get company detail
	company, err := DB.SelectSQL(CONSTANT.CompaniesTable, []string{"name"}, map[string]string{"id": r.Header.Get("company_id")})
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}
	if len(company) == 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.CompanyNotFoundMessage, CONSTANT.ShowDialog, response)
		return
	}

	// get job details
	job, err := DB.SelectSQL(CONSTANT.JobsTable, []string{"name"}, map[string]string{"id": body["job_id"]})
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}
	if len(job) == 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.JobNotFoundMessage, CONSTANT.ShowDialog, response)
		return
	}

	// get application details
	application, err := DB.SelectSQL(CONSTANT.ApplicationsTable, []string{"name"}, map[string]string{"id": body["application_id"]})
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}
	if len(application) == 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.ApplicationNotFoundMessage, CONSTANT.ShowDialog, response)
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
		"meeting_link":   body["meeting_link"],
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

	// create and upload ics file
	os.Mkdir(CONSTANT.ICSS3Path, os.ModePerm)
	icsFileName := CONSTANT.ICSS3Path + UTIL.GenerateRandomID() + ".ics"
	if !UTIL.UploadContentAsFile(icsFileName, []byte(UTIL.BuildICSFile(body["meeting_link"], body["organizer"], body["attendees"], interviewID, body["title"], body["start_at"], body["end_at"], "CONFIRMED", "0"))) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}

	// send interview details email
	DB.InsertWithUniqueID(CONSTANT.EmailsTable, map[string]string{
		"from":           company[0]["name"] + " <" + CONSTANT.NoReplyEmail + ">",
		"to":             body["organizer"] + "," + body["attendees"],
		"title":          body["title"],
		"body":           company[0]["name"],
		"company_id":     r.Header.Get("company_id"),
		"job_id":         body["job_id"],
		"application_id": body["application_id"],
		"attachment":     icsFileName,
		"status":         CONSTANT.EmailTobeSent,
	}, "id")

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func InterviewCancel(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// get interview details
	interview, err := DB.SelectSQL(CONSTANT.InterviewsTable, []string{"*"}, map[string]string{"id": r.FormValue("interview_id")})
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}
	if len(interview) == 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.InterviewNotFoundMessage, CONSTANT.ShowDialog, response)
		return
	}

	// get company detail
	company, err := DB.SelectSQL(CONSTANT.CompaniesTable, []string{"name"}, map[string]string{"id": interview[0]["company_id"]})
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}
	if len(company) == 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.CompanyNotFoundMessage, CONSTANT.ShowDialog, response)
		return
	}

	// cancel interview
	_, err = DB.UpdateSQL(CONSTANT.InterviewsTable, map[string]string{
		"id": r.FormValue("interview_id"),
	}, map[string]string{
		"status": CONSTANT.InterviewCancelled,
	})
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}

	// create and upload ics file
	os.Mkdir(CONSTANT.ICSS3Path, os.ModePerm)
	icsFileName := CONSTANT.ICSS3Path + UTIL.GenerateRandomID() + ".ics"
	if !UTIL.UploadContentAsFile(icsFileName, []byte(UTIL.BuildICSFile(interview[0]["meeting_link"], interview[0]["organizer"], interview[0]["attendees"], interview[0]["id"], interview[0]["title"], interview[0]["start_at"], interview[0]["end_at"], "CANCELLED", "1"))) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}

	// send interview details email
	DB.InsertWithUniqueID(CONSTANT.EmailsTable, map[string]string{
		"from":           company[0]["name"] + " <" + CONSTANT.NoReplyEmail + ">",
		"to":             interview[0]["organizer"] + "," + interview[0]["attendees"],
		"title":          "CANCELLED - " + interview[0]["title"],
		"body":           company[0]["name"],
		"company_id":     r.Header.Get("company_id"),
		"job_id":         interview[0]["job_id"],
		"application_id": interview[0]["application_id"],
		"attachment":     icsFileName,
		"status":         CONSTANT.EmailTobeSent,
	}, "id")

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

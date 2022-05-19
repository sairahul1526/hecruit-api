package admin

import (
	CONSTANT "hecruit-backend/constant"
	DB "hecruit-backend/database"
	"net/http"

	UTIL "hecruit-backend/util"
)

func JobStatusGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// get job status
	jobStatus, err := DB.SelectProcess("select id, name from "+CONSTANT.JobStatusTable+" where job_id = $1 and status = '"+CONSTANT.JobStatusActive+`' order by "order"`, r.FormValue("job_id"))
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}

	response["job_status"] = jobStatus
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

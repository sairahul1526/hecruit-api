package admin

import (
	CONSTANT "hecruit-backend/constant"
	DB "hecruit-backend/database"
	"net/http"
	"strconv"

	UTIL "hecruit-backend/util"
)

func JobStatusAdd(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// read request body
	body, err := UTIL.ReadRequestBodyToArrayMap(r)
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	for i, jobStatus := range body {
		if len(jobStatus["name"]) > 0 {
			DB.InsertWithUniqueID(CONSTANT.JobStatusTable, map[string]string{
				"name":   jobStatus["name"],
				"job_id": r.FormValue("job_id"),
				"order":  strconv.Itoa(i + 1),
				"type":   jobStatus["type"],
			}, "id")
		}
	}

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func JobStatusUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// read request body
	body, err := UTIL.ReadRequestBodyToArrayMap(r)
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	for i, jobStatus := range body {
		if len(jobStatus["name"]) > 0 {
			if len(jobStatus["id"]) > 0 {
				// new id
				DB.UpdateSQL(CONSTANT.JobStatusTable, map[string]string{
					"id":     jobStatus["id"],
					"job_id": r.FormValue("job_id"),
				}, map[string]string{
					"name":  jobStatus["name"],
					"order": strconv.Itoa(i + 1),
					"type":  jobStatus["type"],
				})
			} else {
				// new id
				DB.InsertWithUniqueID(CONSTANT.JobStatusTable, map[string]string{
					"name":   jobStatus["name"],
					"job_id": r.FormValue("job_id"),
					"order":  strconv.Itoa(i + 1),
					"type":   CONSTANT.JobStatusProcess, // always a process status when adding
				}, "id")
			}
		}
	}

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

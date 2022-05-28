package user

import (
	CONSTANT "hecruit-backend/constant"
	DB "hecruit-backend/database"
	"net/http"

	UTIL "hecruit-backend/util"
)

func LocationsGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// get locations in a company
	locations, err := DB.SelectSQL(CONSTANT.LocationsTable, []string{"id", "name"}, map[string]string{"company_id": r.FormValue("company_id")})
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}

	response["locations"] = locations
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

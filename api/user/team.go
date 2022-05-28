package user

import (
	CONSTANT "hecruit-backend/constant"
	DB "hecruit-backend/database"
	"net/http"

	UTIL "hecruit-backend/util"
)

func TeamsGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// get teams in a company
	teams, err := DB.SelectSQL(CONSTANT.TeamsTable, []string{"id", "name"}, map[string]string{"company_id": r.FormValue("company_id")})
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}

	response["teams"] = teams
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

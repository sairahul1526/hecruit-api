package admin

import (
	CONSTANT "hecruit-backend/constant"
	DB "hecruit-backend/database"
	"net/http"
	"strings"

	UTIL "hecruit-backend/util"
)

func TeamsGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// get teams in a company
	teams, err := DB.SelectSQL(CONSTANT.TeamsTable, []string{"id", "name"}, map[string]string{"company_id": r.Header.Get("company_id")})
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}

	if strings.EqualFold(r.FormValue("openings"), "true") {
		// get team ids to get details
		teamIDs := UTIL.ExtractValuesFromArrayMap(teams, "id")

		if len(teamIDs) > 0 {
			// get team openings
			openings, err := DB.SelectProcess("select team_id, count(*) as openings from " + CONSTANT.JobsTable + " where team_id in ('" + strings.Join(teamIDs, "','") + "') and status = '" + CONSTANT.JobActive + "' group by team_id")
			if err != nil {
				UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
				return
			}
			openingsMap := UTIL.ConvertMapToKeyMap(openings, "team_id")

			for _, team := range teams {
				team["openings"] = openingsMap[team["id"]]["openings"]
			}
		}
	}

	response["teams"] = teams
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func TeamAdd(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// read request body
	body, err := UTIL.ReadRequestBodyToMap(r)
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	body["company_id"] = r.Header.Get("company_id")

	// check for required fields
	fieldCheck := UTIL.RequiredFiledsCheck(body, CONSTANT.TeamAddRequiredFields)
	if len(fieldCheck) > 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, fieldCheck+" required", CONSTANT.ShowDialog, response)
		return
	}

	// create team
	teamID, _, err := DB.InsertWithUniqueID(CONSTANT.TeamsTable, map[string]string{
		"name":       body["name"],
		"company_id": body["company_id"],
		"status":     CONSTANT.TeamActive,
	}, "id")
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}

	response["team_id"] = teamID

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func TeamUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// read request body
	body, err := UTIL.ReadRequestBodyToMap(r)
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	_, err = DB.UpdateSQL(CONSTANT.TeamsTable, map[string]string{
		"id":         r.FormValue("team_id"),
		"company_id": r.Header.Get("company_id"),
	}, map[string]string{
		"name": body["name"],
	})
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

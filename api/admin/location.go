package admin

import (
	CONSTANT "hecruit-backend/constant"
	DB "hecruit-backend/database"
	"net/http"
	"strings"

	UTIL "hecruit-backend/util"
)

func LocationsGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// get locations in a company
	locations, err := DB.SelectSQL(CONSTANT.LocationsTable, []string{"id", "name"}, map[string]string{"company_id": r.Header.Get("company_id")})
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}

	if strings.EqualFold(r.FormValue("openings"), "true") {
		// get location ids to get details
		locationIDs := UTIL.ExtractValuesFromArrayMap(locations, "id")

		if len(locationIDs) > 0 {
			// get location openings
			openings, err := DB.SelectProcess("select location_id, count(*) as openings from " + CONSTANT.JobsTable + " where location_id in ('" + strings.Join(locationIDs, "','") + "') and status = '" + CONSTANT.JobActive + "' group by location_id")
			if err != nil {
				UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
				return
			}
			openingsMap := UTIL.ConvertMapToKeyMap(openings, "location_id")

			for _, location := range locations {
				location["openings"] = openingsMap[location["id"]]["openings"]
			}
		}
	}

	response["locations"] = locations
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func LocationAdd(w http.ResponseWriter, r *http.Request) {
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
	fieldCheck := UTIL.RequiredFiledsCheck(body, CONSTANT.LocationAddRequiredFields)
	if len(fieldCheck) > 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, fieldCheck+" required", CONSTANT.ShowDialog, response)
		return
	}

	// create location
	locationID, _, err := DB.InsertWithUniqueID(CONSTANT.LocationsTable, map[string]string{
		"name":       body["name"],
		"company_id": body["company_id"],
		"status":     CONSTANT.LocationActive,
	}, "id")
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}

	response["location_id"] = locationID

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func LocationUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// read request body
	body, err := UTIL.ReadRequestBodyToMap(r)
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	_, err = DB.UpdateSQL(CONSTANT.LocationsTable, map[string]string{
		"id":         r.FormValue("location_id"),
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

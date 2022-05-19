package admin

import (
	CONSTANT "hecruit-backend/constant"
	DB "hecruit-backend/database"
	"net/http"
	"strings"

	UTIL "hecruit-backend/util"
)

func CompanyGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// get company details
	company, err := DB.SelectSQL(CONSTANT.CompaniesTable, []string{"name", "description", "logo", "website", "banner", "owner_user_id", "plan", "plan_amount", "payment_user_details", "status", "jobs_link"}, map[string]string{"id": r.FormValue("company_id")})
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}
	if len(company) == 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.CompanyNotFoundMessage, CONSTANT.ShowDialog, response)
		return
	}
	if !strings.EqualFold(company[0]["status"], CONSTANT.CompanyActive) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.CompanyNotActiveMessage, CONSTANT.ShowDialog, response)
		return
	}

	response["company"] = company[0]
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func CompanyUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// read request body
	body, err := UTIL.ReadRequestBodyToMap(r)
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	_, err = DB.UpdateSQL(CONSTANT.CompaniesTable, map[string]string{
		"id": r.FormValue("company_id"),
	}, map[string]string{
		"name":          body["name"],
		"description":   body["description"],
		"logo":          body["logo"],
		"website":       body["website"],
		"banner":        body["banner"],
		"owner_user_id": body["owner_user_id"],
		"jobs_link":     body["jobs_link"],
	})
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

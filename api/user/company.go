package user

import (
	CONFIG "hecruit-backend/config"
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
	company, err := DB.SelectSQL(CONSTANT.CompaniesTable, []string{"id", "name", "description", "logo", "website", "banner", "status"}, map[string]string{"jobs_link": r.FormValue("company_url")})
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
	response["media_url"] = CONFIG.S3MediaURL
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

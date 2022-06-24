package admin

import (
	CONSTANT "hecruit-backend/constant"
	DB "hecruit-backend/database"
	UTIL "hecruit-backend/util"
	"net/http"
)

func InvoicesGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// get company details
	company, err := DB.SelectSQL(CONSTANT.CompaniesTable, []string{"*"}, map[string]string{"id": r.Header.Get("company_id")})
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}
	if len(company) == 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.CompanyNotFoundMessage, CONSTANT.ShowDialog, response)
		return
	}

	// get invoices
	invoices, err := DB.SelectSQL(CONSTANT.InvoicesTable, []string{"*"}, map[string]string{"company_id": company[0]["id"]})
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}

	response["invoices"] = invoices
	response["company"] = company[0]
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

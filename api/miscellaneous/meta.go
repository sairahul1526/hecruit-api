package miscellaneous

import (
	CONSTANT "hecruit-backend/constant"
	"net/http"

	UTIL "hecruit-backend/util"
)

func MetaGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	response["employment_types"] = CONSTANT.EmploymentTypes
	response["remote_options"] = CONSTANT.RemoteOptions

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

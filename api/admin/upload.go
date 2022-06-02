package admin

import (
	CONSTANT "hecruit-backend/constant"
	"net/http"

	UTIL "hecruit-backend/util"
)

func UploadSignedURL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	var path string
	switch r.FormValue("path_type") {
	case "1":
		path = CONSTANT.CompanyS3Path
	case "2":
		path = CONSTANT.UserS3Path
	default:
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	fileName, url, err := UTIL.GeneratePUTSignedURL(path, r.FormValue("file_type"))
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}

	response["file_name"] = fileName
	response["url"] = url

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

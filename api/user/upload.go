package user

import (
	CONSTANT "hecruit-backend/constant"
	"net/http"

	UTIL "hecruit-backend/util"
)

func UploadSignedURL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	fileName, url, err := UTIL.GeneratePUTSignedURL(CONSTANT.ResumeS3Path, r.FormValue("file_type"))
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}

	response["file_name"] = fileName
	response["url"] = url

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

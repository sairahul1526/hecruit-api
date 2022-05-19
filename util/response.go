package util

import (
	"encoding/json"
	CONSTANT "hecruit-backend/constant"
	"net/http"
)

// SetReponse - set request response with status, message
func SetReponse(w http.ResponseWriter, status string, msg string, msgType string, response map[string]interface{}) {
	// cloudflare caches all responses if this is not set
	w.Header().Set("cache-control", "s-maxage=0")
	w.Header().Set("Status", status)
	w.WriteHeader(getHTTPStatusCode(status))
	response["meta"] = setMeta(status, msg, msgType)
	json.NewEncoder(w).Encode(response)
}

// you will always have meta in any response (GET/POST/PUT/PATCH/DELETE)
// status - HTTP status codes like 200,201,400,500,503
// message - Any message which would be used by app to display or take action
// message_type - 1 : show dialog, 2 : show toast, else nothing
func setMeta(status string, msg string, msgType string) map[string]string {
	if len(msg) == 0 {
		if status == CONSTANT.StatusCodeBadRequest {
			msg = "Bad Request"
		} else if status == CONSTANT.StatusCodeServerError {
			msg = "Server Error"
		}
	}
	return map[string]string{
		"status":       status,
		"message":      msg,
		"message_type": msgType,
	}
}

func getHTTPStatusCode(code string) int {
	switch code {
	case CONSTANT.StatusCodeOk:
		return http.StatusOK
	case CONSTANT.StatusCodeCreated:
		return http.StatusCreated
	case CONSTANT.StatusCodeBadRequest:
		return http.StatusBadRequest
	case CONSTANT.StatusCodeServerError:
		return http.StatusInternalServerError
	}
	return http.StatusOK
}

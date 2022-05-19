package admin

import (
	CONSTANT "hecruit-backend/constant"
	DB "hecruit-backend/database"
	"net/http"
	"strings"

	UTIL "hecruit-backend/util"
)

func UserGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// since user already checked in middleware, no need to check for zero size
	user, _ := DB.SelectSQL(CONSTANT.UsersTable, []string{"name", "email", "role", "photo"}, map[string]string{"user_id": r.Header.Get("user_id")})
	response["user"] = user[0]

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func UserUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	if !strings.EqualFold(r.Header.Get("user_id"), r.FormValue("user_id")) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	// read request body
	body, err := UTIL.ReadRequestBodyToMap(r)
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	change := map[string]string{
		"name":  body["name"],
		"role":  body["role"],
		"photo": body["photo"],
	}
	// check if password is being updated
	if len(body["current_password"]) > 0 {
		// get user
		user, err := DB.SelectSQL(CONSTANT.UsersTable, []string{"password"}, map[string]string{"id": r.FormValue("user_id")})
		if err != nil {
			UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
			return
		}
		if !strings.EqualFold(user[0]["password"], UTIL.GetMD5HashString(body["current_password"])) {
			UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, CONSTANT.PasswordIncorrectMessage, CONSTANT.ShowDialog, response)
			return
		}
		change["password"] = UTIL.GetMD5HashString(body["new_password"])
	}

	_, err = DB.UpdateSQL(CONSTANT.UsersTable, map[string]string{"id": r.FormValue("user_id")}, change)
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func UserLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// read request body
	body, err := UTIL.ReadRequestBodyToMap(r)
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	// check for required fields
	fieldCheck := UTIL.RequiredFiledsCheck(body, CONSTANT.UserLoginRequiredFields)
	if len(fieldCheck) > 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, fieldCheck+" required", CONSTANT.ShowDialog, response)
		return
	}

	// check if user id is valid
	user, err := DB.SelectSQL(CONSTANT.UsersTable, []string{"id", "name", "email", "role", "photo", "company_id", "status"}, map[string]string{"email": body["email"], "password": UTIL.GetMD5HashString(body["password"])})
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}
	if len(user) == 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.UserNotExistMessage, CONSTANT.ShowDialog, response)
		return
	}
	if !strings.EqualFold(user[0]["status"], CONSTANT.UserActive) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.UserNotAllowedMessage, CONSTANT.ShowDialog, response)
		return
	}

	// generate access and refresh token
	// access token - jwt token with short expiry added in header for authorization
	// refresh token - jwt token with long expiry to get new access token if expired
	// if refresh token expired, need to login
	accessToken, err := UTIL.CreateAccessToken(map[string]interface{}{"user_id": user[0]["id"], "company_id": user[0]["company_id"]})
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}
	refreshToken, err := UTIL.CreateRefreshToken(map[string]interface{}{"user_id": user[0]["id"], "company_id": user[0]["company_id"]})
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}

	response["user"] = user[0]
	response["access_token"] = accessToken
	response["refresh_token"] = refreshToken

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func UserSignUp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// read request body
	body, err := UTIL.ReadRequestBodyToMap(r)
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	// check for required fields
	fieldCheck := UTIL.RequiredFiledsCheck(body, CONSTANT.UserSignUpRequiredFields)
	if len(fieldCheck) > 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, fieldCheck+" required", CONSTANT.ShowDialog, response)
		return
	}

	// check if email is valid, based on regex
	if !UTIL.IsEmailValid(body["email"]) {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.UseValidEmailMessage, CONSTANT.ShowDialog, response)
		return
	}

	// check if email already exists
	if DB.CheckIfExists(CONSTANT.UsersTable, map[string]string{"email": body["email"]}) == nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.EmailExistMessage, CONSTANT.ShowDialog, response)
		return
	}

	// create company
	companyID, _, err := DB.InsertWithUniqueID(CONSTANT.CompaniesTable, map[string]string{
		"name":   body["company_name"],
		"status": CONSTANT.CompanyActive,
	}, "id")
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}
	delete(body, "company_name")

	// hash password
	body["company_id"] = companyID
	body["password"] = UTIL.GetMD5HashString(body["password"])
	body["status"] = CONSTANT.UserActive

	userID, _, err := DB.InsertWithUniqueID(CONSTANT.UsersTable, body, "id")
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}

	// update company owner
	DB.UpdateSQL(CONSTANT.CompaniesTable, map[string]string{
		"id": companyID,
	}, map[string]string{
		"owner_user_id": userID,
	})

	// check if user id is valid
	user, err := DB.SelectSQL(CONSTANT.UsersTable, []string{"id", "name", "email", "role", "photo", "company_id", "status"}, map[string]string{"id": userID})
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}
	if len(user) == 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.UserNotExistMessage, CONSTANT.ShowDialog, response)
		return
	}

	// generate access and refresh token
	// access token - jwt token with short expiry added in header for authorization
	// refresh token - jwt token with long expiry to get new access token if expired
	// if refresh token expired, need to login
	accessToken, err := UTIL.CreateAccessToken(map[string]interface{}{"user_id": user[0]["id"], "company_id": user[0]["company_id"]})
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}
	refreshToken, err := UTIL.CreateRefreshToken(map[string]interface{}{"user_id": user[0]["id"], "company_id": user[0]["company_id"]})
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}

	response["user"] = user[0]
	response["access_token"] = accessToken
	response["refresh_token"] = refreshToken

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func UserRefreshToken(w http.ResponseWriter, r *http.Request) {

	var response = make(map[string]interface{})

	// check if user id is valid
	user, err := DB.SelectSQL(CONSTANT.UsersTable, []string{"id", "company_id", "status"}, map[string]string{"id": r.Header.Get("user_id")})
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}
	if len(user) == 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.UserNotExistMessage, CONSTANT.ShowDialog, response)
		return
	}

	// refresh token is already checked in middleware
	// generate new access token
	accessToken, err := UTIL.CreateAccessToken(map[string]interface{}{"user_id": user[0]["id"], "company_id": user[0]["company_id"]})
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}

	response["access_token"] = accessToken

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

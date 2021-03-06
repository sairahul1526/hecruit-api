package admin

import (
	CONFIG "hecruit-backend/config"
	CONSTANT "hecruit-backend/constant"
	DB "hecruit-backend/database"
	"net/http"
	"strings"

	UTIL "hecruit-backend/util"

	"github.com/google/uuid"
)

func UserGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// since user already checked in middleware, no need to check for zero size
	user, err := DB.SelectSQL(CONSTANT.UsersTable, []string{"name", "email", "role", "photo", "company_id"}, map[string]string{"id": r.Header.Get("user_id")})
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}

	// get company details
	company, err := DB.SelectProcess("select name, jobs_link, next_bill_date, plan from " + CONSTANT.CompaniesTable + " where id = '" + user[0]["company_id"] + "'")
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}
	if len(company) == 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.CompanyNotFoundMessage, CONSTANT.ShowDialog, response)
		return
	}

	user[0]["company_name"] = company[0]["name"]
	user[0]["company_plan"] = company[0]["plan"]
	user[0]["company_jobs_link"] = company[0]["jobs_link"]
	user[0]["company_next_bill_date"] = company[0]["next_bill_date"]
	response["user"] = user[0]
	response["media_url"] = CONFIG.S3MediaURL
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func UsersGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// get users of a company
	users, err := DB.SelectSQL(CONSTANT.UsersTable, []string{"id", "name", "email", "role", "photo", "status"}, map[string]string{"company_id": r.Header.Get("company_id")})
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}

	response["users"] = users
	response["media_url"] = CONFIG.S3MediaURL
	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func UserInvite(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// read request body
	body, err := UTIL.ReadRequestBodyToMap(r)
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	// check for required fields
	fieldCheck := UTIL.RequiredFiledsCheck(body, CONSTANT.UserInviteRequiredFields)
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
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.UserAlreadySignedUpMessage, CONSTANT.ShowDialog, response)
		return
	}

	password := strings.Split(uuid.New().String(), "-")[0]
	// create user
	_, _, err = DB.InsertWithUniqueID(CONSTANT.UsersTable, map[string]string{
		"email":          body["email"],
		"company_id":     r.Header.Get("company_id"),
		"password":       UTIL.GetMD5HashString(password),
		"status":         CONSTANT.UserActive,
		"email_verified": CONSTANT.EmailVerified, // since we are sending login details, it's already verified
	}, "id")
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}

	// send login details email
	DB.InsertWithUniqueID(CONSTANT.EmailsTable, map[string]string{
		"from":     "Hecruit <" + CONSTANT.NoReplyEmail + ">",
		"to":       body["email"],
		"reply_to": CONSTANT.SupportEmail,
		"title":    "Login credentials - Hecruit",
		"body": `Hey,

		<br><br>

		Below are the credentials to login to your Hecruit account.

		<br><br>

		Website: https://admin.hecruit.com/

		<br>
		
		Email: ` + body["email"] + `

		<br>
		
		Password: ` + password + `

		<br><br>

		Change your password in the settings.

		<br><br>
		
		Best,

		<br>

		Hecruit Team`,
		"company_id": r.Header.Get("company_id"),
		"status":     CONSTANT.EmailTobeSent,
	}, "id")

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func UserUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

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
	if len(body["new_password"]) > 0 {
		// get user
		user, err := DB.SelectSQL(CONSTANT.UsersTable, []string{"password"}, map[string]string{"id": r.Header.Get("user_id")})
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

	_, err = DB.UpdateSQL(CONSTANT.UsersTable, map[string]string{"id": r.Header.Get("user_id")}, change)
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}

func UserMaintain(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// read request body
	body, err := UTIL.ReadRequestBodyToMap(r)
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	_, err = DB.UpdateSQL(CONSTANT.UsersTable, map[string]string{
		"id":         r.FormValue("user_id"),
		"company_id": r.FormValue("company_id"),
	}, map[string]string{
		"status": body["status"],
	})
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
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.IncorrectCredentialsExistMessage, CONSTANT.ShowDialog, response)
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

	// check if jobs link already exists
	if DB.CheckIfExists(CONSTANT.CompaniesTable, map[string]string{"jobs_link": body["company_jobs_link"]}) == nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.JobsPageExistMessage, CONSTANT.ShowDialog, response)
		return
	}

	// create company
	companyID, _, err := DB.InsertWithUniqueID(CONSTANT.CompaniesTable, map[string]string{
		"name":      body["company_name"],
		"jobs_link": body["company_jobs_link"],
		"status":    CONSTANT.CompanyActive,
	}, "id")
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}

	userID, _, err := DB.InsertWithUniqueID(CONSTANT.UsersTable, map[string]string{
		"company_id": companyID,
		"password":   UTIL.GetMD5HashString(body["password"]),
		"status":     CONSTANT.UserActive,
		"name":       body["name"],
		"email":      body["email"],
	}, "id")
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

	emailToken, err := UTIL.CreateAccessToken(map[string]interface{}{"id": user[0]["id"], "email": body["email"]})
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}

	// send verify email
	DB.InsertWithUniqueID(CONSTANT.EmailsTable, map[string]string{
		"from":     "Hecruit <" + CONSTANT.NoReplyEmail + ">",
		"to":       body["email"],
		"reply_to": CONSTANT.SupportEmail,
		"title":    "Verify " + body["company_name"] + " jobs page - Hecruit",
		"body": `Hey ` + body["name"] + `,

		<br><br>

		Please click on this link (or copy & paste) to verify your account and activate your ` + body["company_name"] + ` job page ????:

		<br><br>
		
		https://api.hecruit.com/admin/email-verify?token=` + emailToken + `

		<br><br>
		
		Best,

		<br>
		
		Hecruit Team`,
		"company_id": companyID,
		"status":     CONSTANT.EmailTobeSent,
	}, "id")

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

func UserEmailVerify(w http.ResponseWriter, r *http.Request) {

	// var response = make(map[string]interface{})

	// parse token and get email
	data, err := UTIL.ParseJWTToken("Token " + r.FormValue("token"))
	if err != nil {
		w.Write([]byte("Server Error"))
		// UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}

	// check if email is valid
	user, err := DB.SelectSQL(CONSTANT.UsersTable, []string{"id", "name", "email", "company_id"}, map[string]string{"id": data["id"].(string)})
	if err != nil {
		w.Write([]byte("Server Error"))
		// UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}
	if len(user) == 0 {
		w.Write([]byte(CONSTANT.UserNotExistMessage))
		// UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.UserNotExistMessage, CONSTANT.ShowDialog, response)
		return
	}

	// update user email verified
	DB.UpdateSQL(CONSTANT.UsersTable, map[string]string{
		"id": user[0]["id"],
	}, map[string]string{
		"email_verified": CONSTANT.EmailVerified,
	})

	// send welcome email when verified
	DB.InsertWithUniqueID(CONSTANT.EmailsTable, map[string]string{
		"from":  "Rahul <" + CONSTANT.RahulEmail + ">",
		"to":    user[0]["email"],
		"title": "Welcome to Hecruit",
		"body": `Hi ` + user[0]["name"] + `,
		
		<br><br>
		
		Rahul here, founder of Hecruit. I want to thank you for joining us and say hello ????.
		
		<br><br>
		
		It???s essential for us that we build a product that you???ll love to use, so if you have questions or feedback just reply to this. This is my email, so I'll try to reply back as soon as possible.
		
		<br><br>

		You can also reach us at support@hecruit.com and on Twitter at @hecruitHQ.
		
		<br><br>

		Thanks so much for joining us!<br>
		Rahul (@sairahul1)<br>
		Founder of Hecruit`,
		"company_id": user[0]["company_id"],
		"status":     CONSTANT.EmailTobeSent,
	}, "id")

	http.Redirect(w, r, "https://admin.hecruit.com/", http.StatusSeeOther)
	// UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "Email successfully Verified", CONSTANT.ShowDialog, response)
}

func UserResetPassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	// read request body
	body, err := UTIL.ReadRequestBodyToMap(r)
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, "", CONSTANT.ShowDialog, response)
		return
	}

	// check for required fields
	fieldCheck := UTIL.RequiredFiledsCheck(body, CONSTANT.UserForgotPasswordRequiredFields)
	if len(fieldCheck) > 0 {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, fieldCheck+" required", CONSTANT.ShowDialog, response)
		return
	}

	// check if email already exists
	if DB.CheckIfExists(CONSTANT.UsersTable, map[string]string{"email": body["email"]}) != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, CONSTANT.UserNotExistMessage, CONSTANT.ShowDialog, response)
		return
	}

	password := strings.Split(uuid.New().String(), "-")[0]
	// update password
	_, err = DB.UpdateSQL(CONSTANT.UsersTable, map[string]string{
		"email": body["email"],
	}, map[string]string{
		"password": UTIL.GetMD5HashString(password),
	})
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeServerError, "", CONSTANT.ShowDialog, response)
		return
	}

	// send login details email
	DB.InsertWithUniqueID(CONSTANT.EmailsTable, map[string]string{
		"from":     "Hecruit <" + CONSTANT.NoReplyEmail + ">",
		"to":       body["email"],
		"reply_to": CONSTANT.SupportEmail,
		"title":    "New login credentials - Hecruit",
		"body": `Hey,

		<br><br>

		Below are the credentials to login to your Hecruit account.

		<br><br>

		Website: https://admin.hecruit.com/

		<br>
		
		Email: ` + body["email"] + `

		<br>
		
		Password: ` + password + `

		<br><br>

		Change your password in the settings.

		<br><br>
		
		Best,

		<br>

		Hecruit Team`,
		"company_id": r.Header.Get("company_id"),
		"status":     CONSTANT.EmailTobeSent,
	}, "id")

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "New password has been sent to your email. Login.", CONSTANT.ShowDialog, response)
}

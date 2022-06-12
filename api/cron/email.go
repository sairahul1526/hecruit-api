package cron

import (
	CONSTANT "hecruit-backend/constant"
	DB "hecruit-backend/database"
	"strings"

	UTIL "hecruit-backend/util"
)

func SendEmailsContinously() {
	defer SendEmailsContinously()
	for {
		// get all emails which are not sent
		emails, err := DB.SelectProcess("select * from " + CONSTANT.EmailsTable + " where status = '" + CONSTANT.EmailTobeSent + "' limit 10")
		if err != nil || len(emails) == 0 { // stop if no emails found
			continue
		}

		UTIL.SendEmails(emails)

		emailIDs := UTIL.ExtractValuesFromArrayMap(emails, "id")

		// update messages to sent status
		DB.ExecuteSQL("update " + CONSTANT.EmailsTable + " set status = " + CONSTANT.EmailSent + " where id in ('" + strings.Join(emailIDs, "','") + "')")
	}
}
